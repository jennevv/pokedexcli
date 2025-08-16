package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Response struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous string   `json:"previous"`
	Results  []Result `json:"results"`
}

type Result struct {
	Name string `json:"name"`
}

var LOCATION_URL string = "https://pokeapi.co/api/v2/location-area"

func commandMap(config *Config) error {
	if config.Next == "" {
		config.Next = LOCATION_URL
	}

	res, err := http.Get(config.Next)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	config.Previous = response.Previous
	config.Next = response.Next

	printMapNames(response)

	return nil
}

func commandMapBack(config *Config) error {
	if config.Previous == "" {
		config.Previous = LOCATION_URL
	}

	res, err := http.Get(config.Previous)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	config.Previous = response.Previous
	config.Next = response.Next

	printMapNames(response)

	return nil
}

func printMapNames(response Response) {
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
}
