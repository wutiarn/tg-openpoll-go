package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
)

func main() {

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	token := flag.String("token", "", "Telegram token")
	flag.Parse()

	if *token == "" {
		logrus.Fatalln("No -token flag provided")
	}
}
