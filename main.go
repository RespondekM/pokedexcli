package main

import (
	"bufio"
	"fmt"
	"github.com/RespondekM/pokedexcli/internal"
	"io"
	"net/http"
	"encoding/json"
	"os"
	"strings")


type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	Next		string `json:"next"`
	Previous	string `json:"previous"`
}

type parameters struct {
	Parameter	string
}

type locationArea struct {
	Name		string `json:"name"`
	URL 		string `json:"url"`
}

type response struct {
	Count		int    `json:"count"`
	Next		string `json:"next"`
	Previous	string `json:"previous"`
	Results 	[]locationArea `json:"results"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	opts := new(config)
	opts.Next = "https://pokeapi.co/api/v2/location-area/"
	opts.Previous = ""
	for true {
		fmt.Print("Pokedex > ")
		if scanner.Scan() {
			words := cleanInput(scanner.Text())
			commandName := words[0]
			command, exists := getCommands()[commandName]
			if exists {
				err := command.callback(opts)
				if err != nil {
					fmt.Println(err)
				}
				continue
			} else {
				fmt.Println("Unknown command")
				continue
			}
		}		
	}
}

func cleanInput(text string) []string {
	result := strings.Split(strings.Trim(strings.ToLower(text)," ")," ")
	return result
}

func commandExit(option *config, parameter *parameters) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(option *config, parameter *parameters) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	commands := getCommands()
	for entry := range(commands) {
		fmt.Print(commands[entry].name, ": ", commands[entry].description, "\n")
	}
	return nil
}

func commandMap(option *config, parameter *parameters) error {
	if option.Next == "" {
		fmt.Println("you're on the last page")
		return nil
	}
	res, err := http.Get(option.Next)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return nil
	}
	if err != nil {
		return err
	}
	
	Response := response{}
	err = json.Unmarshal(body, &Response)
	if err != nil {
    	return err
	}
	option.Previous = option.Next //set the current one as previous
	option.Next = Response.Next
	
	for _, url := range(Response.Results) {
		fmt.Println(url.Name)
	}
	return nil
}

func commandMapb(option *config, parameter *parameters) error {
	if option.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	res, err := http.Get(option.Previous)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return nil
	}
	if err != nil {
		return err
	}
	
	Response := response{}
	err = json.Unmarshal(body, &Response)
	if err != nil {
    	return err
	}
	option.Next = option.Previous //set the current one as previous
	option.Previous = Response.Previous
	
	for _, url := range(Response.Results) {
		fmt.Println(url.Name)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
    "exit": {
        name:        "exit",
        description: "Exit the Pokedex",
        callback:    commandExit,
    },
	"help": {
        name:        "help",
        description: "Displays a help message",
        callback:    commandHelp,
    },
	"map": {
        name:        "map",
        description: "Lists the next 20 location areas",
        callback:    commandMap,
    },
	"mapb": {
        name:        "mapb",
        description: "Lists the previous 20 location areas",
        callback:    commandMapb,
    },
	"explore": {
        name:        "explore",
        description: "Lists all the pokemons on the mentioned location area (specified via id/name)",
        callback:    commandExplore,
    },
	}
	return commands
}