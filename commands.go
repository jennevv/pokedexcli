package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
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

const PokemonURL string = "https://pokeapi.co/api/v2/pokemon/"

type PokemonResponse struct {
	BaseExperience int             `json:"base_experience"`
	Height         int             `json:"height"`
	Weight         int             `json:"weight"`
	Stats          []StatsResponse `json:"stats"`
	Types          []TypeResponse  `json:"types"`
}

type StatsResponse struct {
	BaseStat int          `json:"base_stat"`
	Stat     StatResponse `json:"stat"`
}

type StatResponse struct {
	Name string `json:"name"`
}

type TypesResponse struct {
	Type TypeResponse `json:"type"`
}

type TypeResponse struct {
	Name string `json:"name"`
}

func commandCatch(client pokeapi.PokeClient, config *Config, pokemon string) error {
	fmt.Printf("Throwing a ball at %s...\n", pokemon)

	var val []byte
	var err error

	if cachedVal, ok := config.Cache.Get(PokemonURL + pokemon); ok {
		val = cachedVal
	} else {
		val, err = client.Get(PokemonURL + pokemon)
		if err != nil {
			return err
		}

		config.Cache.Add(PokemonURL+pokemon, val)
	}

	var response PokemonResponse
	err = json.Unmarshal(val, &response)
	if err != nil {
		return err
	}

	r := rand.New(rand.NewSource(42))

	if catchRate(response.BaseExperience) > r.Float32() {
		fmt.Printf("%s was caught!", pokemon)
	}

	return err
}

func catchRate(baseExperience int) float32 {
	// Max base experience is Chanseys: 635
	// Smallest possible catch rate is 1%
	return 1.01 - float32(baseExperience)/635
}
