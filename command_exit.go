package main

import (
	"fmt"
	"os"

	"github.com/jennevv/pokedexcli/internal/pokeapi"
)

func commandExit(client *pokeapi.PokeClient, config *Config, argument string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
