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
