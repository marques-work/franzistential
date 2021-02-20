package logging

import (
	"fmt"
	"io"
	"os"
	"time"
)

var (
	// Trace allows debug messages
	Trace bool

	// Silent prevents printing output
	Silent bool

	// File specifies the log output; defaults to STDERR
	File io.Writer = os.Stderr
)

// Debug logs TRACE level; ignores Silent
func Debug(message string, args ...interface{}) {
	if Trace {
		_forcePrint(File, "[DEBUG] "+message, args...)
	}
}

// Warn logs WARN level
func Warn(message string, args ...interface{}) {
	_print(File, "[WARN] "+message, args...)
}

// WarnMercilessly logs WARN level; ignores Silent
func WarnMercilessly(message string, args ...interface{}) {
	_forcePrint(File, message, args...)
}

// Error logs ERROR level
func Error(message string, args ...interface{}) {
	_print(File, "[ERROR] "+message, args...)
}

// Die logs FATAL level and exits; ignores Silent
func Die(message string, args ...interface{}) {
	_forcePrint(File, "[FATAL] "+message, args...)
	os.Exit(1)
}

func _print(w io.Writer, message string, args ...interface{}) {
	if Silent {
		return
	}

	_forcePrint(w, message, args...)
}

func _forcePrint(w io.Writer, message string, args ...interface{}) {
	message = time.Now().Format(time.RFC3339) + " " + message

	if 0 == len(args) {
		fmt.Fprintln(w, message)
	} else {
		fmt.Fprintf(w, message+"\n", args...)
	}
}
