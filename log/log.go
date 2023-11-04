// logging abstraction with prefixes and predefined logs
package log

import (
	"log"
	"strings"
)

const (
	ansi_reset   = "\033[0m"
	ansi_red     = "\033[91m"
	ansi_yellow  = "\033[93m"
	ansi_blue    = "\033[94m"
	ansi_magenta = "\033[95m"
)

func print(prefix []string, format string, v ...any) {
	log.Printf(strings.Join(prefix, "")+format, v...)
}

func Info(format string, v ...any) {
	print([]string{ansi_blue, "info: ", ansi_reset}, format, v...)
}

func Warn(format string, v ...any) {
	print([]string{ansi_yellow, "warn: ", ansi_reset}, format, v...)
}

func Error(format string, v ...any) {
	print([]string{ansi_red, "err: ", ansi_reset}, format, v...)
}
