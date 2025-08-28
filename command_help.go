package main

import (
	"fmt"

	"github.com/jennevv/pokedexcli/internal/pokeapi"
)

func commandHelp(client *pokeapi.PokeClient, config *Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}
