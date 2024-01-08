package logger

//
//import (
//	"fmt"
//	"github.com/sirupsen/logrus"
//)
//
//func Trace(args ...interface{}) {
//	logEntry.Log(logrus.TraceLevel, args...)
//}
//
//func Debug(args ...interface{}) {
//	logEntry.Log(logrus.DebugLevel, args...)
//}
//
//func Print(args ...interface{}) {
//	logEntry.Info(args...)
//}
//
//func Info(args ...interface{}) {
//	logEntry.Log(logrus.InfoLevel, args...)
//}
//
//func Warn(args ...interface{}) {
//	logEntry.Log(logrus.WarnLevel, args...)
//}
//
//func Warning(args ...interface{}) {
//	logEntry.Warn(args...)
//}
//
//func Error(args ...interface{}) {
//	logEntry.Log(logrus.ErrorLevel, args...)
//}
//
//func Fatal(args ...interface{}) {
//	logEntry.Log(logrus.FatalLevel, args...)
//	logEntry.Logger.Exit(1)
//}
//
//func Panic(args ...interface{}) {
//	logEntry.Log(logrus.PanicLevel, args...)
//}
//
//// Entry Printf family functions
//
//func Logf(level logrus.Level, format string, args ...interface{}) {
//	if logEntry.Logger.IsLevelEnabled(level) {
//		logEntry.Log(level, fmt.Sprintf(format, args...))
//	}
//}
//
//func Tracef(format string, args ...interface{}) {
//	logEntry.Logf(logrus.TraceLevel, format, args...)
//}
//
//func Debugf(format string, args ...interface{}) {
//	logEntry.Logf(logrus.DebugLevel, format, args...)
//}
//
//func Infof(format string, args ...interface{}) {
//	logEntry.Logf(logrus.InfoLevel, format, args...)
//}
//
//func Printf(format string, args ...interface{}) {
//	logEntry.Infof(format, args...)
//}
//
//func Warnf(format string, args ...interface{}) {
//	logEntry.Logf(logrus.WarnLevel, format, args...)
//}
//
//func Warningf(format string, args ...interface{}) {
//	logEntry.Warnf(format, args...)
//}
//
//func Errorf(format string, args ...interface{}) {
//	logEntry.Logf(logrus.ErrorLevel, format, args...)
//}
//
//func Fatalf(format string, args ...interface{}) {
//	logEntry.Logf(logrus.FatalLevel, format, args...)
//	logEntry.Logger.Exit(1)
//}
//
//func Panicf(format string, args ...interface{}) {
//	logEntry.Logf(logrus.PanicLevel, format, args...)
//}
//
//// Entry Println family functions
//
//func Logln(level logrus.Level, args ...interface{}) {
//	if logEntry.Logger.IsLevelEnabled(level) {
//		logEntry.Log(level, sprintlnn(args...))
//	}
//}
//
//func Traceln(args ...interface{}) {
//	logEntry.Logln(logrus.TraceLevel, args...)
//}
//
//func Debugln(args ...interface{}) {
//	logEntry.Logln(logrus.DebugLevel, args...)
//}
//
//func Infoln(args ...interface{}) {
//	logEntry.Logln(logrus.InfoLevel, args...)
//}
//
//func Println(args ...interface{}) {
//	logEntry.Infoln(args...)
//}
//
//func Warnln(args ...interface{}) {
//	logEntry.Logln(logrus.WarnLevel, args...)
//}
//
//func Warningln(args ...interface{}) {
//	logEntry.Warnln(args...)
//}
//
//func Errorln(args ...interface{}) {
//	logEntry.Logln(logrus.ErrorLevel, args...)
//}
//
//func Fatalln(args ...interface{}) {
//	logEntry.Logln(logrus.FatalLevel, args...)
//	logEntry.Logger.Exit(1)
//}
//
//func Panicln(args ...interface{}) {
//	logEntry.Logln(logrus.PanicLevel, args...)
//}
//
//// Sprintlnn => Sprint no newline. This is to get the behavior of how
//// fmt.Sprintln where spaces are always added between operands, regardless of
//// their type. Instead of vendoring the Sprintln implementation to spare a
//// string allocation, we do the simplest thing.
//func sprintlnn(args ...interface{}) string {
//	msg := fmt.Sprintln(args...)
//	return msg[:len(msg)-1]
//}
