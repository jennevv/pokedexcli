package main

import (
	"encoding/json"
	"fmt"

	"github.com/jennevv/pokedexcli/internal/pokeapi"
)

type ExploreResponse struct {
	Encounters []Encounter `json:"pokemon_encounters"`
}

type Encounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func commandExplore(client *pokeapi.PokeClient, config *Config, location string) error {
	locationURL := LocationURL + "/" + location

	var val []byte
	var err error

	if cachedVal, ok := config.Cache.Get(locationURL); ok {
		val = cachedVal
	} else {
		val, err = client.Get(locationURL)
		if err != nil {
			return err
		}

		config.Cache.Add(locationURL, val)
	}

	var response ExploreResponse
	err = json.Unmarshal(val, &response)
	if err != nil {
		return err
	}

	printPokemonNames(response)

	return nil
}

func printPokemonNames(response ExploreResponse) {
	for _, encounter := range response.Encounters {
		pokemonName := encounter.Pokemon.Name
		fmt.Println(pokemonName)
	}
}
