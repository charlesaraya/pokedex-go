package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charlesaraya/pokedex-go/internal/pokedex"
)

const (
	DATA_DIR       string = "data"
	TEST_DIR       string = "test_data"
	SAVE_FILE_NAME string = "pokedex.json"
)

func Save(p *pokedex.Pokedex, dirName string) error {
	// Ensure folder exists
	dirPath := filepath.Join(".", dirName)
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

func Load(dirName string) (*pokedex.Pokedex, error) {
	var pokedex *pokedex.Pokedex
	filePath := filepath.Join(".", dirName, SAVE_FILE_NAME)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %w", err)
	}
	if err := json.Unmarshal(data, &pokedex); err != nil {
		return nil, fmt.Errorf("error: unmarshal operation failed: %w", err)
	}
	return pokedex, nil
}
