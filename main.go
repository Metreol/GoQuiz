package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var debugMode bool

type problem struct {
	question string
	answer   string
}

func main() {
	// HANDLE USER INPU:
	inputFile, limit := flagHandler()
	debugPrint(fmt.Sprintf("Reading quiz from %v, time limit:%d\n", inputFile, limit))

	// OBTAIN QUIZ FROM FILE
	fileLines := readLinesFromFile(inputFile)
	problems := convertLinesToProblems(fileLines)
	debugQuiz(problems)

	// RUN QUIZ
	// Set up the timer which...
	timer := time.NewTimer(time.Duration(limit) * time.Second)
	//.. runs on a channel accessed as follows
	timerCh := timer.C
	// Set up a channel to listen for answers
	answerCh := make(chan string)

	var score int = 0
	for i, nextProb := range problems {
		fmt.Printf("Question %d:\n\t%v = ", i+1, nextProb.question)
		// go keyword puts this anonymous func in a goroutine, so the program continues while this goroutine waits for the users answer
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer // sends user answer to the answerCh channel
		}() // The () here means the annonymous function is run here.

		// select lets us listen to multiple channels so we can act on whichever case is met first (sort of like a while loop?)
		select {
		case <-timerCh: // Is true if closed (timer has run out).
			fmt.Printf("\nTime ran out! You scored %d/%d", score, len(problems))
			return
		case answer := <-answerCh: // Is true when the answerCh is sent a string value and we assign that value to answer.
			if answer == nextProb.answer {
				score++
				debugPrint("Correct!\n")
			} else {
				debugPrint("Wrong!\n")
			}
		}
	}
	// Print score
	fmt.Printf("You scored %d/%d\n", score, len(problems))
}

// Deals with options we want to provide the script from the command line.
func flagHandler() (string, int) {
	csvFilePath := flag.String("csv", "problem.csv", "The path to a csv file. Must have the format of 'question,answer'")

	debugModePointer := flag.Bool("debug", false, "Turns on debug mode. Doesn't run test but instead prints full quiz contents.")

	timeLimit := flag.Int("time-limit", 30, "The time limit for the quiz in seconds.")

	flag.Parse() // If this isn't called after all flags are defined they will contain default values when accessed.
	debugMode = *debugModePointer
	return *csvFilePath, *timeLimit
}

// Open file at path inputFile for reading (exits if any errors when attempting this) and reads all lines from the opened file
func readLinesFromFile(inputFile string) [][]string {
	file, err := os.Open(inputFile)
	if err != nil {
		exit(fmt.Sprintf("There was an error when attemptring to read %v:\n%v", inputFile, err), 1)
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the csv file %v", inputFile), 2)
	}

	return lines
}

// Converts the lines from the given file ([][]string) to an problem slice ([]problem)
func convertLinesToProblems(lines [][]string) []problem {
	result := make([]problem, len(lines))

	for i, line := range lines {
		// Trimming spaces from answer of each line in file  as the same is done
		// by Scanf when reading user provided Having answer from standard input.
		result[i] = problem{line[0], strings.TrimSpace(line[1])}
	}

	return result
}

// Exits the program with given exitCode after printing message
func exit(message string, exitCode int) {
	fmt.Println(message)
	os.Exit(exitCode)
}

// DEBUG FUNCS
func debugPrint(message string) {
	if debugMode {
		fmt.Print(message)
	}
}

func debugQuiz(problems []problem) {
	if debugMode {
		for i, nextProb := range problems {
			fmt.Printf("Question %d: %v = %v\n", i+1, nextProb.question, nextProb.answer)
		}
	}
}
