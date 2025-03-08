package logger

import (
	"fmt"
	"log"
)

func Error(v ...any) {
	log.Println("[ERROR]", fmt.Sprint(v...))
}

func Warn(v ...any) {
	log.Println("[WARN]", fmt.Sprint(v...))
}

func Info(v ...any) {
	log.Println("[INFO]", fmt.Sprint(v...))
}
