package logger

import "go.uber.org/zap"

var logger *zap.Logger

func Init(dev bool) error {
	l, err := newLogger(dev)
	if err != nil {
		return err
	}
	logger = l
	return nil
}

func newLogger(dev bool) (*zap.Logger, error) {
	var err error
	var l *zap.Logger
	if dev {
		l, err = zap.NewDevelopment()
	} else {
		l, err = zap.NewProduction()
	}
	if err != nil {
		return nil, err
	}
	return l, nil
}

func Get() *zap.Logger {
	return logger
}

func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}
