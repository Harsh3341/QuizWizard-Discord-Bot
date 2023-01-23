package bot

import (
	"fmt"       // for printing
	"math/rand" // for shuffling options
	"strconv"   // for converting int to string
	"strings"   // for splitting strings

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"                         // for Discord bot
	"github.com/harsh3341/3rd-Semester-Mini-Project/api"    // for trivia questions
	"github.com/harsh3341/3rd-Semester-Mini-Project/config" // for bot token and prefix
)

// BotID is the ID of the bot.
var (
	BotID            string
	goBot            *discordgo.Session
	totalQuestionT   int // number of questions to ask
	currentQuestionT int
	scoreT           int
	totalQuestionQ   int // number of questions to ask
	currentQuestionQ int
	scoreQ           int
)

var buttons = []string{":one:", ":two:", ":three:", ":four:", ":five:", ":six:", ":seven:", ":eight:", ":nine:", ":keycap_ten:"}

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
	case config.BotPrefix + "startt":
		if len(args) != 3 {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Error", "```Usage: !startt [number of questions] [difficulty]```", 0xb40000))
			return
		}

		numQuestions, err := strconv.Atoi(args[1]) // convert string to int
		if err != nil {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Error", "```Number of questions must be an integer```", 0xb40000))
			return
		}

		difficulty := args[2]

		if numQuestions < 1 || numQuestions > 50 {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Error", "```Number of questions must be between 1 and 50```", 0xb40000))
			return
		}
		api.FetchTrivia(numQuestions, difficulty)
		startTrivia(s, m, numQuestions)
	case config.BotPrefix + "startq":
		if len(args) != 3 {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Error", "```Usage: !startq [number of questions] [category]```", 0xb40000))
			return
		}

		numQuestions, err := strconv.Atoi(args[1]) // convert string to int
		if err != nil {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Error", "```Number of questions must be an integer```", 0xb40000))
			return
		}

		category := args[2]

		if numQuestions < 1 || numQuestions > 20 {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Error", "```Number of questions must be between 1 and 20```", 0xb40000))
			return
		}
		api.FetchQuiz(numQuestions, category)
		startQuiz(s, m, numQuestions)

	case config.BotPrefix + "answert":
		if len(args) < 2 {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Error", "```Usage: !answert [answer]```", 0xb40000))
			return
		}

		answerTrivia(s, m, strings.Join(args[1:], " ")) // join the arguments to get the answer

	case config.BotPrefix + "answerq":

		if len(args) < 2 {
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Error", "```Usage: !answerq [answer]```", 0xb40000))
			return
		}

		answerQuiz(s, m, strings.Join(args[1:], " ")) // join the arguments to get the answer
	case config.BotPrefix + "help":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Help", "```Commands: \n!startt [number of questions] [difficulty] \n!startq [number of questions] [category] \n!answert [answer] \n!answerq [answer] \n!ping \n!clear```", 0x00b4b4))

	case config.BotPrefix + "ping":
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewEmbed().SetTitle("No worries, I'm Alive!🚀").SetColor(0x00ff00).MessageEmbed)

	case config.BotPrefix + "clear":

		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewEmbed().SetTitle("Clearing Chat...").SetColor(0x00ff00).MessageEmbed)
		messages, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, message := range messages {
			s.ChannelMessageDelete(m.ChannelID, message.ID)
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewEmbed().SetTitle("Chat Cleared!").SetColor(0x00ff00).MessageEmbed)

	}
}

