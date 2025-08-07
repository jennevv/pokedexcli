package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		err := scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
		scanner.Scan()
		input := scanner.Text()
		inputSlice := cleanInput(input)

		firstWord := strings.ToLower(inputSlice[0])

		fmt.Printf("Your command was: %s\n", firstWord)
	}
}

func cleanInput(text string) []string {
	return strings.Fields(text)
}
