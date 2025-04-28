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

func TestKeys(t *testing.T) {
	cases := []struct {
		input    []byte
		expected string
	}{
		{
			input:    []byte{27, 91, 65}, // Up arrow
			expected: KEY_UP,
		},
		{
			input:    []byte{27, 91, 66}, // Down arrow
			expected: KEY_DOWN,
		},
		{
			input:    []byte{27, 91, 68}, // Left arrow
			expected: KEY_LEFT,
		},
		{
			input:    []byte{27, 91, 67}, // Right arrow
			expected: KEY_RIGHT,
		},
		{
			input:    []byte{127, 0, 0}, // Backspace
			expected: KEY_BACKSPACE,
		},
		{
			input:    []byte{10, 0, 0}, // Enter
			expected: KEY_ENTER,
		},
		{
			input:    []byte{32, 0, 0}, // Printable lower bound
			expected: KEY_PRINTABLE,
		},
		{
			input:    []byte{126, 0, 0}, // Printable upper bound
			expected: KEY_PRINTABLE,
		},
		{
			input:    []byte{128, 0, 0}, // Unknown
			expected: KEY_UNKNOWN,
		},
	}
	for _, c := range cases {
		if key := getKey(c.input); key != c.expected {
			t.Errorf("got %q expected %v", key, c.expected)
		}
	}
}
