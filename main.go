package main

import (
  "time"
  "encoding/json"
  "bufio"
  "strings"
  "fmt"
  "os"
  "io"
  "log"
  "net/http"
  "github.com/adibbelel/pokedexcli/internal/pokecache"
)

type cliCommands struct {
  name string
  description string
  callback func(*config, *pokecache.Cache) error
}

type config struct {
  Next     *string `json:"next"`
  Previous *string `json:"previous"`
  Results []LocationArea `json:"results"`
}

type LocationArea struct {
  Name string `json:"name"`
  URL  string `json:"url"`
}


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
  }
}

func cleanInput(text string) []string {
  text = strings.TrimSpace(text)
	output := strings.Split(strings.ToLower(text), " ")
	return output
}

func commandExit(cfg *config, cache *pokecache.Cache) error {
  fmt.Print("Closing the Pokedex... Goodbye!")
  os.Exit(0)
  return nil
}

func commandHelp(cfg *config, cache *pokecache.Cache) error {
  fmt.Println("Welcome to the Pokedex!")
  fmt.Println("Usage:")
  fmt.Println("")
  for _, value := range getCommands() {
    fmt.Printf("%s: %s\n", value.name, value.description)
  }
  return nil
}

func printLocations (areas []LocationArea) {
  for _, area := range areas {
    fmt.Println(area.Name)
  }  
}

func commandMap(cfg *config, cache *pokecache.Cache) error {
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
    printLocations(cacheConfig.Results)
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
  printLocations(config1.Results)

  return nil
}

func commandMapb(cfg *config, cache *pokecache.Cache) error {
  if cfg == nil || cfg.Previous == nil {
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
  for _, area := range cfg.Results {
    fmt.Println(area.Name)
  }  

  return nil
}

func main() {
  cfg := &config{}
  cache := pokecache.NewCache(5 * time.Minute)
  go cache.ReapLoop(1 * time.Minute)
  scanner := bufio.NewScanner(os.Stdin)
  for { 
    fmt.Print("Pokedex > ")
    scanner.Scan()
    input := cleanInput(scanner.Text())
    
    if len(input) == 0 {
      continue
    }
    
    command := input[0]
    cmd, value := getCommands()[command]
    
    if value {
      err := cmd.callback(cfg, cache)
      if err != nil {
        fmt.Println("Error:", err)
      }

      continue
    } 
    
    fmt.Println("Unknown command")
    
  }
}

