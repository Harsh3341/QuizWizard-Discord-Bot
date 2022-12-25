package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Trivia represents a single trivia question.
type Trivia struct {
	Question string
	Answer   string
}

// TriviaList is a list of trivia questions.
var TriviaList []Trivia

// fetchTrivia fetches trivia questions from the Open Trivia Database API.
func FetchTrivia() {
	url := "https://opentdb.com/api.php?amount=1&encode=base64"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
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

		TriviaList = append(TriviaList, Trivia{string(question), string(answer)})

	}

}
