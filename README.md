# CLI Pokédex
A Pokédex built as a command-line REPL! This project fetches real Pokémon data using the PokéAPI, allowing users to explore, catch, and inspect Pokémon — all through simple commands in your terminal.

## Features

### Explore the Pokémon World
- `map`: Displays the names of 20 location areas from the Pokémon world. Each call shows the next 20.
- `mapb`: (map back) Displays the previous 20 locations.

### Location Exploration
- `explore <location>`: Lists all the Pokémon found in a given location area.

### Catch Pokémon
 - `catch <pokemon>`: Attempts to catch a Pokémon by name using a simulated Pokéball throw. Successful catches add the Pokémon to your personal Pokédex.

### Inspect Your Pokémon
- `inspect <pokemon>`: View details (name, height, weight, stats, types) for any Pokémon you've successfully caught.

### Personal Pokédex
- `pokedex`: Lists all Pokémon you have caught so far.

### Help and Exit Commands
- `help`: Displays instructions and a list of available commands.
- `exit`: Safely exits the REPL.

### Caching for Speed
Responses from the PokéAPI are cached for faster access. Ensure safe concurrent access. Old cache entries are cleaned automatically using a Ticker-based system.

## Commands Reference

| Command              | Description                         |
|----------------------|-------------------------------------|
| `help`               | Show available commands            |
| `exit`               | Exit the REPL                      |
| `map`                | View the next 20 location areas    |
| `mapb`               | View the previous 20 location areas |
| `explore <location>` | List Pokémon in a given location   |
| `catch <pokemon>`    | Try to catch a Pokémon             |
| `inspect <pokemon>`  | View details about a caught Pokémon |
| `pokedex`            | List all caught Pokémon            |