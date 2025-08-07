package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(text)
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		err := scanner.Err()
		if err != nil {
			fmt.Println(err)
		}
		scanner.Scan()
		input := scanner.Text()
		inputSlice := cleanInput(input)

		firstWord := strings.ToLower(inputSlice[0])

		if command, ok := commands[firstWord]; ok {
			err := command.callback()
			if err != nil {
				fmt.Printf("error calling %s: %v", command.name, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
