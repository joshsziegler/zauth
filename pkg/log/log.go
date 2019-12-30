package log

import (
	"fmt"
	"log"
)

func init() {
	// Don't print date, time, file name, or line number as a prefix when logging
	log.SetFlags(0)

	// Print the date, time, file name, and line number as a prefix when logging
	// log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func Info(args ...interface{}) {
	log.Print("I: " + fmt.Sprint(args...))
}

func Infof(format string, args ...interface{}) {
	log.Print("I: " + fmt.Sprintf(format, args...))
}

func Debug(args ...interface{}) {
	log.Print("D: " + fmt.Sprint(args...))
}

func Debugf(format string, args ...interface{}) {
	log.Print("D: " + fmt.Sprintf(format, args...))
}

func Error(args ...interface{}) {
	log.Print("E: " + fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	log.Print("E: " + fmt.Sprintf(format, args...))
}

func Fatal(args ...interface{}) {
	log.Fatal("F: " + fmt.Sprint(args...))
}

func Fatalf(format string, args ...interface{}) {
	log.Fatal("F: " + fmt.Sprintf(format, args...))
}
