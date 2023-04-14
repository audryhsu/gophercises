package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

type problem struct {
	question string
	answer   string
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
	return quiz
}

func (q *Quiz) Score() {
	percentCorrect := float64(q.correct) / float64(q.totalQuestions) * 100.0
	pprint(fmt.Sprintf("You answered %d out of %d questions --> %.2f%%", q.correct, q.totalQuestions, percentCorrect))
}

// Start begins the quiz
func (q *Quiz) Start() {
	file, err := os.Open(q.file)
	if err != nil {
		log.Fatalf("Failed to open csv file: %s", q.file)
	}
	reader := csv.NewReader(file)

	// ReadAll returns a slice of records, where each record is a slice of fields
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Failed to parse csv file")
	}

	// parse lines into a slice of problems
	problems, totalQuestions := q.ParseLines(lines)
	q.totalQuestions = totalQuestions

	// start timer
	timer := time.NewTimer(q.timeLimit)

	// loop over problems
problemLoop:
	for i, p := range problems {
		pprint(fmt.Sprintf("~~* Problem #%d *~~", i+1))
		answerCh := make(chan string)
		// make this a go channel?
		go q.AskQuestion(p, answerCh)

		select {
		case <-timer.C:
			pprint("time's up!")
			break problemLoop
		case answer := <-answerCh:
			if answer == p.answer {
				q.correct++
				pprint("You got it!")
			} else {
				pprint(fmt.Sprintf("Sorry, that's wrong. Correct answer: %v", p.answer))
			}
		}
	}
}

func (q *Quiz) ParseLines(lines [][]string) (problems []problem, totalLines int) {
	for _, line := range lines {
		p := problem{
			question: line[0],
			answer:   line[len(line)-1],
		}

		totalLines++
		problems = append(problems, p)
	}
	return problems, totalLines
}

func (q *Quiz) AskQuestion(p problem, answerChan chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	pprint(p.question + "?")

	for scanner.Scan() {
		input := scanner.Text()
		answerChan <- strings.TrimSpace(input)
		break
	}
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in format question, answer")
	limit := flag.Int("limit", 30, "time limit for the quiz in seconds")
	flag.Parse()

	pprint(fmt.Sprintf("Starting %s quiz with %d time limit...", *csvFilename, *limit))

	config := Config{
		isTimed:       true,
		timeLimit:     time.Second * time.Duration(*limit),
		questionsFile: *csvFilename,
	}
	quiz := New(&config)
	quiz.Start()
	quiz.Score()
}

// pretty print
func pprint(message string) {
	fmt.Printf("\n--- %s\n", message)
}
