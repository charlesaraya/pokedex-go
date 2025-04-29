package main

import (
	"testing"

	"github.com/charlesaraya/pokedex-go/pokeapi"
)

func TestPokedex(t *testing.T) {
	var pokedex = pokeapi.NewPokedex()

	cases := []struct {
		input    pokeapi.Pokemon
		expected string
	}{
		{
			input:    pokeapi.Pokemon{Name: "pikachu"},
			expected: "pikachu",
		},
		{
			input:    pokeapi.Pokemon{Name: "clefairy"},
			expected: "clefairy",
		},
	}

	for _, c := range cases {
		if _, ok := pokedex.Get(c.expected); ok {
			t.Errorf("Error [Pokedex.Get]: shouldn't have got any entry: %s", c.expected)
		}
		pokedex.Add(c.input)
		pokedexEntry, ok := pokedex.Get(c.expected)
		if !ok {
			t.Errorf("Error [Pokedex.Get]: Failed to Get a pokedex entry for key: %s", c.expected)
		}
		if pokedexEntry.Pokemon.Name != c.expected {
			t.Errorf("Error [Pokedex.Get]: Pokedex entry does not match expected Pokemon name: %s != %s", pokedexEntry.Pokemon.Name, c.expected)
		}
	}
}
