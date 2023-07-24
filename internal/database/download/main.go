package main

import (
	log "github.com/sirupsen/logrus"
)

//go:generate go run .

func main() {
	action, err := NewDownloader("Nintendo - Nintendo Entertainment System")
	if err != nil {
		log.Panic(err)
	}

	if err := action.Run(); err != nil {
		log.Panic(err)
	}
}
