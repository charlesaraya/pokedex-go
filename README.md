# PokédexGo CLI

A Pokédex built as a command-line REPL! This project fetches real Pokémon data using the PokéAPI, allowing users to explore, catch, and inspect Pokémon — all through simple commands in your terminal.

## Motivation

PokédexGo was born from the idea of turning a simple programming challenge into a hands-on experience in building interactive command-line applications with Go. It goes far beyond just fetching data from an API. It simulates a full Pokédex experience where Pokémon trainers can explore locations, encounter and catch Pokémon, and manage their collection across sessions. 

The goal was to create something practical, nostalgic, and technically challenging with a fun, engaging context.

## Getting Started

### Installing

This assumes you already have a working Go environment, if not please see [this](https://go.dev/doc/install) page first.

```bash
go install github.com/charlesaraya/pokedex-go/cmd/pokedex@latest
```

### Run

```bash
pokedex
```

## Features

### Explore the Pokémon World

- `map`: Displays the names of 20 location areas from the Pokémon world. Each call shows the next 20.
- `mapb`: (map back) Displays the previous 20 locations.
- `visit <location-area>`: Visits a given location area in the Pokémon world.
- `explore <location-area>`: Lists all Pokémon that live in a given location area.

### Encounter and Catch Pokémon

 - `encounter`: Encounters a random Pokémon in the area based on their encounter chance. Call `catch`right away (without Pokémon name) before it escapes.
 - `catch [<pokemon>]`: Attempts to catch a Pokémon by name using a simulated Pokéball throw. Successful catches will add the Pokémon to your personal Pokédex.

### Personal Pokédex and Inspect Your Pokémon

- `pokedex`: Lists all Pokémon you have caught so far.
- `inspect <pokemon>`: View details (name, height, weight, stats, types) for any Pokémon you've successfully caught.

### Save your Progress

Your caught Pokémon are now saved to disk and automatically loaded on startup, allowing you to continue where you left off across sessions. No more starting over—your journey is saved!

- `save`: Saves your Pokedex and all caught Pokémon you have caught so far.
- `load`: Loads the last saved Pokedex to resume the game.

### Help and Exit Commands

- `help`: Displays instructions and a list of available commands.
- `exit`: Safely exits the REPL.

### Caching for Speed

Responses from the PokéAPI are cached for faster access. Ensure safe concurrent access. Old cache entries are cleaned automatically using a Ticker-based system.

### CLI Enhancements for Text Navigation and Editing

This enhancement allows for basic text manipulation in the terminal:

- **Move Cursor Left/Right**: Navigate the cursor within the input text using the left and right arrow keys.
- **Backspace**: Delete the character at the cursor position using the backspace key.
- **Move Cursor Up/Down**: Navigate the command history using the up and down arrow keys.

These features provide a more interactive and user-friendly experience when entering commands in the terminal, similar to standard text editing behavior.

## Commands Reference

| Command                | Description                         |
|------------------------|-------------------------------------|
| `help`                 | Show available commands             |
| `exit`                 | Exit the REPL                       |
| `whereami [-l \| -r]`  | Shows your current location area, location (`-l`) or region (`-r`).     |
| `map`                  | View the next 20 location areas     |
| `mapb`                 | View the previous 20 location areas |
| `visit <location>`     | Visit an existing location area     |
| `explore <location>`   | List Pokémon in a given location    |
| `encounter`            | Encounters a Pokémon in the area    |
| `catch [<pokemon>]`    | Try to catch a Pokémon              |
| `inspect <pokemon>`    | View details about a caught Pokémon |
| `pokedex`              | List all caught Pokémon             |
| `save`                 | Save your Pokedex                   |
| `load`                 | Load your latest saved Pokedex      |

## Improvement Ideas

- Introduce Pokémon battles, allowing players to simulate fights between their caught Pokémon.
- Implement a party system where players can manage a team of Pokémon that can gain experience and level up.
- Enable Pokémon evolution, allowing caught Pokémon to evolve after meeting certain conditions (e.g., time-based or level-based).
- Improve exploration by providing navigational choices (e.g., "left" or "right") rather than manually typing out area names.
- Introduce different types of Poké Balls (e.g., Poké Balls, Great Balls, Ultra Balls) with varying catch rates to make capturing more strategic.

## Contributing

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request to the `main` branch.
