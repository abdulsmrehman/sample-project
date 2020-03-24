package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {

	fmt.Println("Program to conduct online quiz")

	//flag is mainly used for passing cmd line arguments while running go executables.
	csvFileName := flag.String("csv", "sampleData.csv", "al file in the folrm of QUIZ!")
	timeLimit := flag.Int("limit", 30, "time limit which tells how long the quiz lasts in SECONDS..")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		exit("Failed to open the File")
	}

	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the CSV file")
	}

	problem := parseLines(lines)
	correct := evaluateQuiz(problem, timeLimit)
	fmt.Printf("\nYOU SCORED %d out of %d.\n", correct, len(problem))
}

func evaluateQuiz(problem []Problem, timeLimit *int) (correct int) {
	correct = 0
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second) // mainly used to add timer to the QUIZ.
	answerCh := make(chan string,2)
	breakLoop:
	for i, p := range problem {

		go func() {
			fmt.Printf("Quiz %d:%s \n", i+1, p.q)
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("Oops Sorry, TimedOUT!!!")
			break breakLoop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	return
}

type Problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []Problem {
	ret := make([]Problem, len(lines))

	for i, line := range lines {
		ret[i] = Problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

func exit(msg string) {
	fmt.Println(msg + "\n")
	os.Exit(1)
}
