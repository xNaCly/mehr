// logging abstraction with prefixes and predefined logs
package log

import (
	"fmt"
	"log"
)

const (
	ansi_reset   = "\033[0m"
	ansi_red     = "\033[91m"
	ansi_yellow  = "\033[93m"
	ansi_blue    = "\033[94m"
	ansi_magenta = "\033[95m"
)

const (
	errPrefix  = ansi_red + "err: " + ansi_reset
	infoPrefix = ansi_blue + "info: " + ansi_reset
	warnPrefix = ansi_yellow + "warn: " + ansi_reset
)

func printf(prefix string, format string, v ...any) {
	log.Printf(prefix+format, v...)
}

func Infof(format string, v ...any) {
	printf(infoPrefix, format, v...)
}

func Warnf(format string, v ...any) {
	printf(warnPrefix, format, v...)
}

func Warn(v ...any) {
	fmt.Print(warnPrefix)
	log.Println(v...)
}

func Errorf(format string, v ...any) {
	printf(errPrefix, format, v...)
}

func Error(v ...any) {
	fmt.Print(errPrefix)
	log.Println(v...)
}
