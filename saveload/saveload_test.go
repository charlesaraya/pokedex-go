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
	dirPath := filepath.Join(".", TEST_DIR)

	t.Run("save pokedex", func(t *testing.T) {
		if err := Save(pokedex, TEST_DIR); err != nil {
			t.Errorf("error saving pokedex")
		}
		filePath := filepath.Join(dirPath, SAVE_FILE_NAME)
		_, err := os.ReadFile(filePath)
		if err != nil {
			t.Errorf("error save file  does not exist")
		}
	})

	t.Run("load pokedex", func(t *testing.T) {
		var got *pokeapi.Pokedex
		got, err = Load(TEST_DIR)
		if err != nil {
			t.Errorf("error loading test data")
		}
		pokemonNames := got.GetAll()
		if len(pokemonNames) == 0 {
			t.Errorf("error no pokemons in pokedex")
		}
		want := "pikachu"
		pokedexEntry, ok := got.Get(want)
		if !ok {
			t.Errorf("error got %v, want %v", pokedexEntry.Pokemon.Name, want)
		}
	})
	if err := os.RemoveAll(dirPath); err != nil {
		t.Errorf("error removing test data dir")
	}
}
