package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/wutiarn/tg-openpoll-go"
)

func main() {
	token := flag.String("token", "", "Telegram token")
	debug := flag.Bool("debug", false, "Debug mode")
	flag.Parse()

	if *token == "" {
		logrus.Fatalln("No -token flag provided")
	}

	if *debug == true {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	}

	openpoll.Run(*token, *debug)
}
