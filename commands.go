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
  "math/rand"
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

    "catch": {
      name: "catch",
      description: "input 'catch <pokemon name>' to throw pokeball",
      callback: commandCatch,
    },

    "inspect": {
      name: "inspect",
      description: "check caught pokemon stats",
      callback: commandInspect,
    },

    "pokedex": {
      name: "pokedex",
      description: "open your pokedex",
      callback: commandPokedex,
    },
  }
}

func getResponseBody(url string) []byte {
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

  return body
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

func commandMap(cfg *config, cache *pokecache.Cache, parameter string, pokedex map[string]PokeStats) error {
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

  body := getResponseBody(url)
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

func commandMapb(cfg *config, cache *pokecache.Cache, parameter string, pokedex map[string]PokeStats) error {
  if cfg.Previous == nil {
    return fmt.Errorf("you're on the first page")
  } 

  url := *cfg.Previous
  body := getResponseBody(url)
  config1 := config{}
  jsonErr := json.Unmarshal(body, &config1)
  if jsonErr != nil {
    log.Fatal(jsonErr)
  }

  *cfg = config1
  printLocations(cfg.Results)

  return nil
}

func commandExit(cfg *config, cache *pokecache.Cache, parameter string, pokedex map[string]PokeStats) error {
  fmt.Println("Closing the Pokedex... Goodbye!")
  os.Exit(0)

  return nil
}

func commandHelp(cfg *config, cache *pokecache.Cache, parameter string, pokedex map[string]PokeStats) error {
  fmt.Println("Welcome to the Pokedex!")
  fmt.Println("Usage:")
  fmt.Println("")
  for _, value := range getCommands() {
    fmt.Printf("%s: %s\n", value.name, value.description)
  }
  return nil
}

func commandExplore(cfg *config, cache *pokecache.Cache, location string, pokedex map[string]PokeStats) error {
  fmt.Printf("Exploring %s...\n", location)
  url := "https://pokeapi.co/api/v2/location-area/" + location 

  body := getResponseBody(url)

  config1 := LocationArea{}
  jsonErr := json.Unmarshal(body, &config1)
  if jsonErr != nil {
    log.Fatal(jsonErr)
  }

  printPokemon(config1.PokemonEncounter)

  return nil
}

func commandCatch (cfg *config, cache *pokecache.Cache, name string, pokedex map[string]PokeStats) error {
  fmt.Printf("Throwing a Pokeball at %s...\n", name)
  url := "https://pokeapi.co/api/v2/pokemon/" + name

  body := getResponseBody(url)

  config1 := PokeStats{}
  jsonErr := json.Unmarshal(body, &config1)
  if jsonErr != nil {
    log.Fatal(jsonErr)
  }

  probablility := rand.Intn(config1.BaseExperience) % 10

  if probablility > 1 {
    fmt.Printf("%s was caught!\n", name)
    pokedex[name] = config1
    return nil
  }

  fmt.Printf("%s escaped!\n", name)
  
  return nil
}

func commandInspect(cfg *config, cache *pokecache.Cache, name string, pokedex map[string]PokeStats) error {
  pokestats, exists := pokedex[name]
  if !exists {
    fmt.Println("you have not caught that pokemon")
  }

  fmt.Println("Name: ", pokestats.Name)
  fmt.Println("Height: ", pokestats.Height)
  fmt.Println("Weight: ", pokestats.Weight)
  fmt.Println("Stats: ")
  for _, stat := range pokestats.Stats {
    fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
  }
  fmt.Println("Types:")
  for _, types := range pokestats.Types {
    fmt.Printf("  - %s\n", types.Type.Name)
  }
  fmt.Println("")

  return nil
} 

func commandPokedex(cfg *config, cache *pokecache.Cache, name string, pokedex map[string]PokeStats) error {
  if len(pokedex) == 0 {
    fmt.Println("You have not caught any Pokemon")
    return nil
  }
  fmt.Println("Your Pokedex:")
  for name, _ := range pokedex {
    fmt.Printf("  - %s\n", name)
  }

  return nil
}
