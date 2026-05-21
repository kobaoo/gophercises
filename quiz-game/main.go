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
	csvFileName := flag.String("cs", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	
	file, err := os.Open(*csvFileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file %s\n", *csvFileName))
	}
	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	problems := parseLines(lines)
	
	var rightAnswers, correct int 

  problemLoop:
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s -> ", i+1, p.q )
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C:
			break problemLoop
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	finish(fmt.Sprintf("Quiz completed!\nYour score is %d out of %d.", rightAnswers, len(lines)))
}

func parseLines(lines [][]string) []problems {
	ret := make([]problems, len(lines))
	for i, line := range lines {
		ret[i] = problems{
			a : line[1],
			q : strings.TrimSpace(line[0]),
		}
	}
	return ret
}

type problems struct {
	q string
	a string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func finish(msg string) {
	fmt.Println(msg)
	os.Exit(0)
}