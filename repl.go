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
	Pokedex  map[string]Pokemon
}

type Pokemon struct {
	PokedexNo    int
	Name         string
	Species      string
	Types        []string
	Stats        map[string]int
	Height       int
	Weight       int
	NumberCaught int
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
			description: "Displays this help message",
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch the given pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func cleanInput(text string) []string {
	return strings.Fields(text)
}

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{
		Cache:   *pokecache.NewCache(time.Duration(5 * time.Second)),
		Pokedex: make(map[string]Pokemon),
	}
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

		var command string
		if len(inputSlice) == 0 {
			command = "help"
		} else {
			command = strings.ToLower(inputSlice[0])
		}

		var argument string
		if len(inputSlice) > 1 {
			argument = strings.ToLower(inputSlice[1])
		} else {
			argument = ""
		}

		if c, ok := commands[command]; ok {
			err := c.callback(client, &config, argument)
			if err != nil {
				fmt.Printf("error calling %s: %v\n", c.name, err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}
