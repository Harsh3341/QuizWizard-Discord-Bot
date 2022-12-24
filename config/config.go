package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Token     string
	BotPrefix string
)

func ReadConfig() error {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = os.Getenv("TOKEN")
	BotPrefix = os.Getenv("BOTPREFIX")

	return nil

}
