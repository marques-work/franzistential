package logging

import (
	"fmt"
	"io"
	"os"
)

var (
	Trace  *bool
	Silent *bool
)

func Debug(message string, args ...interface{}) {
	if *Trace {
		_print(os.Stderr, message, args...)
	}
}

func Warn(message string, args ...interface{}) {
	_print(os.Stderr, "[WARN] "+message, args...)
}

func Error(message string, args ...interface{}) {
	_print(os.Stderr, "[ERROR] "+message, args...)
}

func Die(message string, args ...interface{}) {
	*Silent = false // Always print out fatal errors before exiting
	Error(message, args...)
	os.Exit(1)
}

func _print(w io.Writer, message string, args ...interface{}) {
	if *Silent {
		return
	}

	if 0 == len(args) {
		fmt.Fprintln(w, message)
	} else {
		fmt.Fprintf(w, message+"\n", args...)
	}
}
