package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	action, err := NewDownloader()
	if err != nil {
		log.Panic(err)
	}

	if err := action.Run(); err != nil {
		log.Panic(err)
	}
}