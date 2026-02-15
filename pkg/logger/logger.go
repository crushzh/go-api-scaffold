package logger

import (
	"os"

	"go-api-scaffold/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var globalLogger *zap.SugaredLogger

// Init initializes the global logger
func Init(cfg *config.LogConfig) error {
	l, err := New(cfg)
	if err != nil {
		return err
	}
	globalLogger = l
	return nil
}

// New creates a new logger instance
func New(cfg *config.LogConfig) (*zap.SugaredLogger, error) {
	level := parseLevel(cfg.Level)

	// Encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	// Console output
	if cfg.Output == "console" || cfg.Output == "both" {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	}

	// File output
	if (cfg.Output == "file" || cfg.Output == "both") && cfg.FilePath != "" {
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, level))
	}

	if len(cores) == 0 {
		// Default to console output
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	}

	core := zapcore.NewTee(cores...)
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return zapLogger.Sugar(), nil
}

// L returns the global logger
func L() *zap.SugaredLogger {
	if globalLogger == nil {
		// Use default logger when not initialized
		l, _ := zap.NewDevelopment()
		return l.Sugar()
	}
	return globalLogger
}

// Sync flushes the logger buffer
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

// Convenience methods
func Debug(args ...interface{})                    { L().Debug(args...) }
func Debugf(template string, args ...interface{})  { L().Debugf(template, args...) }
func Info(args ...interface{})                     { L().Info(args...) }
func Infof(template string, args ...interface{})   { L().Infof(template, args...) }
func Warn(args ...interface{})                     { L().Warn(args...) }
func Warnf(template string, args ...interface{})   { L().Warnf(template, args...) }
func Error(args ...interface{})                    { L().Error(args...) }
func Errorf(template string, args ...interface{})  { L().Errorf(template, args...) }
func Fatal(args ...interface{})                    { L().Fatal(args...) }
func Fatalf(template string, args ...interface{})  { L().Fatalf(template, args...) }

func Debugw(msg string, keysAndValues ...interface{}) { L().Debugw(msg, keysAndValues...) }
func Infow(msg string, keysAndValues ...interface{})  { L().Infow(msg, keysAndValues...) }
func Warnw(msg string, keysAndValues ...interface{})   { L().Warnw(msg, keysAndValues...) }
func Errorw(msg string, keysAndValues ...interface{})  { L().Errorw(msg, keysAndValues...) }

func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
