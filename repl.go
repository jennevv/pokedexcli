package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jennevv/pokedexcli/internal/pokeapi"
	"github.com/jennevv/pokedexcli/internal/pokecache"
)

type Config struct {
	Next     string
	Previous string
	Cache    pokecache.Cache
}

func (c *Config) UpdateNavigation(response LocationResponse) {
	c.Next = response.Next
	c.Previous = response.Previous
}

type cliCommand struct {
	name        string
	description string
	callback    func(*pokeapi.PokeClient, *Config, string) error
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
			description: "Show the 20 previous locations",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore the given location",
			callback:    commandExplore,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(text)
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{Cache: *pokecache.NewCache(time.Duration(5 * time.Second))}
	client := pokeapi.NewClient()

	for {
		fmt.Print("Pokedex > ")
		err := scanner.Err()
		if err != nil {
			fmt.Println(err)
		}
		scanner.Scan()
		input := scanner.Text()
		inputSlice := cleanInput(input)

		command := strings.ToLower(inputSlice[0])

		var argument string
		if len(inputSlice) > 1 {
			argument = strings.ToLower(inputSlice[1])
		} else {
			argument = ""
		}

		if c, ok := commands[command]; ok {
			err := c.callback(client, &config, argument)
			if err != nil {
				fmt.Printf("error calling %s: %v", c.name, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
