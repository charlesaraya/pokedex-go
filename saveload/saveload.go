package saveload

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charlesaraya/pokedex-go/pokeapi"
)

const (
	DATA_DIR       string = "data"
	TEST_DIR       string = "test_data"
	SAVE_FILE_NAME string = "pokedex.json"
)

func Save(p *pokeapi.Pokedex, dir string) error {
	// Ensure folder exists
	dirPath := filepath.Join(".", dir)
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error: creating folder %w", err)
	}

	// Create the JSON file
	filePath := filepath.Join(dirPath, SAVE_FILE_NAME)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error: creating file %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	if err := encoder.Encode(p); err != nil {
		return fmt.Errorf("error: encoding json %w", err)
	}
	return nil
}

func Load() (*pokeapi.Pokedex, error) {
	return nil, nil
}
