package log

import (
	"fmt"
	"log"
)

var debug = false

func EnableDebug()  {
	debug = true
}

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

func Debug(l interface{})  {
	if debug {
		debug := fmt.Sprintf("[%s] %v", "DEBUG", l)
		log.Println(debug)
	}
}

func Debugf(format string, a ...interface{}) {
	if debug {
		tmp := fmt.Sprintf(format, a...)
		debug := fmt.Sprintf("[%s] %s", "Error", tmp)
		log.Println(debug)
	}
}