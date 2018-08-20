package log

import (
	"fmt"
	"log"
)

func Info(l interface{}) {
	info := fmt.Sprintf("[%s] %v", "INFO", l)
	log.Println(info)
}

func Infof(l string, a ...interface{}) {
	tmp := fmt.Sprintf(l, a...)
	info := fmt.Sprintf("[%s] %s", "INFO", tmp)
	log.Println(info)
}

func Error(l interface{}) {
	err := fmt.Sprintf("[%s] %v", "Error", l)
	log.Println(err)
}

func Errorf(format string, a ...interface{}) {
	tmp := fmt.Sprintf(format, a...)
	err := fmt.Sprintf("[%s] %s", "Error", tmp)
	log.Println(err)
}