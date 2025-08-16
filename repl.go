package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
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
		"map": {
			name:        "map",
			description: "Show the 20 next locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the 20 next locations",
			callback:    commandMapBack,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(text)
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{}

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
			err := command.callback(&config)
			if err != nil {
				fmt.Printf("error calling %s: %v", command.name, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
