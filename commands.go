package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"time"

	"github.com/jennevv/pokedexcli/internal/pokeapi"
)

func commandHelp(client *pokeapi.PokeClient, config *Config, argument string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Commands:")

	var commandKeys []string

	for key := range commands {
		if key == "help" {
			continue
		}
		commandKeys = append(commandKeys, key)
	}

	// Print help command first
	fmt.Printf("%8s  %s\n", commands["help"].name, commands["help"].description)

	// Print commands alphabetically
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
	Pokemon EncounterPokemon `json:"pokemon"`
}

type EncounterPokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func commandExplore(client *pokeapi.PokeClient, config *Config, location string) error {
	if len(location) == 0 {
		return errors.New("no location provided")
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
	ID             int            `json:"id"`
	Name           string         `json:"name"`
	BaseExperience int            `json:"base_experience"`
	Height         int            `json:"height"`
	Weight         int            `json:"weight"`
	Species        PokemonSpecies `json:"species"`
	Stats          []PokemonStats `json:"stats"`
	Types          []PokemonTypes `json:"types"`
}

type PokemonSpecies struct {
	Name string `json:"name"`
}

type PokemonStats struct {
	BaseStat int         `json:"base_stat"`
	Stat     PokemonStat `json:"stat"`
}

type PokemonStat struct {
	Name string `json:"name"`
}

type PokemonTypes struct {
	Type PokemonType `json:"type"`
}

type PokemonType struct {
	Name string `json:"name"`
}

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func commandCatch(client *pokeapi.PokeClient, config *Config, pokemon string) error {
	if len(pokemon) == 0 {
		return errors.New("no pokemon specified")
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon)

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

	if catchRate(response.BaseExperience) > rng.Float32() {
		fmt.Printf("%s was caught!\n", pokemon)
		addToPokedex(config.Pokedex, pokemon, response)
	} else {
		fmt.Printf("%s escaped!\n", pokemon)
	}

	return err
}

func catchRate(baseExperience int) float32 {
	// Max base experience is Chanseys: 635
	// Smallest possible catch rate is 0.1%
	return 1.0 - float32(baseExperience)/395 + 0.001
}

func addToPokedex(pokedex map[string]Pokemon, pokemonKey string, pokemonResponse PokemonResponse) {
	// pokemon already caught and in pokedex
	if _, ok := pokedex[pokemonKey]; ok {
		pokemon := pokedex[pokemonKey]
		pokemon.NumberCaught++
		pokedex[pokemonKey] = pokemon
		return
	}

	var types []string
	stats := make(map[string]int)

	for _, typeInfo := range pokemonResponse.Types {
		types = append(types, typeInfo.Type.Name)
	}

	for _, statInfo := range pokemonResponse.Stats {
		stats[statInfo.Stat.Name] = statInfo.BaseStat
	}

	pokedex[pokemonKey] = Pokemon{
		PokedexNo:    pokemonResponse.ID,
		Name:         pokemonResponse.Name,
		Species:      pokemonResponse.Species.Name,
		Types:        types,
		Stats:        stats,
		Height:       pokemonResponse.Height,
		Weight:       pokemonResponse.Weight,
		NumberCaught: 1,
	}
}

func commandInspect(client *pokeapi.PokeClient, config *Config, pokemon string) error {
	if len(pokemon) == 0 {
		return errors.New("no pokemon specified")
	}

	if pokemonInfo, ok := config.Pokedex[pokemon]; !ok {
		fmt.Println("you have not caught that pokemon")
	} else {
		printPokemonInfo(pokemonInfo)
	}
	return nil
}

func printPokemonInfo(pokemonInfo Pokemon) {
	fmt.Printf("Name: %s\n", pokemonInfo.Name)
	fmt.Printf("Height: %d cm\n", pokemonInfo.Height*10)
	fmt.Printf("Weight: %.1f kg\n", float32(pokemonInfo.Weight)/10.0)
	fmt.Println("Stats:")
	for s, statVal := range pokemonInfo.Stats {
		fmt.Printf("%4s- %s: %d\n", "", s, statVal)
	}
	fmt.Println("Types:")
	for _, t := range pokemonInfo.Types {
		fmt.Printf("%4s- %s\n", "", t)
	}
}

func commandPokedex(client *pokeapi.PokeClient, config *Config, argument string) error {
	if len(config.Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty.")
	} else {
		for _, pokemon := range config.Pokedex {
			fmt.Printf("%4s- %s\n", "", pokemon.Name)
		}
	}
	return nil
}
