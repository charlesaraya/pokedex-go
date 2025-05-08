package commands

import (
	"testing"
	"time"

	"github.com/charlesaraya/pokedex-go/internal/cache"
	"github.com/charlesaraya/pokedex-go/internal/pokedex"
)

func TestCommands(t *testing.T) {
	registry := GetRegistry()
	duration, _ := time.ParseDuration("1s")
	Cache := cache.NewCache(duration)
	Cache.Pokedex = pokedex.NewPokedex()

	t.Run("run whereami command", func(t *testing.T) {
		command := registry[CMD_WHEREAMI]
		if err := command.Command(command.Config, Cache); err != nil {
			t.Errorf("error %q command", CMD_WHEREAMI)
		}
	})

	t.Run("run explore command", func(t *testing.T) {
		command := registry[CMD_EXPLORE]
		if err := command.Command(command.Config, Cache); err != nil {
			t.Errorf("error %q command", CMD_EXPLORE)
		}
	})
}
