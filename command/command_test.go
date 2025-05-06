package command

import (
	"bytes"
	"testing"
	"time"

	"github.com/charlesaraya/pokedex-go/pokeapi"
)

type cacheData struct {
	key   string
	value []byte
}

func TestCache(t *testing.T) {
	var duration, _ = time.ParseDuration("0.5s")
	var cache = NewCache(duration)

	cases := []struct {
		input    cacheData
		expected cacheData
	}{
		{
			input: cacheData{
				key:   "map",
				value: []byte("content"),
			},
			expected: cacheData{
				key:   "map",
				value: []byte("content"),
			},
		},
	}

	for _, c := range cases {
		cache.Add(c.input.key, c.input.value)
		entry, ok := cache.Get(c.expected.key)
		if ok == false {
			t.Errorf("Error [Cache.Get]: couldn't Get cached entry for key: %v", c.input.key)
		}
		if bytes.Equal(entry.Val, c.expected.value) == false {
			t.Errorf("Error [Cache.Get]: cached entry value does not match expected value: %v != %v", entry.Val, c.expected.value)
		}
		// sleeps for 1 sec
		time.Sleep(2 * duration)

		_, ok = cache.Get(c.expected.key)
		if ok != false {
			t.Errorf("Error [Cache.reapLoop]: reapLoop should have cleaned the cache after duration %vs", duration)
		}
	}
}

func TestCommands(t *testing.T) {
	registry := GetRegistry()
	duration, _ := time.ParseDuration("1s")
	cache := NewCache(duration)
	cache.Pokedex = pokeapi.NewPokedex()

	t.Run("run whereami command", func(t *testing.T) {
		command := registry[CMD_WHEREAMI]
		if err := command.Command(command.Config, cache); err != nil {
			t.Errorf("error %q command", CMD_WHEREAMI)
		}
	})

	t.Run("run explore command", func(t *testing.T) {
		command := registry[CMD_EXPLORE]
		if err := command.Command(command.Config, cache); err != nil {
			t.Errorf("error %q command", CMD_EXPLORE)
		}
	})
}
