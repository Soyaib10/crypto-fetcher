package logger

import "fmt"

func LogInfo(format string, a ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", a...)
}

func LogError(format string, a ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", a...)
}
