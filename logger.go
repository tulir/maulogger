package maulog

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// Level is the severity level of a log entry.
type Level struct {
	Name            string
	Severity, Color int
}

// GetColor gets the ANSI escape color code for the log level.
func (lvl Level) GetColor() []byte {
	if lvl.Color == -1 {
		return []byte{}
	}
	return []byte(fmt.Sprintf("\x1b[3%dm", lvl.Color))
}

var (
	// Debug is the level for debug messages.
	Debug = Level{Name: "DEBUG", Color: 6, Severity: 0}
	// Info is the level for basic log messages.
	Info = Level{Name: "INFO", Color: -1, Severity: 10}
	// Warn is the level saying that something went wrong, but the program will continue operating mostly normally.
	Warn = Level{Name: "WARN", Color: 3, Severity: 50}
	// Error is the level saying that something went wrong and the program may not operate as expected, but will still continue.
	Error = Level{Name: "ERROR", Color: 1, Severity: 100}
	// Fatal is the level saying that something went wrong and the program will not operate normally.
	Fatal = Level{Name: "FATAL", Color: 5, Severity: 9001}
)

// PrintDebug tells if debug messages (severity lower than 10) should be printed.
var PrintDebug = false

// FileTimeformat is the time format used in log file names.
var FileTimeformat = "2006-01-02"

// FileformatArgs is an undocumented integer.
var FileformatArgs = 3

// Fileformat is the format used for log file names.
var Fileformat func(string, int) string

// Timeformat is the time format used in logging.
var Timeformat = "15:04:05 02.01.2006"

var writer *bufio.Writer
var lines int

// Init ...
func Init() {
	now := time.Now().Format(FileTimeformat)
	i := 1
	for ; ; i++ {
		if _, err := os.Stat(Fileformat(now, i)); os.IsNotExist(err) {
			break
		}
	}
	file, err := os.OpenFile(Fileformat(now, i), os.O_WRONLY|os.O_CREATE|os.O_EXCL|os.O_TRUNC|os.O_APPEND, 0700)
	if err != nil {
		panic(err)
	}
	if file == nil {
		panic(os.ErrInvalid)
	}
	writer = bufio.NewWriter(file)
}

// Debugf formats and logs a debug message.
func Debugf(message string, args ...interface{}) {
	logln(Debug, fmt.Sprintf(message, args...))
}

// Printf formats and logs a string in the Info log level.
func Printf(message string, args ...interface{}) {
	Infof(message, args...)
}

// Infof formats and logs a string in the Info log level.
func Infof(message string, args ...interface{}) {
	logln(Info, fmt.Sprintf(message, args...))
}

// Warnf formats and logs a string in the Warn log level.
func Warnf(message string, args ...interface{}) {
	logln(Warn, fmt.Sprintf(message, args...))
}

// Errorf formats and logs a string in the Error log level.
func Errorf(message string, args ...interface{}) {
	logln(Error, fmt.Sprintf(message, args...))
}

// Fatalf formats and logs a string in the Fatal log level.
func Fatalf(message string, args ...interface{}) {
	logln(Fatal, fmt.Sprintf(message, args...))
}

// Logf formats and logs a message in the given log level.
func Logf(level Level, message string, args ...interface{}) {
	logln(level, fmt.Sprintf(message, args...))
}

// Debugln logs a debug message.
func Debugln(message string, args ...interface{}) {
	logln(Debug, fmt.Sprintf(message, args...))
}

// Println logs a string in the Info log level.
func Println(args ...interface{}) {
	Infoln(args...)
}

// Infoln logs a string in the Info log level.
func Infoln(args ...interface{}) {
	logln(Info, fmt.Sprintln(args...))
}

// Warnln logs a string in the Warn log level.
func Warnln(args ...interface{}) {
	logln(Warn, fmt.Sprintln(args...))
}

// Errorln logs a string in the Error log level.
func Errorln(args ...interface{}) {
	logln(Error, fmt.Sprintln(args...))
}

// Fatalln logs a string in the Fatal log level.
func Fatalln(args ...interface{}) {
	logln(Fatal, fmt.Sprintln(args...))
}

// Logln logs a message in the given log level.
func Logln(level Level, args ...interface{}) {
	logln(level, fmt.Sprintln(args...))
}

func logln(level Level, message string) {
	msg := []byte(fmt.Sprintf("[%[1]s] [%[2]s] %[3]s\n", time.Now().Format(Timeformat), level.Name, message))

	_, err := writer.Write(msg)
	if err != nil {
		panic(err)
	}
	lines++
	if lines == 5 {
		lines = 0
		writer.Flush()
	}
	if level.Severity >= Error.Severity {
		os.Stderr.Write(level.GetColor())
		os.Stderr.Write(msg)
		os.Stderr.Write([]byte("\x1b[0m"))
	} else if level.Severity >= Info.Severity || PrintDebug {
		os.Stdout.Write(level.GetColor())
		os.Stdout.Write(msg)
		os.Stdout.Write([]byte("\x1b[0m"))
	}
}

// Shutdown cleans up the logger.
func Shutdown() {
	writer.Flush()
}
