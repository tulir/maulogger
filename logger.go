// mauLogger - A logger for Go programs
// Copyright (C) 2016 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package maulogger ...
package maulogger

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// LoggerFileFormat ...
type LoggerFileFormat func(now string, i int) string

// Logger ...
type Logger struct {
	PrintLevel         int
	FlushLineThreshold int
	FileTimeFormat     string
	FileFormat         LoggerFileFormat
	TimeFormat         string
	FileMode           os.FileMode
	DefaultSub         *Sublogger

	writer *bufio.Writer
	lines  int
}

// GeneralLogger contains advanced logging functions and also implements io.Writer
type GeneralLogger interface {
	Write(p []byte) (n int, err error)
	Log(level Level, parts ...interface{})
	Logln(level Level, parts ...interface{})
	Logf(level Level, message string, args ...interface{})
	Debug(parts ...interface{})
	Debugln(parts ...interface{})
	Debugf(message string, args ...interface{})
	Info(parts ...interface{})
	Infoln(parts ...interface{})
	Infof(message string, args ...interface{})
	Warn(parts ...interface{})
	Warnln(parts ...interface{})
	Warnf(message string, args ...interface{})
	Error(parts ...interface{})
	Errorln(parts ...interface{})
	Errorf(message string, args ...interface{})
	Fatal(parts ...interface{})
	Fatalln(parts ...interface{})
	Fatalf(message string, args ...interface{})
}

// Create a Logger
func Create() *Logger {
	var log = &Logger{
		PrintLevel:         10,
		FileTimeFormat:     "2006-01-02",
		FileFormat:         func(now string, i int) string { return fmt.Sprintf("%[1]s-%02[2]d.log", now, i) },
		TimeFormat:         "15:04:05 02.01.2006",
		FileMode:           0600,
		FlushLineThreshold: 5,
		lines:              0,
	}
	log.DefaultSub = log.CreateSublogger("", LevelInfo)
	return log
}

// SetWriter formats the given parts with fmt.Sprint and log them with the SetWriter level
func (log *Logger) SetWriter(w *bufio.Writer) {
	log.writer = w
}

// OpenFile formats the given parts with fmt.Sprint and log them with the OpenFile level
func (log *Logger) OpenFile() error {
	now := time.Now().Format(log.FileTimeFormat)
	i := 1
	for ; ; i++ {
		if _, err := os.Stat(log.FileFormat(now, i)); os.IsNotExist(err) {
			break
		} else if i == 99 {
			i = 1
			break
		}
	}
	file, err := os.OpenFile(log.FileFormat(now, i), os.O_WRONLY|os.O_CREATE|os.O_APPEND, log.FileMode)
	if err != nil {
		return err
	} else if file == nil {
		return os.ErrInvalid
	}
	log.writer = bufio.NewWriter(file)
	return nil
}

// Close formats the given parts with fmt.Sprint and log them with the Close level
func (log *Logger) Close() {
	if log.writer != nil {
		log.writer.Flush()
	}
}

// Raw formats the given parts with fmt.Sprint and log them with the Raw level
func (log *Logger) Raw(level Level, module, message string) {
	var msg []byte
	if len(module) == 0 {
		msg = []byte(fmt.Sprintf("[%s] [%s] %s", time.Now().Format(log.TimeFormat), level.Name, message))
	} else {
		msg = []byte(fmt.Sprintf("[%s] [%s/%s] %s", time.Now().Format(log.TimeFormat), module, level.Name, message))
	}

	if log.writer != nil {
		_, err := log.writer.Write(msg)
		if err != nil {
			fmt.Println("Failed to write to log file:", err)
			return
		}
		log.lines++
		if log.lines == log.FlushLineThreshold {
			log.lines = 0
			log.writer.Flush()
		}
	}

	if level.Severity >= log.PrintLevel {
		if level.Severity >= LevelError.Severity {
			os.Stderr.Write(level.GetColor())
			os.Stderr.Write(msg)
			os.Stderr.Write(level.GetReset())
		} else {
			os.Stdout.Write(level.GetColor())
			os.Stdout.Write(msg)
			os.Stdout.Write(level.GetReset())
		}
	}
}
