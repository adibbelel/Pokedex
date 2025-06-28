package main

import (
  "strings"
  "fmt"
  "io"
  "net/http"
  "github.com/adibbelel/pokedexcli/internal/pokecache"
  "encoding/json"
  "os"
  "log"
)

func getCommands() map[string]cliCommands{
   return map[string]cliCommands{
    "exit": {
      name: "exit",
      description: "Exit the Pokedex",
      callback: commandExit,
    },

    "help": {
      name: "help",
      description: "Displays a help message",
      callback: commandHelp,
    },

    "map": {
      name: "map",
      description: "Displays Pokemon world locations",
      callback: commandMap,
    },

    "mapb": {
      name: "mapb",
      description: "Goes back to previous world locations",
      callback: commandMapb,
    },

    "explore": {
      name: "explore",
      description: "input 'explore <location name>' to explore entered location",
      callback: commandExplore,
    },
  }
}

func cleanInput(text string) []string {
  text = strings.TrimSpace(text)
	output := strings.Split(strings.ToLower(text), " ")
	return output
}

func printPokemon (PokemonEncounter []PokemonEncounters) {
  if len(PokemonEncounter) == 0 {
    fmt.Println("No Pokemon found in this location")
  }

  fmt.Println("Found Pokemon:")
  for _, pokemons := range PokemonEncounter {
    fmt.Println("  - ", pokemons.Pokemon.Name)
  }
}

func printLocations (areas []LocationArea) {
  for _, area := range areas {
    fmt.Println(area.Name)
  }  
}

func commandMap(cfg *config, cache *pokecache.Cache, parameter string) error {
  url := "https://pokeapi.co/api/v2/location-area"

  if cfg != nil && cfg.Next != nil {
    url = *cfg.Next
  }

  if cached, ok := cache.Get(url); ok {
    var cacheConfig config
    jsonErr := json.Unmarshal(cached, &cacheConfig)
    if jsonErr != nil {
      log.Fatal(jsonErr)
    }

    *cfg = cacheConfig
    printLocations(cfg.Results)
  }

  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }

  body, err := io.ReadAll(res.Body)
  res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
  cache.Add(url, body)

  config1 := config{}
  jsonErr := json.Unmarshal(body, &config1)
  if jsonErr != nil {
    log.Fatal(jsonErr)
  }

  *cfg = config1
  printLocations(cfg.Results)

  return nil
}

func commandMapb(cfg *config, cache *pokecache.Cache, parameter string) error {

  if cfg.Previous == nil {
    return fmt.Errorf("you're on the first page")
  } 
  res, err := http.Get(*cfg.Previous)
  if err != nil {
    log.Fatal(err)
  }

  body, err := io.ReadAll(res.Body)
  res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

  config1 := config{}
  jsonErr := json.Unmarshal(body, &config1)
  if jsonErr != nil {
    log.Fatal(jsonErr)
  }

  *cfg = config1
  printLocations(cfg.Results)

  return nil
}

func commandExit(cfg *config, cache *pokecache.Cache, parameter string) error {
  fmt.Print("Closing the Pokedex... Goodbye!")
  os.Exit(0)
  return nil
}

func commandHelp(cfg *config, cache *pokecache.Cache, parameter string) error {
  fmt.Println("Welcome to the Pokedex!")
  fmt.Println("Usage:")
  fmt.Println("")
  for _, value := range getCommands() {
    fmt.Printf("%s: %s\n", value.name, value.description)
  }
  return nil
}

func commandExplore(cfg *config, cache *pokecache.Cache, location string) error {
  fmt.Printf("Exploring %s...\n", location)
  url := "https://pokeapi.co/api/v2/location-area/" + location 

  res, err := http.Get(url)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

  config1 := LocationArea{}
  jsonErr := json.Unmarshal(body, &config1)
  if jsonErr != nil {
    log.Fatal(jsonErr)
  }

  printPokemon(config1.PokemonEncounter)

  return nil
}
