package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
)

var debugMode bool

type problem struct {
	question string
	answer   string
}

// MAIN
func main() {
	inputFile := flagHandler()

	debugPrint(fmt.Sprintf("Reading quiz from %v\n", inputFile))

	file, err := os.Open(inputFile)
	if err != nil {
		exit(fmt.Sprintf("There was an error when attemptring to read %v:\n%v", inputFile, err), 1)
	}

	reader := csv.NewReader(file)
	fileLines, err := reader.ReadAll()
	if err != nil {
		exit(fmt.Sprintf("Failed to parse the csv file %v", inputFile), 2)
	}

	problems := convertLinesToProblems(fileLines)

	debugQuiz(problems)

	var score int = 0
	for i, nextProb := range problems {
		fmt.Printf("Question %d:\n\t%v = ", i+1, nextProb.question)
		var answer string
		answerCount, err := fmt.Scanf("%s\n", &answer)

		debugPrint(fmt.Sprintf("Answers Read: %d, Error: %v\n", answerCount, err))

		if answer == nextProb.answer {
			score++
			debugPrint("Correct!\n")
		} else {
			debugPrint("Wrong!\n")
		}
	}

	fmt.Printf("You scored %d/%d", score, len(problems))

}

func convertLinesToProblems(lines [][]string) []problem {
	result := make([]problem, len(lines))

	for i, line := range lines {
		// Trimming spaces from answer of each line in file  as the same is done
		// by Scanf when reading user provided Having answer from standard input.
		result[i] = problem{line[0], strings.TrimSpace(line[1])}
	}

	return result
}

func exit(message string, exitCode int) {
	fmt.Println(message)
	os.Exit(exitCode)
}

// Deals with options we want to provide the script from the command line.
func flagHandler() string {
	csvFilePath := flag.String("csv", "problem.csv", "a csv file in the format of 'question,answer'")

	debugModePointer := flag.Bool("debug", false, "Turns on debug mode. Doesn't run test but instead prints full quiz contents.")

	flag.Parse() // If this isn't called after all flags are defined they will contain default values when accessed.
	debugMode = *debugModePointer
	return *csvFilePath
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
