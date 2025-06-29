package main

import (
  "time"
  "bufio"
  "fmt"
  "os"
  "github.com/adibbelel/pokedexcli/internal/pokecache"
)

func main() {
  Pokedex := make(map[string]PokeStats)
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
    var parameter string
    if len(input) > 1 {
      parameter = input[1]
    }
    cmd, value := getCommands()[command]
    
    if value {
      err := cmd.callback(cfg, cache, parameter, Pokedex)
      if err != nil {
        fmt.Println("Error:", err)
      }

      continue
    } 
    
    fmt.Println("Unknown command")
    
  }
}

