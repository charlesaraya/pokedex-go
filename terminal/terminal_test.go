package terminal

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
		actual := CleanInput(c.input)
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
		if key := GetKey(c.input); key != c.expected {
			t.Errorf("got %q expected %v", key, c.expected)
		}
	}
}

func TestBuffer(t *testing.T) {
	t.Run("init buffer", func(t *testing.T) {
		gotBuffer, gotCursor := InitBuffer()
		wantBuffer, wantCursor := []rune{}, 0
		if string(gotBuffer) != string(wantBuffer) {
			t.Errorf("got %q want %q", gotBuffer, wantBuffer)
		}
		if gotCursor != wantCursor {
			t.Errorf("got %q want %q", gotCursor, wantCursor)
		}
	})

	t.Run("update buffer", func(t *testing.T) {
		gotBuffer := []rune("hell")
		gotCursor := 4
		AddToBuffer(rune('o'), &gotBuffer, &gotCursor)
		wantBuffer := []rune("hello")
		wantCursor := 5
		if string(gotBuffer) != string(wantBuffer) {
			t.Errorf("got %q want %q", gotBuffer, wantBuffer)
		}
		if gotCursor != wantCursor {
			t.Errorf("got %q want %q", gotCursor, wantCursor)
		}
	})
	t.Run("reset buffer", func(t *testing.T) {
		gotBuffer := []rune("hell")
		gotCursor := 4
		ResetBuffer(&gotBuffer, &gotCursor)
		wantBuffer := []rune{}
		wantCursor := 0
		if string(gotBuffer) != string(wantBuffer) {
			t.Errorf("got %q want %q", gotBuffer, wantBuffer)
		}
		if gotCursor != wantCursor {
			t.Errorf("got %q want %q", gotCursor, wantCursor)
		}
	})
}
