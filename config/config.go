package config

import (
	"fmt" // for printing
	"os"  // for reading environment variables

	"github.com/joho/godotenv" // for reading .env file
)

// Token is the bot token.
var (
	Token     string
	BotPrefix string
	APITOKEN  string
)

// ReadConfig reads the .env file and sets the bot token and prefix.
func ReadConfig() error {

	// load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = os.Getenv("TOKEN")
	BotPrefix = os.Getenv("BOTPREFIX")
	APITOKEN = os.Getenv("APITOKEN")

	return nil

}
