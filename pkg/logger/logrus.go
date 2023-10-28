package logger

import (
	"io"
	log2 "log"

	"github.com/sirupsen/logrus"
)

var _ Logger = (*LogrusWrapper)(nil)

type LogrusWrapper struct {
	log *logrus.Entry
}

func New(writer io.Writer) *LogrusWrapper {
	log := logrus.New()

	log.SetOutput(writer)
	log.SetNoLock()

	log.SetLevel(logrus.DebugLevel)

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
	})

	log2.SetPrefix("std: ")
	log2.SetOutput(log.Writer())

	return &LogrusWrapper{
		log: logrus.NewEntry(log),
	}
}

func (l *LogrusWrapper) Warn(kv ...interface{}) {
	l.log.Warnln(kv...)
}

func (l *LogrusWrapper) Error(kv ...interface{}) {
	logWithFileLine := l.log.WithFields(logrus.Fields{})

	logWithFileLine.Errorln(kv...)
}

func (l *LogrusWrapper) Debug(kv ...interface{}) {
	l.log.Debugln(kv...)
}

func (l *LogrusWrapper) Info(kv ...interface{}) {
	l.log.Infoln(kv...)
}

func (l *LogrusWrapper) Warnf(s string, kv ...interface{}) {
	l.log.Warnf(s, kv...)
}

func (l *LogrusWrapper) Errorf(s string, kv ...interface{}) {
	l.log.Errorf(s, kv...)
}

func (l *LogrusWrapper) Debugf(s string, kv ...interface{}) {
	l.log.Debugf(s, kv...)
}

func (l *LogrusWrapper) Infof(s string, kv ...interface{}) {
	l.log.Infof(s, kv...)
}
