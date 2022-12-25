package bot

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/harsh3341/3rd-Semester-Mini-Project/config"
)

var (
	BotID           string
	goBot           *discordgo.Session
	TriviaList      []Trivia // list of trivia questions
	totalQuestion   int      // number of questions to ask
	currentQuestion int
	score           int
)

func init() {
	// just some random questions(no offense)
	TriviaList = []Trivia{
		{
			Question: "Who is the Siyar of B2",
			Answer:   "Kushal",
		},
		{
			Question: "Who is the queen of B2",
			Answer:   "Harsh Sahu",
		},
		{
			Question: "Who is the Captain of B2",
			Answer:   "Divyansh",
		},
		{
			Question: "Who is the Siyar 2.O of B2",
			Answer:   "Hemendra",
		},
		{
			Question: "Where do we study",
			Answer:   "SSTC",
		},
	}
}

// Trivia represents a single trivia question.
type Trivia struct {
	Question string
	Answer   string
}

func Start() {
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
	}

	BotID = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Bot is running...")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotID {
		return
	}

	// check if the message is a command
	if !strings.HasPrefix(m.Content, config.BotPrefix) {
		return
	}

	args := strings.Split(m.Content, " ")

	// if m.Content == config.BotPrefix+"ping" {
	// 	_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	// }

	switch args[0] {
	case config.BotPrefix + "start":
		if len(args) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: !start [number of questions]")
			return
		}

		numQuestions, err := strconv.Atoi(args[1]) // convert string to int
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error: Invalid number of questions")
			return
		}

		if numQuestions < 1 || numQuestions > len(TriviaList) {
			s.ChannelMessageSend(m.ChannelID, "Error: Number of questions must be between 1 and "+strconv.Itoa(len(TriviaList)))
			return
		}

		startTrivia(s, m, numQuestions)
	case config.BotPrefix + "answer":
		if len(args) != 2 {
			s.ChannelMessageSend(m.ChannelID, "Usage: !answer [answer]")
			return
		}
		answerTrivia(s, m, args[1])
	case config.BotPrefix + "help":
		s.ChannelMessageSend(m.ChannelID, "Usage: !start [number of questions], !answer [answer], !help")

	case config.BotPrefix + "ping":
		s.ChannelMessageSend(m.ChannelID, "No worries, I'm Alive!")

	case config.BotPrefix + "clear":
		s.ChannelMessageSend(m.ChannelID, "Clearing the chat")
		messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, message := range messages {
			s.ChannelMessageDelete(m.ChannelID, message.ID)
		}

	}
}

func startTrivia(s *discordgo.Session, m *discordgo.MessageCreate, numQuestions int) {

	rand.Shuffle(len(TriviaList), func(i, j int) {
		TriviaList[i], TriviaList[j] = TriviaList[j], TriviaList[i]
	})

	currentQuestion = 0
	score = 0
	totalQuestion = numQuestions

	s.ChannelMessageSend(m.ChannelID, TriviaList[currentQuestion].Question)
}

func answerTrivia(s *discordgo.Session, m *discordgo.MessageCreate, answer string) {
	if answer == strings.ToLower(TriviaList[currentQuestion].Answer) {
		score++

		s.ChannelMessageSend(m.ChannelID, "Correct! Your score is "+strconv.Itoa(score))
	} else {
		s.ChannelMessageSend(m.ChannelID, "Incorrect! Your score is "+strconv.Itoa(score))
	}

	currentQuestion++

	if currentQuestion >= totalQuestion {
		s.ChannelMessageSend(m.ChannelID, "Trivia finished! Your score is "+strconv.Itoa(score))
		return
	}

	s.ChannelMessageSend(m.ChannelID, TriviaList[currentQuestion].Question)

}
