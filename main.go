package main

import (
  "bufio"
  "strings"
  "fmt"
  "os"
)

type cliCommands struct {
  name string
  description string
  callback func() error
}

var commands map[string]cliCommands

func init() {
   commands = map[string]cliCommands{
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
  }
}

func cleanInput(text string) []string {
  text = strings.TrimSpace(text)
	output := strings.Split(strings.ToLower(text), " ")
	return output
}

func commandExit() error {
  fmt.Print("Closing the Pokedex... Goodbye!")
  os.Exit(0)
  return nil
}

func commandHelp() error {
  fmt.Println("Welcome to the Pokedex!")
  fmt.Println("Usage:")
  fmt.Println("")
  for _, value := range commands {
    fmt.Printf("%s: %s\n", value.name, value.description)
  }
  return nil
}

func main() {
 scanner := bufio.NewScanner(os.Stdin)
  for { 
    fmt.Print("Pokedex > ")
    scanner.Scan()
    input := cleanInput(scanner.Text())
    
    if len(input) == 0 {
      continue
    }
    
    command := input[0]
    cmd, value := commands[command]
    
    if value {
      err := cmd.callback()
      if err != nil {
        fmt.Println("Error:", err)
      }

      continue
    } 
    
    fmt.Println("Unknown command")
    
  }
}

 


