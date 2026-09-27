package logger

import "go.uber.org/zap"

type Logger struct {
    zap *zap.Logger
}

func New() *Logger {
    logger, err := zap.NewDevelopment()
    if err != nil {
        panic("logger creation error")
    }

    return &Logger{
        zap: logger,
    }
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
    l.zap.Info(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
    l.zap.Error(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
    l.zap.Warn(msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
    l.zap.Debug(msg, fields...)
}