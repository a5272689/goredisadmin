package controllers

import (
	"os"
	"log"
)

func NewLogger() *log.Logger {
	logger:=log.New(os.Stdout,"[goredisadmin] ",log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	return logger
}

var Logger=NewLogger()