// startTrivia starts the trivia game.
func startTrivia(s *discordgo.Session, m *discordgo.MessageCreate, numQuestions int) {

	// shuffle the questions
	rand.Shuffle(len(api.TriviaList), func(i, j int) {
		api.TriviaList[i], api.TriviaList[j] = api.TriviaList[j], api.TriviaList[i]
	})

	currentQuestionT = 0
	scoreT = 0
	totalQuestionT = numQuestions

	// send the first question
	var Options []string

	for i := 0; i < len(api.TriviaList[currentQuestionT].Options); i++ {
		Options = append(Options, buttons[i]+" "+api.TriviaList[currentQuestionT].Options[i])
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Question "+strconv.Itoa(currentQuestionT+1)+"\n\n"+api.TriviaList[currentQuestionT].Question+"\n\n", strings.Join(Options, "\n"), 0x00ff00))
}

func startQuiz(s *discordgo.Session, m *discordgo.MessageCreate, numQuestions int) {

	// shuffle the questions
	rand.Shuffle(len(api.QuizList), func(i, j int) {
		api.QuizList[i], api.QuizList[j] = api.QuizList[j], api.QuizList[i]
	})

	currentQuestionQ = 0
	scoreQ = 0
	totalQuestionQ = numQuestions

	// send the first question
	var Options []string

	for i := 0; i < len(api.QuizList[currentQuestionQ].Answers); i++ {
		if api.QuizList[currentQuestionQ].Answers[i] != "" {
			Options = append(Options, buttons[i]+" "+api.QuizList[currentQuestionQ].Answers[i])
		}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Question "+strconv.Itoa(currentQuestionQ+1)+"\n\n"+api.QuizList[currentQuestionQ].Question+"\n\n", strings.Join(Options, "\n"), 0x00ff00))
}

// answerTrivia checks if the answer is correct.
func answerTrivia(s *discordgo.Session, m *discordgo.MessageCreate, answer string) {

	// check if the answer is correct
	if strings.ToLower(answer) == strings.ToLower(api.TriviaList[currentQuestionT].Answer) {
		scoreT++

		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Correct!", "```Your score is "+strconv.Itoa(scoreT)+"```", 0x00ff00))
	} else {

		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Incorrect!", "```Correct answer is "+api.TriviaList[currentQuestionT].Answer+"\nYour score is "+strconv.Itoa(scoreT)+"```", 0xb40000))
	}

	currentQuestionT++

	// check if the trivia is finished
	if currentQuestionT >= totalQuestionT {
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Trivia Finished!", "```Your score is "+strconv.Itoa(scoreT)+"```", 0x00ff00))
		return
	}

	var Options []string

	for i := 0; i < len(api.TriviaList[currentQuestionT].Options); i++ {
		Options = append(Options, buttons[i]+" "+api.TriviaList[currentQuestionT].Options[i])
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Question "+strconv.Itoa(currentQuestionT+1)+"\n\n"+api.TriviaList[currentQuestionT].Question+"\n\n", strings.Join(Options, "\n"), 0x00ff00))
}

func answerQuiz(s *discordgo.Session, m *discordgo.MessageCreate, answer string) {

	// check if the answer is correct

	for i := 0; i < len(api.QuizList[currentQuestionQ].Correct_Answers); i++ {

		if strings.ToLower(answer) == strings.ToLower(api.QuizList[currentQuestionQ].Correct_Answers[i]) {

			scoreQ++
			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Correct!", "```Your score is "+strconv.Itoa(scoreQ)+"```", 0x00ff00))

		} else {

			s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Incorrect!", "```Correct answer is "+api.QuizList[currentQuestionQ].Correct_Answers[i]+"\nYour score is "+strconv.Itoa(scoreQ)+"```", 0xb40000))
		}
	}

	currentQuestionQ++

	if currentQuestionQ >= totalQuestionQ {
		s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Quiz Finished!", "```Your score is "+strconv.Itoa(scoreQ)+"```", 0x00ff00))
		return
	}

	var Options []string

	for i := 0; i < len(api.QuizList[currentQuestionQ].Answers); i++ {
		if api.QuizList[currentQuestionQ].Answers[i] != "" {
			Options = append(Options, buttons[i]+" "+api.QuizList[currentQuestionQ].Answers[i])
		}
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed.NewGenericEmbedAdvanced("Question "+strconv.Itoa(currentQuestionQ+1)+"\n\n"+api.QuizList[currentQuestionQ].Question+"\n\n", strings.Join(Options, "\n"), 0x00ff00))

}
