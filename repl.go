package main


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
