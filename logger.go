// mauLogger - A logger for Go programs
// Copyright (c) 2020 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package maulogger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"
)

type BasicLogger struct {
	Sublogger

	TimeFormat  string
	PrintLevel  int
	FileLevel   int
	FileMode    os.FileMode
	FileName    string
	MaxFiles    int
	MaxFileSize int64

	writer   io.Writer
	file     *os.File
	fileSize int64
}

// Logger contains advanced logging functions and also implements io.Writer
type Logger interface {
	Sub(module ...string) Logger
	WithDefaultLevel(level Level) Logger
	GetParent() Logger

	Write(p []byte) (n int, err error)

	Log(level Level, parts ...interface{})
	Logf(level Level, message string, args ...interface{})

	Debug(parts ...interface{})
	Debugf(message string, args ...interface{})
	Info(parts ...interface{})
	Infof(message string, args ...interface{})
	Warn(parts ...interface{})
	Warnf(message string, args ...interface{})
	Error(parts ...interface{})
	Errorf(message string, args ...interface{})
	Fatal(parts ...interface{})
	Fatalf(message string, args ...interface{})
}

// Create a Logger
func Create() Logger {
	var log = &BasicLogger{
		PrintLevel:  LevelWarn.Severity,
		FileLevel:   LevelInfo.Severity,
		TimeFormat:  "15:04:05 02.01.2006",
		FileName:    "mau.log",
		FileMode:    0600,
		MaxFiles:    10,
		MaxFileSize: 10 * 1024 * 1024,
	}
	log.Sublogger = Sublogger{
		topLevel:     log,
		parent:       nil,
		Module:       "",
		DefaultLevel: LevelDebug,
	}
	return log
}

func (log *BasicLogger) formatName(x int) string {
	if x == 0 {
		return log.FileName
	}
	return fmt.Sprintf("%s.%d", log.FileName, x)
}

func (log *BasicLogger) rename(x int) {
	name := log.formatName(x)
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return
	}
	if x >= log.MaxFiles-1 {
		err := os.Remove(name)
		if err != nil {
			log.Errorf("Failed to remove log file %s: %v", name, err)
		}
	} else {
		log.rename(x + 1)
		err := os.Rename(name, log.formatName(x+1))
		if err != nil {
			log.Errorf("Failed to rotate log file %s: %v", name, err)
		}
	}
}

func (log *BasicLogger) rotate() {
	var buf bytes.Buffer
	log.writer = &buf
	err := log.file.Close()
	if err != nil {
		log.Error("Failed to close old log file:", err)
	}
	log.rename(0)
	log.file, err = os.OpenFile(log.FileName, os.O_WRONLY|os.O_CREATE, log.FileMode)
	if err != nil {
		log.Error("Failed to open new log file:", err)
	}
	log.fileSize = 0
	log.writer = log.file
	n, err := buf.WriteTo(log.writer)
	log.fileSize += n
}

// OpenFile formats the given parts with fmt.Sprint and logs the result with the OpenFile level
func (log *BasicLogger) Open() error {
	info, _ := os.Stat(log.FileName)
	if info != nil {
		log.fileSize = info.Size()
	}
	var err error
	log.file, err = os.OpenFile(log.FileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, log.FileMode)
	if err != nil {
		return err
	} else if log.file == nil {
		return os.ErrInvalid
	}
	log.writer = log.file
	return nil
}

// Close formats the given parts with fmt.Sprint and logs the result with the Close level
func (log *BasicLogger) Close() error {
	if log.file != nil {
		return log.file.Close()
	}
	return nil
}

// Raw formats the given parts with fmt.Sprint and logs the result with the Raw level
func (log *BasicLogger) Raw(level Level, module, message string) {
	var mod string
	if len(module) == 0 {
		mod = module + "/" + level.Name
	} else {
		mod = level.Name
	}

	if level.Severity >= log.PrintLevel {
		var file io.Writer
		if level.Severity >= LevelError.Severity {
			file = os.Stderr
		} else {
			file = os.Stdout
		}
		_, _ = fmt.Fprintf(file, "%s[%s] [%s]%s %s", level.GetColor(), time.Now().Format(log.TimeFormat), mod, level.GetReset(), message)
	}

	if log.writer != nil && level.Severity >= log.FileLevel {
		n, _ := fmt.Fprintf(log.writer, "[%s] [%s] %s", time.Now().Format(log.TimeFormat), mod, message)
		log.fileSize += int64(n)
		if log.fileSize > log.MaxFileSize {
			log.rotate()
		}
	}
}
