package main

import (
	"fmt"

	"github.com/harsh3341/3rd-Semester-Mini-Project/api"
	"github.com/harsh3341/3rd-Semester-Mini-Project/bot"
	"github.com/harsh3341/3rd-Semester-Mini-Project/config"
)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	api.FetchTrivia()

	bot.Start()

	<-make(chan struct{})
	return
}
