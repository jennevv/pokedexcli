package main

import (
	"encoding/json"
	"fmt"

	"github.com/jennevv/pokedexcli/internal/pokeapi"
)

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
