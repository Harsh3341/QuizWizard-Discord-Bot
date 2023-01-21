package api

import (
	"encoding/base64" // for decoding base64 strings
	"encoding/json"   // for decoding JSON
	"fmt"             // for printing
	"io/ioutil"       // for reading response body
	"math/rand"       // for shuffling options
	"net/http"        // for making HTTP requests
	"reflect"
	"strconv" // for converting int to string

	"github.com/harsh3341/3rd-Semester-Mini-Project/config"
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

type Quiz struct {
	ID                     int
	Question               string
	Description            string
	Answers                []string
	MultipleCorrectAnswers string
	Correct_Answers        []string

	CorrectAnswer string
	Explanation   string
	Tip           string
	Category      string
	Difficulty    string
}

var QuizList []Quiz

// TriviaList is a list of trivia questions.
var TriviaList []Trivia

// fetchTrivia fetches trivia questions from the Open Trivia Database API.
func FetchTrivia(num int, difficulty string) {
	// make HTTP request to API
	url := "https://opentdb.com/api.php?amount=" + strconv.Itoa(num) + "&difficulty=" + difficulty + "&encode=base64"

	// get response body from API
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close() // close response body when function returns

	// read response body
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

	// decode JSON
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

	// decode base64 strings and add to TriviaList
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

		// decode incorrect answers
		var incorrect_answers []string
		for _, incorrect_answer := range result.Incorrect_Answers {
			incorrect_answer, err := base64.StdEncoding.DecodeString(incorrect_answer)
			if err != nil {
				fmt.Println(err.Error())
			}
			incorrect_answers = append(incorrect_answers, string(incorrect_answer))
		}

		// add correct answer to options
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

		// fmt.Println(TriviaList)

	}

}

func FetchQuiz(num int, category string) {

	url := "https://quizapi.io/api/v1/questions?category=" + category + "&limit=" + strconv.Itoa(num) + "&apiKey=" + config.APITOKEN

	// get response body from API
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer resp.Body.Close() // close response body when function returns

	// read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	var data []struct {
		ID          int    `json:"id"`
		Question    string `json:"question"`
		Description string `json:"description"`
		Answers     struct {
			A string `json:"answer_a"`
			B string `json:"answer_b"`
			C string `json:"answer_c"`
			D string `json:"answer_d"`
			E string `json:"answer_e"`
			F string `json:"answer_f"`
		} `json:"answers"`
		MultipleCorrectAnswers string `json:"multiple_correct_answers"`
		CorrectAnswers         struct {
			A string `json:"answer_a_correct"`
			B string `json:"answer_b_correct"`
			C string `json:"answer_c_correct"`
			D string `json:"answer_d_correct"`
			E string `json:"answer_e_correct"`
			F string `json:"answer_f_correct"`
		} `json:"correct_answers"`

		CorrectAnswer string `json:"correct_answer"`
		Explanation   string `json:"explanation"`
		Tip           string `json:"tip"`
		Tags          []struct {
			Name string `json:"name"`
		} `json:"tags"`
		Category   string `json:"category"`
		Difficulty string `json:"difficulty"`
	}

	// decode JSON
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(err.Error())
	}

	//store data values in QuizList
	for _, result := range data {
		answers := []string{result.Answers.A, result.Answers.B, result.Answers.C, result.Answers.D, result.Answers.E, result.Answers.F}

		values := reflect.ValueOf(result.CorrectAnswers)
		correct_answer := []string{}
		for i := 0; i < values.NumField(); i++ {
			valueField := values.Field(i)
			if valueField.String() == "true" {
				correct_answer = append(correct_answer, answers[i])
			}
		}

		QuizList = append(QuizList, Quiz{
			ID:              result.ID,
			Question:        result.Question,
			Description:     result.Description,
			Answers:         answers,
			Correct_Answers: correct_answer,
			Explanation:     result.Explanation,
			Tip:             result.Tip,
			Category:        result.Category,
			Difficulty:      result.Difficulty,
		})
	}
}
