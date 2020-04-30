// mauLogger - A logger for Go programs
// Copyright (c) 2020 Tulir Asokan
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package maulogger

// DefaultLogger ...
var DefaultLogger = Create().(*BasicLogger)

// OpenFile formats the given parts with fmt.Sprint and logs the result with the OpenFile level
func Open() error {
	return DefaultLogger.Open()
}

// Close formats the given parts with fmt.Sprint and logs the result with the Close level
func Close() error {
	return DefaultLogger.Close()
}

// Sub creates a Sublogger
func Sub(module ...string) Logger {
	return DefaultLogger.Sub(module...)
}

// Raw formats the given parts with fmt.Sprint and logs the result with the Raw level
func Raw(level Level, module, message string) {
	DefaultLogger.Raw(level, module, message)
}

// Log formats the given parts with fmt.Sprint and logs the result with the given level
func Log(level Level, parts ...interface{}) {
	DefaultLogger.Log(level, parts...)
}

// Logf formats the given message and args with fmt.Sprintf and logs the result with the given level
func Logf(level Level, message string, args ...interface{}) {
	DefaultLogger.Logf(level, message, args...)
}

// Debug formats the given parts with fmt.Sprint and logs the result with the Debug level
func Debug(parts ...interface{}) {
	DefaultLogger.Debug(parts...)
}

// Debugf formats the given message and args with fmt.Sprintf and logs the result with the Debug level
func Debugf(message string, args ...interface{}) {
	DefaultLogger.Debugf(message, args...)
}

// Info formats the given parts with fmt.Sprint and logs the result with the Info level
func Info(parts ...interface{}) {
	DefaultLogger.Info(parts...)
}

// Infof formats the given message and args with fmt.Sprintf and logs the result with the Info level
func Infof(message string, args ...interface{}) {
	DefaultLogger.Infof(message, args...)
}

// Warn formats the given parts with fmt.Sprint and logs the result with the Warn level
func Warn(parts ...interface{}) {
	DefaultLogger.Warn(parts...)
}

// Warnf formats the given message and args with fmt.Sprintf and logs the result with the Warn level
func Warnf(message string, args ...interface{}) {
	DefaultLogger.Warnf(message, args...)
}

// Error formats the given parts with fmt.Sprint and logs the result with the Error level
func Error(parts ...interface{}) {
	DefaultLogger.Error(parts...)
}

// Errorf formats the given message and args with fmt.Sprintf and logs the result with the Error level
func Errorf(message string, args ...interface{}) {
	DefaultLogger.Errorf(message, args...)
}

// Fatal formats the given parts with fmt.Sprint and logs the result with the Fatal level
func Fatal(parts ...interface{}) {
	DefaultLogger.Fatal(parts...)
}

// Fatalf formats the given message and args with fmt.Sprintf and logs the result with the Fatal level
func Fatalf(message string, args ...interface{}) {
	DefaultLogger.Fatalf(message, args...)
}
