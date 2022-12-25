package main

import (
	"fmt"

	"github.com/harsh3341/3rd-Semester-Mini-Project/bot"
	"github.com/harsh3341/3rd-Semester-Mini-Project/config"
)

func main() {

	// read the config file
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// start the bot
	bot.Start()

	// wait for the bot to exit
	<-make(chan struct{})
	return
}
