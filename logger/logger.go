package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitializeLogger(level string, logpath string) {
	switch level {
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	case "file":
		fileLogger(logpath)
	}

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func fileLogger(logpath string) {
	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}
	Log.SetLevel(logrus.DebugLevel)
}

/*
  log.Trace("Something very low level.")
  log.Debug("Useful debugging information.")
  log.Info("Something noteworthy happened!")
  log.Warn("You should probably take a look at this.")
  Log.Fatal("Something failed but I'm not quitting.")
  // Calls os.Exit(1) after logging
  log.Fatal("Bye.")
  // Calls panic() after logging
  log.Panic("I'm bailing.")
*/
