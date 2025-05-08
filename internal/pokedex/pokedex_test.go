package pokedex

import "testing"

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
