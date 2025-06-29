package main

import (
  "github.com/adibbelel/pokedexcli/internal/pokecache"
)

type cliCommands struct {
  name string
  description string
  callback func(*config, *pokecache.Cache, string, map[string]PokeStats) error
}

type config struct {
  Next     *string `json:"next"`
  Previous *string `json:"previous"`
  Results []LocationArea `json:"results"`
}

type LocationArea struct {
	Name                 string `json:"name"`
	PokemonEncounter []PokemonEncounters `json:"pokemon_encounters"`
}

type PokeStats struct {
  BaseExperience int `json:"base_experience"`
}


type PokemonEncounters struct {
  Pokemon struct {
    Name string `json:"name"`
  } `json:"pokemon"`
  VersionDetails []struct {
    Version struct {
      Name string `json:"name"`
      URL  string `json:"url"`
    } `json:"version"`
    MaxChance        int `json:"max_chance"`
    EncounterDetails []struct {
      MinLevel        int   `json:"min_level"`
      MaxLevel        int   `json:"max_level"`
      ConditionValues []any `json:"condition_values"`
      Chance          int   `json:"chance"`
      Method          struct {
        Name string `json:"name"`
        URL  string `json:"url"`
      } `json:"method"`
    } `json:"encounter_details"`
  } `json:"version_details"`
} 
