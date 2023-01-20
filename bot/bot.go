package bot

import (
	"fmt"       // for printing
	"math/rand" // for shuffling options
	"strconv"   // for converting int to string
	"strings"   // for splitting strings

	"github.com/bwmarrin/discordgo"                         // for Discord bot
	"github.com/harsh3341/3rd-Semester-Mini-Project/api"    // for trivia questions
	"github.com/harsh3341/3rd-Semester-Mini-Project/config" // for bot token and prefix
)

// BotID is the ID of the bot.
var (
	BotID           string
	goBot           *discordgo.Session
	totalQuestion   int // number of questions to ask
	currentQuestion int
	score           int
)

// starts the trivia game.
func Start() {

	// create a new Discord session using the provided bot token.
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// get the account information for the bot.
	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	// add a handler for messages.
	goBot.AddHandler(messageHandler)

	// open a websocket connection to Discord and begin listening.
	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	api.FetchQuiz(1, "linux")

	fmt.Println("Bot is running...")

}

// messageHandler handles messages sent to the bot.
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages sent by the bot itself
	if m.Author.ID == BotID {
		return
	}

	// check if the message is a command
	if !strings.HasPrefix(m.Content, config.BotPrefix) {
		return
	}

	// split message into command and arguments
	args := strings.Split(m.Content, " ")

	// handle commands
	switch args[0] {
	case config.BotPrefix + "start":
		if len(args) != 3 {
			s.ChannelMessageSend(m.ChannelID, "```Usage: !start [number of questions] [Difficulty]```")
			return
		}

		numQuestions, err := strconv.Atoi(args[1]) // convert string to int
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "```Error: Invalid number of questions```")
			return
		}

		difficulty := args[2]

		if numQuestions < 1 || numQuestions > 50 {
			s.ChannelMessageSend(m.ChannelID, "```Error: Number of questions must be between 1 and "+strconv.Itoa(len(api.TriviaList))+"```")
			return
		}
		api.FetchTrivia(numQuestions, difficulty)
		startTrivia(s, m, numQuestions)
	case config.BotPrefix + "answer":
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "```Usage: !answer [answer]```")
			return
		}

		answerTrivia(s, m, strings.Join(args[1:], " ")) // join the arguments to get the answer
	case config.BotPrefix + "help":
		s.ChannelMessageSend(m.ChannelID, "```Usage: !start [number of questions], !answer [answer], !help```")

	case config.BotPrefix + "ping":
		s.ChannelMessageSend(m.ChannelID, "```No worries, I'm Alive!🚀```")

	case config.BotPrefix + "clear":
		s.ChannelMessageSend(m.ChannelID, "```Clearing the chat```")
		messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, message := range messages {
			s.ChannelMessageDelete(m.ChannelID, message.ID)
		}
		s.ChannelMessageSend(m.ChannelID, "```Chat Cleared```")

	}
}

// startTrivia starts the trivia game.
func startTrivia(s *discordgo.Session, m *discordgo.MessageCreate, numQuestions int) {

	// shuffle the questions
	rand.Shuffle(len(api.TriviaList), func(i, j int) {
		api.TriviaList[i], api.TriviaList[j] = api.TriviaList[j], api.TriviaList[i]
	})

	currentQuestion = 0
	score = 0
	totalQuestion = numQuestions

	// send the first question
	s.ChannelMessageSend(m.ChannelID, "```Q. "+api.TriviaList[currentQuestion].Question+"\nOptions: "+strings.Join(api.TriviaList[currentQuestion].Options, "\n")+"```")
}

// answerTrivia checks if the answer is correct.
func answerTrivia(s *discordgo.Session, m *discordgo.MessageCreate, answer string) {

	// check if the answer is correct
	if strings.ToLower(answer) == strings.ToLower(api.TriviaList[currentQuestion].Answer) {
		score++

		s.ChannelMessageSend(m.ChannelID, "```Correct! Your score is "+strconv.Itoa(score)+"```")
	} else {
		s.ChannelMessageSend(m.ChannelID, "```Incorrect!, Correct answer is "+api.TriviaList[currentQuestion].Answer+"\n Your score is "+strconv.Itoa(score)+"```")
	}

	currentQuestion++

	// check if the trivia is finished
	if currentQuestion >= totalQuestion {
		s.ChannelMessageSend(m.ChannelID, "```Trivia finished! Your score is "+strconv.Itoa(score)+"```")
		return
	}

	s.ChannelMessageSend(m.ChannelID, "```Q. "+api.TriviaList[currentQuestion].Question+"\nOptions: "+strings.Join(api.TriviaList[currentQuestion].Options, "\n")+"```")

}
