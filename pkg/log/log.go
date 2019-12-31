// Package log provides four levels of logging: Info, Debug, Error, and Fatal
// along with the format versions (e.g. Info() and Infof()). These print the log
// level in front of each message as a single character. This keeps the log
// messages shorter, but still easy to skim for errors or debug messages.
//
// The logger will add a the current date and time when running in development
// or testing environments. In production, only the message type prefix and
// message is printed. We assume the system's logger will add or store the date
// and time to each message (e.g. like syslog and systemd do).
//
// Example 1:
//    log.Init(environment.Dev)
//    log.Info("Hello")
//    2019/12/30 16:14:10 I: Hello
//
// Example 2:
//    log.Init(environment.Prod)
//    log.Info("Hello")
//    I: Hello
//
// Example 3:
//     log.Init(environment.Prod)
//     log.Info("Showing other message prefixes")
//     I: Showing other message prefixes
//     log.Debug("Showing other message prefixes")
//     [ Debug messages are not printed in production! ]
//     log.Error("Showing other message prefixes")
//     E: Showing other message prefixes
//
// Example 4:
//     log.Inif(environment.Prod)
//     log.Infof("Using formatting for variable Foo: %+v", foo)
//     I: Using formatting for variable Foo: FooStruct{Bar: 5, Valid: True}
//

package log

import (
	"fmt"
	"log"

	e "github.com/joshsziegler/zauth/pkg/environment"
)

// Init sets the prefix depending on the environment (dev, test, or prod).
func Init(environment int) {
	env = environment
	if env == e.Prod {
		log.SetFlags(0)
	}
}

var env int

// Info logs a normal, informative message if not in testing. For information
// that should not be logged in production, used Debug() instead.
func Info(args ...interface{}) {
	if env != e.Test {
		log.Print("I: " + fmt.Sprint(args...))
	}
}

func Infof(format string, args ...interface{}) {
	if env != e.Test {
		log.Print("I: " + fmt.Sprintf(format, args...))
	}
}

// Debug logs a message that should only be shown only in development.
func Debug(args ...interface{}) {
	if env == e.Dev {
		log.Print("D: " + fmt.Sprint(args...))
	}
}

func Debugf(format string, args ...interface{}) {
	if env != e.Dev {
		log.Print("D: " + fmt.Sprintf(format, args...))
	}
}

// Error logs a message that is a recoverable error. For unrecoverable errors,
// use Fatal() instead.
func Error(args ...interface{}) {
	log.Print("E: " + fmt.Sprint(args...))
}

func Errorf(format string, args ...interface{}) {
	log.Print("E: " + fmt.Sprintf(format, args...))
}

// Fatal logs an unrecoverable error. For recoverable errors, use Error()
// instead.
func Fatal(args ...interface{}) {
	log.Fatal("F: " + fmt.Sprint(args...))
}

func Fatalf(format string, args ...interface{}) {
	log.Fatal("F: " + fmt.Sprintf(format, args...))
}
