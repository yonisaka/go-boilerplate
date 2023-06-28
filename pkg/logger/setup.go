package logger

import (
	"github.com/sirupsen/logrus"
)

func SetJSONFormatter() {
	logrus.SetFormatter(&Formatter{
		ChildFormatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg: MessageKey,
			},
		},
		Line:         true,
		Package:      false,
		File:         true,
		BaseNameOnly: false,
	})
}

func Setup(cfg Config) {
	conf = &cfg

	if cfg.Debug {
		logrus.SetLevel(logrus.TraceLevel)
		return
	}

	lvl, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
}
