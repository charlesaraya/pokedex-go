package saveload

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/charlesaraya/pokedex-go/pokeapi"
)

func TestSaveGame(t *testing.T) {
	pokedex := pokeapi.NewPokedex()
	pikachu, err := pokeapi.GetPokemon(pokeapi.ENDPOINT_POKEMON + "pikachu")
	if err != nil {
		fmt.Printf("error: GetPokemon failed.")
	}
	pokedex.Add(pikachu)

	t.Run("save pokedex", func(t *testing.T) {
		if err := Save(pokedex, TEST_DIR); err != nil {
			t.Errorf("error saving pokedex")
		}
		dirPath := filepath.Join("..", TEST_DIR)
		filePath := filepath.Join(dirPath, SAVE_FILE_NAME)
		_, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("error save file  does not exist")
		}
		if err := os.RemoveAll(dirPath); err != nil {
			t.Errorf("error removing test data dir")
		}
	})

	/* t.Run("load pokedex", func(t *testing.T) {
		var got *pokeapi.Pokedex
		got, err = Load()
		want := "Pikachu"
		pokedexEntry, _ := got.Get("Pikachu")
		p := pokedexEntry.Pokemon.Name
		if p != want {
			t.Errorf("error: got %v, want %v", p, want)
		}
	}) */
}
