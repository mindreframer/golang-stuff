package log

import (
	"fmt"
	"os"
	"time"
)

func Fatal(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "%v goltunnel: %s\n", time.Now(), fmt.Sprintf(s, a...))
	os.Exit(2)
}

func Log(msg string, r ...interface{}) {
	fmt.Printf("%v - %s\n", time.Now(), fmt.Sprintf(msg, r...))
}

func Info(msg string, r ...interface{}) {
	fmt.Printf("\033[1;34m%s\033[0m\n", fmt.Sprintf(msg, r...))
}
