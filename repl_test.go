package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "   Hello   World   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "   ",
			expected: []string{},
		},
		{
			input:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. ",
			expected: []string{"lorem", "ipsum", "dolor", "sit", "amet,", "consectetur", "adipiscing", "elit."},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check lengths
		if len(actual) != len(c.expected) {
			t.Errorf("Error [cleanInput]: cleaned input slice length (%v) doesn't match the expected (%v).\n", len(actual), len(c.expected))
		}
		// Check cleaned words
		for i := range actual {
			if actual[i] != c.expected[i] {
				t.Errorf("Error [cleanInput]: cleaned input[%v] doesn't match the expected", i)
			}
		}
	}
}

func TestPokedex(t *testing.T) {
	var pokedex = NewPokedex()

	cases := []struct {
		input    Pokemon
		expected string
	}{
		{
			input:    Pokemon{Name: "pikachu"},
			expected: "pikachu",
		},
		{
			input:    Pokemon{Name: "clefairy"},
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
