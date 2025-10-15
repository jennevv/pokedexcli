package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"

	"github.com/jennevv/pokedexcli/internal/pokeapi"
)

func commandHelp(client *pokeapi.PokeClient, config *Config, argument string) error {
	fmt.Println("Welcome to the Pokedex!\n")
	fmt.Println("Commands:")

	var commandKeys []string

	for key := range commands {
		commandKeys = append(commandKeys, key)
	}

	slices.Sort(commandKeys)

	for _, key := range commandKeys {
		command := commands[key]
		fmt.Printf("%8s  %s\n", command.name, command.description)
	}
	return nil
}

func commandExit(client *pokeapi.PokeClient, config *Config, argument string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

const LocationURL string = "https://pokeapi.co/api/v2/location-area"

type LocationResponse struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Result `json:"results"`
}

type Result struct {
	Name string `json:"name"`
}

func commandMap(client *pokeapi.PokeClient, config *Config, argument string) error {
	if config.Next == "" {
		config.Next = LocationURL
	}

	var val []byte
	var err error

	if cachedVal, ok := config.Cache.Get(config.Next); ok {
		val = cachedVal
	} else {
		val, err = client.Get(config.Next)
		if err != nil {
			return err
		}
	}

	config.Cache.Add(config.Next, val)

	var response LocationResponse
	err = json.Unmarshal(val, &response)
	if err != nil {
		return err
	}

	config.UpdateNavigation(response)

	printMapNames(response)

	return nil
}

func commandMapBack(client *pokeapi.PokeClient, config *Config, argument string) error {
	if config.Previous == "" {
		config.Previous = LocationURL
	}

	var val []byte
	var err error

	if cachedVal, ok := config.Cache.Get(config.Previous); ok {
		val = cachedVal
	} else {
		val, err = client.Get(config.Previous)
		if err != nil {
			return err
		}
	}

	config.Cache.Add(config.Previous, val)

	var response LocationResponse
	err = json.Unmarshal(val, &response)
	if err != nil {
		return err
	}

	config.UpdateNavigation(response)

	printMapNames(response)

	return nil
}

func printMapNames(response LocationResponse) {
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
}

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
	if len(location) == 0 {
		return errors.New("no location provided\n")
	}

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
