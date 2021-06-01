package logger

import "github.com/sirupsen/logrus"

func Log(err error) {
	logrus.WithField("error", err)
}
func Info(info string) {
	logrus.Info(info)
}
