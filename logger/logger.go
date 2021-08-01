package logger

import (
	"context"
	"time"

	"github.com/fatih/color"
	"gorm.io/gorm/logger"
)

type Logger interface {
	Tracef(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

type DBLogger struct {
	l Logger
}

func (dl *DBLogger) LogMode(logger.LogLevel) logger.Interface {
	return dl
}

func (dl *DBLogger) Info(ctx context.Context, msg string, i ...interface{}) {
	dl.l.Infof(msg, i...)
}
func (dl *DBLogger) Warn(ctx context.Context, msg string, i ...interface{}) {
	dl.l.Warnf(msg, i...)
}

func (dl *DBLogger) Error(ctx context.Context, msg string, i ...interface{}) {
	dl.l.Errorf(msg, i...)
}

func (dl *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	elapsedStr := color.CyanString("[%fms] ", float64(elapsed.Nanoseconds())/1e6)
	rowsStr := color.YellowString("[rows:%d] ", rows)
	if err != nil {
		sql = sql + color.RedString(" [%s]", err.Error())
	}
	if rows == -1 {
		dl.l.Trace(sql)
	} else {
		dl.l.Trace(elapsedStr, rowsStr, sql)
	}
}

func NewDBLogger(l Logger) *DBLogger {
	return &DBLogger{l}
}
