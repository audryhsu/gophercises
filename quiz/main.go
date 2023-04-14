package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Config struct {
	isTimed       bool
	timeLimit     time.Duration
	questionsFile string
}

type Quiz struct {
	correct        int
	totalQuestions int
	file           string // path to file?
	timeLimit      time.Duration
	isTimed        bool
	timer          *time.Timer
}

func New(config *Config) *Quiz {
	quiz := &Quiz{
		file:           config.questionsFile,
		totalQuestions: 0,
		correct:        0,
	}

	if config.isTimed {
		quiz.isTimed = true
		quiz.timeLimit = config.timeLimit
	}

	var processAndCountQuestions = csvProcessor(quiz, quiz.CountQuestion())
	processAndCountQuestions()

	fmt.Printf("number of questions: %d\n", quiz.totalQuestions)
	return quiz
}

func (q *Quiz) Score() {
	percentCorrect := float64(q.correct/q.totalQuestions) * 100
	fmt.Printf("You answered %d out of %d questions --> %f score\n", q.correct, q.totalQuestions, percentCorrect)
}

// Start begins the quiz
func (q *Quiz) Start() {
	// start the timer
	var timerChan chan bool
	if q.isTimed {
		q.timer = time.NewTimer(q.timeLimit)

		// create a channel to receive timer events
		timerChan = make(chan bool)
		// go routine monitors the timer by continuously waiting for timer event
		go func() {
			for {
				select {
				case <-q.timer.C:
					fmt.Println("\nTime's up!")
					// send a value on channel and return
					timerChan <- true
					return
				default:
				}
			}
		}()
	}

	processAndAskQuestion := csvProcessor(q, q.AskQuestion())
	processAndAskQuestion()

	// stop timer if it's still running and quiz is complete
	if q.isTimed && !q.timer.Stop() {
		<-timerChan
	}
}

// CountQuestion returns a function with a closure that includes the quiz it is processing.
func (q *Quiz) CountQuestion() func(row []string) {
	return func(row []string) {
		q.totalQuestions++
	}
}

func (q *Quiz) AskQuestion() func(row []string) {
	return func(row []string) {
		question := row[0]
		answer := row[len(row)-1]

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println(question)
		for scanner.Scan() {
			input := scanner.Text()

			if input == answer {
				q.correct++
				fmt.Println("You got it!")
			} else {
				fmt.Printf("Sorry, that's wrong. Correct answer: %v\n", answer)
			}
			break
		}
	}
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in format question, answer")
	limit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	config := Config{
		isTimed:       true,
		timeLimit:     time.Second * time.Duration(*limit),
		questionsFile: *csvFilename,
	}
	quiz := New(&config)
	quiz.Start()
	quiz.Score()
}

// // // //  HELPERS

// csvProcessor returns a function processes each row while reading a file from open to close
func csvProcessor(q *Quiz, processor func(row []string)) func() {
	return func() {
		file, err := os.Open(q.file)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()

		// NewReader takes an io.Reader, and file implements the io.Reader interface
		csvReader := csv.NewReader(file)
		csvReader.Comma = ','

		for {
			// loop over each record until EOF; a record is a slice of fields
			row, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			if q.timer != nil {

				if q.timer.Stop() {
					fmt.Println("timer hasn't stopped yet")
				}
				//// check if timer fired yet
				//select {
				//case <-q.timer.C:
				//	fmt.Println("Time's up!")
				//	return
				//default:
				//	// continue
				//}
			}
			if processor != nil {
				processor(row)
			}
		}
	}
}
