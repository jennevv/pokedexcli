package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Next     string
	Previous string
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

type LocationInfo struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
}

func commandMap(config *Config) error {
	locations := make(map[int]string, 20)
	for i := 0; i < 20; i++ {
		res, err := http.Get(config.Next)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return err
		}

		var locationInfo LocationInfo
		err = json.Unmarshal(body, &locationInfo)
		if err != nil {
			return err
		}

		locations[i] = locationInfo.Location.Name

		config.Previous = config.Next
		config.Next, err = incrementURL(config.Next)
		if err != nil {
			return err
		}
	}

	printMap(locations)

	return nil
}

func commandMapBack(config *Config) error {
	locations := make(map[int]string, 20)
	for i := 0; i < 20; i++ {
		res, err := http.Get(config.Previous)
		if err != nil {
			return err
		}
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if err != nil {
			return err
		}

		var locationInfo LocationInfo
		err = json.Unmarshal(body, &locationInfo)
		if err != nil {
			return err
		}

		locations[len(locations)-i-1] = locationInfo.Location.Name

		config.Next = config.Previous
		config.Previous, err = decrementURL(config.Previous)
		if err != nil {
			return err
		}
	}

	printMap(locations)

	return nil
}

func printMap(locations map[int]string) {
	for _, location := range locations {
		fmt.Println(location)
	}
}

func incrementURL(url string) (string, error) {
	splitURL := strings.Split(url, "/")

	currentNum, err := strconv.Atoi(splitURL[len(splitURL)-1])
	if err != nil {
		return "", err
	}

	splitURL[len(splitURL)-1] = string(currentNum + 1)

	newURL := strings.Join(splitURL, "/")
	return newURL, nil
}

func decrementURL(url string) (string, error) {
	splitURL := strings.Split(url, "/")

	currentNum, err := strconv.Atoi(splitURL[len(splitURL)-1])
	if err != nil {
		return "", err
	}

	splitURL[len(splitURL)-1] = string(currentNum - 1)

	newURL := strings.Join(splitURL, "/")
	return newURL, err
}
