package helpers

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

// LoadFixture loads test data from JSON file
func LoadFixture(filename string, v interface{}) error {
    path := filepath.Join("../testdata", filename)
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read fixture file %s: %w", filename, err)
    }

    if err := json.Unmarshal(data, v); err != nil {
        return fmt.Errorf("failed to unmarshal fixture data: %w", err)
    }

    return nil
}

// SaveFixture saves test data to JSON file
func SaveFixture(filename string, v interface{}) error {
    data, err := json.MarshalIndent(v, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal fixture data: %w", err)
    }

    path := filepath.Join("testdata", filename)
    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("failed to write fixture file %s: %w", filename, err)
    }

    return nil
}