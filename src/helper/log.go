package helper

import (
	"fmt"
	"time"
)

func Log(format string, args ...any) {
	fmt.Printf("%s %s\n", Cyan("GitUp:"), fmt.Sprintf(format, args...))
}

func LogFail(format string, args ...any) {
	fmt.Printf("%s %s\n", Red(Bold("✗ FAIL")), fmt.Sprintf(format, args...))
}

func LogOk(format string, args ...any) {
	fmt.Printf("%s %s\n", Green(Bold("✓")), fmt.Sprintf(format, args...))
}

func LogWarn(format string, args ...any) {
	fmt.Printf("%s %s\n", Yellow(Bold("⚠ WARN")), fmt.Sprintf(format, args...))
}

func LogError(format string, args ...any) {
	now := time.Now()
	timestamp := Dim(fmt.Sprintf("[%02d:%02d:%02d]", now.Hour(), now.Minute(), now.Second()))

	fmt.Printf("%s %s %s\n", timestamp, Red(Bold("ERROR")), fmt.Sprintf(format, args...))
}
