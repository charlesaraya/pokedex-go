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

### CLI Enhancements for Text Navigation and Editing
This enhancement allows for basic text manipulation in the terminal:

- **Move Cursor Left/Right**: Navigate the cursor within the input text using the left and right arrow keys.
- **Backspace**: Delete the character at the cursor position using the backspace key.

These features provide a more interactive and user-friendly experience when entering commands in the terminal, similar to standard text editing behavior.

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

## Improvement Ideas

- Enhance the CLI to support using the "up" arrow key to scroll through and reuse previous commands.
- Introduce Pokémon battles, allowing players to simulate fights between their caught Pokémon.
- Expand unit test coverage to ensure better reliability and stability.
- Refactor the codebase for better organization, readability, and testability.
- Implement a party system where players can manage a team of Pokémon that can gain experience and level up.
- Enable Pokémon evolution, allowing caught Pokémon to evolve after meeting certain conditions (e.g., time-based or level-based).
- Persist user data by saving the Pokédex to disk, enabling progress to be saved between sessions.
- Improve exploration by providing navigational choices (e.g., "left" or "right") rather than manually typing out area names.
- Add random wild Pokémon encounters while exploring different locations.
- Introduce different types of Poké Balls (e.g., Poké Balls, Great Balls, Ultra Balls) with varying catch rates to make capturing more strategic.