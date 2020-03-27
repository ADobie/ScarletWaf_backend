package tool

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

// singleton ,for the logger is concurrent safe
var (
	logger = logrus.New()
	once   sync.Once
)

func GetLogger() *logrus.Logger {
	once.Do(func() {
		logger.SetFormatter(&logrus.JSONFormatter{})
		logger.SetOutput(os.Stdout)
	})
	return logger
}
