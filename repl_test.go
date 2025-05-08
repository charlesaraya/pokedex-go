package main

import (
	"testing"

	"github.com/charlesaraya/pokedex-go/internal/api"
)

func TestPokedex(t *testing.T) {
	var pokedex = api.NewPokedex()

	cases := []struct {
		input    api.Pokemon
		expected string
	}{
		{
			input:    api.Pokemon{Name: "pikachu"},
			expected: "pikachu",
		},
		{
			input:    api.Pokemon{Name: "clefairy"},
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
