package main

import (
	"cyoa/internal/models"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	filename := flag.String("file", "gopher.json", "path for json file")
	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	var stories models.Story
	if err := json.NewDecoder(file).Decode(&stories); err != nil {
		log.Fatalf("Failed to decode json: %s", err)
	}
	var input string
	printChapter(stories["intro"])
	for {
		fmt.Scanf("%s\n", &input)
		if input == "restart" {
			input = "intro"
		} else if input == "exit" {
			break
		}
		printChapter(stories[input])
	}
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

func printChapter(chap models.Chapter) {
	fmt.Printf("Chapter: %s\n", chap.Title)
	for _, p := range chap.Paragraphs {
		fmt.Printf("%s\n", p)
	}
	if len(chap.Options) == 0 {
		fmt.Println("Type " + Yellow + "restart" + Reset + " or " + Yellow + "exit" + Reset)
		return
	}
	fmt.Println("Choose Your Own Adventure !!!")
	for _, o := range chap.Options {
		fmt.Println("\n", o.Text)
		fmt.Println("To choose print: " + Green + o.Chapter + Reset)
		fmt.Println()
	}
}
