package main

import (
	"fmt"
	"net"
)

type Logger struct {
	ch   chan string
	conn net.Conn
}

func (l *Logger) msg(str string) {
	l.ch <- str
}
func (l *Logger) br() {
	l.ch <- "--------------------"
}

func (l *Logger) init() {
	conn, err := net.Dial("tcp", "localhost:7777")
	if err != nil {
		return
	}
	l.conn = conn
	defer l.conn.Close()

	for {
		msg := <-l.ch
		fmt.Fprintln(conn, msg)
	}
}
