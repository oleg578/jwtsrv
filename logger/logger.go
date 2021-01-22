package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	lstd *log.Logger
)

func Fatal(v ...interface{}) {
	_ = lstd.Output(2, fmt.Sprint(v...))
	os.Exit(1)
}

func Fatalf(format string, v ...interface{}) {
	_ = lstd.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Print(v ...interface{}) {
	_ = lstd.Output(2, fmt.Sprint(v...))
}

func Printf(format string, v ...interface{}) {
	_ = lstd.Output(2, fmt.Sprintf(format, v...))
}

func GetLogger() *log.Logger {
	return lstd
}

func Start(logPath, prefix string) error {
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	lstd = log.New(f, prefix, log.LstdFlags|log.Lshortfile)
	lstd.SetPrefix(prefix)
	return nil
}
