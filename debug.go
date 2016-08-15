package main

import (
	"flag"
	"fmt"
	"log"
)

var debug debugT

type debugT struct {
	enabled bool
}

func (d debugT) Printf(format string, args ...interface{}) {
	if d.enabled {
		msg := fmt.Sprintf(format, args...)
		log.Output(2, msg)
	}
}

func (d debugT) Println(args ...interface{}) {
	if d.enabled {
		msg := fmt.Sprintln(args...)
		log.Output(2, msg)
	}
}

func init() {
	flag.BoolVar(&debug.enabled, "debug", false, "show debug")
}
