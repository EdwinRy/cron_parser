package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func getCronArg() (string, error) {
	if len(os.Args) != 2 {
		return "", fmt.Errorf("invalid arguments")
	}
	return os.Args[1], nil
}

func printUsage() {
	fmt.Println("Usage: cronParser \"<cron string>\"")
	fmt.Println("Example:")
	fmt.Printf("\tcronParser \"*/15 0 1,15 * 1-5 /usr/bin/find\"\n\n")
	fmt.Println("\tOutput: ")
	fmt.Println("\tminute         0 15 30 45")
	fmt.Println("\thour           0")
	fmt.Println("\tday of month   1 15")
	fmt.Println("\tmonth          1 2 3 4 5 6 7 8 9 10 11 12")
	fmt.Println("\tday of week    1 2 3 4 5")
	fmt.Printf("\tcommand        /usr/bin/find\n\n")
}

func CronTaskCompile(cronStr string) (*CronTask, error) {

	_, debug := os.LookupEnv("DEBUG")

	// Convert raw string into a list of tokens
	tokens, err := Tokenize(cronStr)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize your cron string: %v", err)
	}
	if debug {
		fmt.Println(tokens)
	}

	// Convert tokens into an abstract syntax tree
	ast, err := Parse(tokens)
	if err != nil {
		return nil, fmt.Errorf("could not parse cron task: %v", err)
	}
	if debug {
		b, err := json.MarshalIndent(ast, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("could not debug the AST: %v", err)
		}
		fmt.Println(string(b))
	}

	// Convert the abstract syntax tree into a semantic cron task object
	task, err := GetCronTask(ast)
	if err != nil {
		return nil, fmt.Errorf("failed to extract valid cron task from syntax: %v", err)
	}
	return task, nil
}

func main() {
	cronStr, err := getCronArg()
	if err != nil {
		printUsage()
		return
	}

	cronTask, err := CronTaskCompile(cronStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(cronTask)
}
