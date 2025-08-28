package main

import (
	"fmt"

	"github.com/jennevv/pokedexcli/internal/pokeapi"
)

const LOCATION_URL string = "https://pokeapi.co/api/v2/location-area"

func commandMap(client *pokeapi.PokeClient, config *Config) error {
	if config.Next == "" {
		config.Next = LOCATION_URL
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

	response, err := client.Unmarshal(val)
	if err != nil {
		return err
	}

	config.UpdateNavigation(response)

	printMapNames(response)

	return nil
}

func commandMapBack(client *pokeapi.PokeClient, config *Config) error {
	if config.Previous == "" {
		config.Previous = LOCATION_URL
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

	response, err := client.Unmarshal(val)
	if err != nil {
		return err
	}

	config.UpdateNavigation(response)

	printMapNames(response)

	return nil
}

func printMapNames(response pokeapi.Response) {
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
}
