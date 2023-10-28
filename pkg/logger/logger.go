package logger

type Logger interface {
	Warn(kv ...interface{})
	Error(kv ...interface{})
	Debug(kv ...interface{})
	Info(kv ...interface{})

	Warnf(s string, kv ...interface{})
	Errorf(s string, kv ...interface{})
	Debugf(s string, kv ...interface{})
	Infof(s string, kv ...interface{})
}

func WithPrefix(log Logger, kv ...string) Logger {
	entry := log.(*LogrusWrapper).log

	if len(kv) == 2 {
		entry = entry.WithField(kv[0], kv[1])
	}

	return &LogrusWrapper{
		log: entry,
	}
}

func WithPrefixMap(log Logger, f map[string]interface{}) Logger {
	entry := log.(*LogrusWrapper).log
	entry = entry.WithFields(f)

	return &LogrusWrapper{
		log: entry,
	}
}
