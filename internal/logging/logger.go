// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/

package logging

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger
type Logger struct {
	zerolog.Logger
}

// Config holds logger configuration
type Config struct {
	Level      string
	Pretty     bool
	TimeFormat string
	Output     io.Writer
}

// DefaultConfig returns the default logger configuration
func DefaultConfig() Config {
	return Config{
		Level:      "info",
		Pretty:     true,
		TimeFormat: time.RFC3339,
		Output:     os.Stdout,
	}
}

// NewLogger creates a new Logger with the given configuration
func NewLogger(cfg Config) *Logger {
	// Set time format
	zerolog.TimeFieldFormat = cfg.TimeFormat

	// Parse level
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	var output io.Writer = cfg.Output
	if output == nil {
		output = os.Stdout
	}

	// Use console writer for pretty output
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        output,
			TimeFormat: cfg.TimeFormat,
		}
	}

	logger := zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{Logger: logger}
}

// NewDefaultLogger creates a logger with default configuration
func NewDefaultLogger() *Logger {
	return NewLogger(DefaultConfig())
}

// WithComponent adds a component field to the logger
func (l *Logger) WithComponent(component string) *Logger {
	return &Logger{
		Logger: l.Logger.With().Str("component", component).Logger(),
	}
}

// WithVM adds a VM ID field to the logger
func (l *Logger) WithVM(vmID string) *Logger {
	return &Logger{
		Logger: l.Logger.With().Str("vm_id", vmID).Logger(),
	}
}

// WithUser adds a user ID field to the logger
func (l *Logger) WithUser(userID string) *Logger {
	return &Logger{
		Logger: l.Logger.With().Str("user_id", userID).Logger(),
	}
}

// WithRequestID adds a request ID field to the logger
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		Logger: l.Logger.With().Str("request_id", requestID).Logger(),
	}
}

// Global logger instance
var globalLogger *Logger

// Init initializes the global logger
func Init(cfg Config) {
	globalLogger = NewLogger(cfg)
}

// GetLogger returns the global logger
func GetLogger() *Logger {
	if globalLogger == nil {
		globalLogger = NewDefaultLogger()
	}
	return globalLogger
}

// Convenience functions for global logger
func Debug() *zerolog.Event { return GetLogger().Debug() }
func Info() *zerolog.Event  { return GetLogger().Info() }
func Warn() *zerolog.Event  { return GetLogger().Warn() }
func Error() *zerolog.Event { return GetLogger().Error() }
func Fatal() *zerolog.Event { return GetLogger().Fatal() }
