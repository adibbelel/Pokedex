package main

import (
  "bufio"
  "strings"
  "fmt"
  "os"
)

func main() {
  scanner := bufio.NewScanner(os.Stdin)
  for { 
    fmt.Print("Pokedex > ")
    scanner.Scan()
    input := scanner.Text()
    cleanText := cleanInput(input)
    fmt.Printf("Your command was: %s\n", cleanText[0])
  }
}

func cleanInput(text string) []string {
  text = strings.TrimSpace(text)
	output := strings.Split(strings.ToLower(text), " ")
	return output
}


