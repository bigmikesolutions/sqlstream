package containers

import "log"

type logger struct{}

func (c logger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
