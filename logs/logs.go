package logs

import (
	"fmt"
	"log"
	"os"
)

type PackageName string

var (
	infoLogger  = log.New(os.Stdout, "Info - ", log.LstdFlags|log.Lmicroseconds)
	warnLogger  = log.New(os.Stdout, "Warning - ", log.LstdFlags|log.Lmicroseconds)
	errorLogger = log.New(os.Stdout, "Error - ", log.LstdFlags|log.Lmicroseconds)
)

// Info prints information logs in the console
func Info(pkgTitle PackageName, message string) {
	infoMessage := fmt.Sprintf("[%s] %s", pkgTitle, message)

	infoLogger.Println(infoMessage)
}

// Infof prints formated information logs in the console
func Infof(pkgTitle PackageName, message string, v ...any) {
	Info(pkgTitle, fmt.Sprintf(message, v...))
}

// Warn prints warning logs in the console
func Warn(pkgTitle PackageName, message string) {
	warnMessage := fmt.Sprintf("[%s] %s", pkgTitle, message)

	warnLogger.Println(warnMessage)
}

// Warnf prints formated warning messages in the console
func Warnf(pkgTitle PackageName, message string, v ...any) {
	Warn(pkgTitle, fmt.Sprintf(message, v...))
}

// Error prints error logs in the console
func Error(pkgTitle PackageName, message string) {
	errorMessage := fmt.Sprintf("[%s] %s", pkgTitle, message)

	errorLogger.Println(errorMessage)
}

// Errorf prints formated error logs in the console
func Errorf(pkgTitle PackageName, message string, v ...any) {
	Error(pkgTitle, fmt.Sprintf(message, v...))
}
