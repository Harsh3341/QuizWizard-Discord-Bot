package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
)

// Trivia represents a single trivia question.
type Trivia struct {
	Category          string
	Type              string
	Difficulty        string
	Question          string
	Answer            string
	Incorrect_Answers []string
	Options           []string
}

// TriviaList is a list of trivia questions.
var TriviaList []Trivia

// fetchTrivia fetches trivia questions from the Open Trivia Database API.
func FetchTrivia(num int) {
	url := "https://opentdb.com/api.php?amount=" + strconv.Itoa(num) + "&encode=base64"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	// get response code from API
	var response struct {
		Response_Code int
	}
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err.Error())
	}

	// if response code is not 0, then there was an error
	if response.Response_Code != 0 {
		fmt.Println("Error: Response code is not 0")
	}

	var data struct {
		Results []struct {
			Category          string
			Type              string
			Difficulty        string
			Question          string
			Correct_Answer    string
			Incorrect_Answers []string
		}
	}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err.Error())
	}

	for _, result := range data.Results {
		question, err := base64.StdEncoding.DecodeString(result.Question)
		if err != nil {
			fmt.Println(err.Error())
		}
		answer, err := base64.StdEncoding.DecodeString(result.Correct_Answer)
		if err != nil {
			fmt.Println(err.Error())
		}

		xyz, err := base64.StdEncoding.DecodeString(result.Type)
		if err != nil {
			fmt.Println(err.Error())
		}

		difficulty, err := base64.StdEncoding.DecodeString(result.Difficulty)
		if err != nil {
			fmt.Println(err.Error())
		}

		category, err := base64.StdEncoding.DecodeString(result.Category)
		if err != nil {
			fmt.Println(err.Error())
		}

		var incorrect_answers []string
		for _, incorrect_answer := range result.Incorrect_Answers {
			incorrect_answer, err := base64.StdEncoding.DecodeString(incorrect_answer)
			if err != nil {
				fmt.Println(err.Error())
			}
			incorrect_answers = append(incorrect_answers, string(incorrect_answer))
		}

		options := append(incorrect_answers, string(answer))
		// shuffle options
		for i := range options {
			j := rand.Intn(i + 1)
			options[i], options[j] = options[j], options[i]
		}

		// add question to TriviaList
		TriviaList = append(TriviaList, Trivia{
			Category:          string(category),
			Type:              string(xyz),
			Difficulty:        string(difficulty),
			Question:          string(question),
			Answer:            string(answer),
			Incorrect_Answers: incorrect_answers,
			Options:           options,
		})

	}

	fmt.Println("TriviaList: ", TriviaList)

}
