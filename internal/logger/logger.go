package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logEntry *Logger

func GetLogger() Logger {
	return *logEntry
}

type Logger struct {
	*logrus.Entry
	lmbjk *lumberjack.Logger
}

func (entry *Logger) Close() {
	err := entry.lmbjk.Close()
	if err != nil {
		panic(err)
	}
}

func getLogger(lgr *logrus.Entry, lmbjk *lumberjack.Logger) *Logger {
	return &Logger{lgr, lmbjk}
}

type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err := w.Write([]byte(line))
		if err != nil {
			return err
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func init() {
	l := logrus.New()
	// TODO -> trace level from config
	l.SetLevel(logrus.TraceLevel)
	l.SetReportCaller(true)

	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: time.DateTime, // RFC3339
	}

	logsDir, exists := os.LookupEnv("LOCAL_LOGS_DIR")
	if !exists {
		logsDir = "logs"
	}
	err := os.MkdirAll(logsDir, 0750)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	//_, err = os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	//if err != nil {
	//	panic(err)
	//}

	// Create a new logger that rotates the log file
	lmbj := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/all.log", logsDir), // "/var/log/myapp.log",
		MaxSize:    1,                                  // Max size in megabytes
		MaxBackups: 3,                                  // Max number of old log files to keep
		MaxAge:     3,                                  // Max age in days
		Compress:   true,
	}

	l.SetOutput(lmbj)
	l.AddHook(&writerHook{
		Writer:    []io.Writer{os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	lgr := logrus.NewEntry(l)
	logEntry = getLogger(lgr, lmbj)
}

// TODO elasticSearch
