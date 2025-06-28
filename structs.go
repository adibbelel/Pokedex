package main

import (
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
