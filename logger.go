package maulog

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Level is the severity level of a log entry.
type Level int

// ToString converts a level to a string
func (lvl Level) ToString() string {
	switch lvl {
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "GENERIC"
	}
}

// GetColor gets the ANSI escape color code for the log level.
func (lvl Level) GetColor() []byte {
	switch lvl {
	case Warn:
		return []byte("\x1b[33m")
	case Error:
		return []byte("\x1b[31m")
	case Fatal:
		return []byte("\x1b[35m")
	default:
		return []byte{}
	}
}

const (
	// Info is the level for basic log messages.
	Info Level = iota
	// Warn is the level saying that something went wrong, but the program will continue operating mostly normally.
	Warn Level = iota
	// Error is the level saying that something went wrong and the program may not operate as expected, but will still continue.
	Error Level = iota
	// Fatal is the level saying that something went wrong and the program will not operate normally.
	Fatal Level = iota
)

var (
	// FileTimeformat is the time format used in log file names.
	FileTimeformat = "2006-01-02"
	// Fileformat is the format used for log file names.
	Fileformat = "logs/%[1]s-%[2]d.log"
	// Timeformat is the time format used in logging.
	Timeformat = "15:04:05 02.01.2006"
)

var writer *bufio.Writer
var lines int

func init() {
	now := time.Now().Format(FileTimeformat)
	i := 1
	for ; ; i++ {
		if _, err := os.Stat(fmt.Sprintf(Fileformat, now, i)); os.IsNotExist(err) {
			break
		}
	}
	file, err := os.OpenFile(fmt.Sprintf(Fileformat, now, i), os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_TRUNC|os.O_APPEND, 0700)
	if err != nil {
		panic(err)
	}
	if file == nil {
		panic(os.ErrInvalid)
	}
	writer = bufio.NewWriter(file)
}

// Printf ...
func Printf(message string, args ...interface{}) {
	Infof(message, args...)
}

// Println ...
func Println(args ...interface{}) {
	Infoln(args...)
}

// Infof ...
func Infof(message string, args ...interface{}) {
	logln(Info, fmt.Sprintf(message, args...))
}

// Warnf ...
func Warnf(message string, args ...interface{}) {
	logln(Warn, fmt.Sprintf(message, args...))
}

// Errorf ...
func Errorf(message string, args ...interface{}) {
	logln(Error, fmt.Sprintf(message, args...))
}

// Fatalf ...
func Fatalf(message string, args ...interface{}) {
	logln(Fatal, fmt.Sprintf(message, args...))
}

// Logf formats and logs a message.
func Logf(level Level, message string, args ...interface{}) {
	logln(level, fmt.Sprintf(message, args...))
}

// Infoln ...
func Infoln(args ...interface{}) {
	logln(Info, fmt.Sprintln(args...))
}

// Warnln ...
func Warnln(args ...interface{}) {
	logln(Warn, fmt.Sprintln(args...))
}

// Errorln ...
func Errorln(args ...interface{}) {
	logln(Error, fmt.Sprintln(args...))
}

// Fatalln ...
func Fatalln(args ...interface{}) {
	logln(Fatal, fmt.Sprintln(args...))
}

// Logln logs a message.
func Logln(level Level, args ...interface{}) {
	logln(level, fmt.Sprintln(args...))
}

func logln(level Level, message string) {
	msg := []byte(fmt.Sprintf("[%[1]s] [%[2]s] %[3]s\n", time.Now().Format(Timeformat), level.ToString(), message))

	_, err := writer.Write(msg)
	if err != nil {
		panic(err)
	}
	lines++
	if lines == 5 {
		lines = 0
		writer.Flush()
	}

	if level >= Error {
		os.Stderr.Write(level.GetColor())
		os.Stderr.Write(msg)
		os.Stderr.Write([]byte("\x1b[0m"))
	} else {
		os.Stdout.Write(level.GetColor())
		os.Stdout.Write(msg)
		os.Stdout.Write([]byte("\x1b[0m"))
	}
}

// Shutdown cleans up the logger.
func Shutdown() {
	writer.Flush()
}