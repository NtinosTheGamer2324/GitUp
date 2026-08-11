package helper

import (
	"fmt"
	"time"
)

func Log(format string, args ...any) {
	fmt.Printf("GitUp: "+format+"\n", args...)
}

func LogFail(format string, args ...any) {
	fmt.Printf("[\033[31mFAIL\033[0m] GitUp: "+format+"\n", args...)
}

func LogOk(format string, args ...any) {
	fmt.Printf("[\033[32mOK\033[0m] GitUp: "+format+"\n", args...)
}

func LogError(format string, args ...any) {
	timestamp := time.Now()

	args = append(
		[]any{
			timestamp.Hour(),
			timestamp.Minute(),
			timestamp.Second(),
		},
		args...,
	)

	fmt.Printf(
		"[%02d:%02d:%02d \033[31mERROR\033[0m] "+format+"\n",
		args...,
	)
}
