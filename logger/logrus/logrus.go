package logrus

import (
	log "github.com/sirupsen/logrus"
)

type Logrus struct {
	*log.Logger
}

func NewLogrus(logger *log.Logger) *Logrus {
	return &Logrus{logger}
}

//func (l *Logrus) Tracef(format string, args ...interface{}) {
//	l.logger.Tracef(format, args...)
//}
//
//func (l *Logrus) Debugf(format string, args ...interface{}) {
//	l.logger.Debugf(format, args...)
//}
//
//func (l *Logrus) Infof(format string, args ...interface{}) {
//	l.logger.Infof(format, args...)
//}
//
//func (l *Logrus) Printf(format string, args ...interface{}) {
//	l.logger.Printf(format, args...)
//}
//
//func (l *Logrus) Warnf(format string, args ...interface{}) {
//	l.logger.Warnf(format, args...)
//}
//
//func (l *Logrus) Errorf(format string, args ...interface{}) {
//	l.logger.Errorf(format, args...)
//}
//
//func (l *Logrus) Fatalf(format string, args ...interface{}) {
//	l.logger.Fatalf(format, args...)
//}
//
//func (l *Logrus) Panicf(format string, args ...interface{}) {
//	l.logger.Panicf(format, args...)
//}
//
//func (l *Logrus) Trace(args ...interface{}) {
//	l.logger.Trace(args...)
//}
//
//func (l *Logrus) Debug(args ...interface{}) {
//	l.logger.Debug(args...)
//}
//
//func (l *Logrus) Info(args ...interface{}) {
//	l.logger.Info(args...)
//}
//
//func (l *Logrus) Print(args ...interface{}) {
//	l.logger.Print(args...)
//}
//
//func (l *Logrus) Warn(args ...interface{}) {
//	l.logger.Warn(args...)
//}
//
//func (l *Logrus) Error(args ...interface{}) {
//	l.logger.Error(args...)
//}
//
//func (l *Logrus) Fatal(args ...interface{}) {
//	l.logger.Fatal(args...)
//}
//
//func (l *Logrus) Panic(args ...interface{}) {
//	l.logger.Panic(args...)
//}
