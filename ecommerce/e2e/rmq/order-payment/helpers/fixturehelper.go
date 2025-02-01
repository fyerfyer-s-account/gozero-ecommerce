package helpers

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "runtime"
)

// FixtureHelper handles test fixture operations
type FixtureHelper struct {
    basePath string
}

// NewFixtureHelper creates a new fixture helper
func NewFixtureHelper() *FixtureHelper {
    // Get the path to the testdata directory relative to this file
    _, filename, _, _ := runtime.Caller(0)
    basePath := filepath.Join(filepath.Dir(filename), "..", "testdata")
    
    return &FixtureHelper{
        basePath: basePath,
    }
}

// LoadFixture loads test data from JSON file
func (h *FixtureHelper) LoadFixture(filename string, v interface{}) error {
    if v == nil {
        return fmt.Errorf("target variable cannot be nil")
    }

    path := filepath.Join(h.basePath, filename)
    
    // Verify file exists
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return fmt.Errorf("fixture file %s does not exist: %w", filename, err)
    }

    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read fixture file %s: %w", filename, err)
    }

    // Validate JSON
    if !json.Valid(data) {
        return fmt.Errorf("invalid JSON in fixture file %s", filename)
    }

    // Unmarshal data
    if err := json.Unmarshal(data, v); err != nil {
        return fmt.Errorf("failed to unmarshal fixture data from %s: %w", filename, err)
    }

    return nil
}

// SaveFixture saves test data to JSON file
func (h *FixtureHelper) SaveFixture(filename string, v interface{}) error {
    if v == nil {
        return fmt.Errorf("source data cannot be nil")
    }

    // Marshal with indentation for readability
    data, err := json.MarshalIndent(v, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal fixture data for %s: %w", filename, err)
    }

    path := filepath.Join(h.basePath, filename)

    // Ensure directory exists
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory for fixture %s: %w", filename, err)
    }

    // Write file
    if err := os.WriteFile(path, data, 0644); err != nil {
        return fmt.Errorf("failed to write fixture file %s: %w", filename, err)
    }

    return nil
}

// LoadFixtureString loads fixture data as a string
func (h *FixtureHelper) LoadFixtureString(filename string) (string, error) {
    path := filepath.Join(h.basePath, filename)
    
    data, err := os.ReadFile(path)
    if err != nil {
        return "", fmt.Errorf("failed to read fixture file %s: %w", filename, err)
    }
    
    return string(data), nil
}

// ValidateFixture checks if the fixture file exists and contains valid JSON
func (h *FixtureHelper) ValidateFixture(filename string) error {
    path := filepath.Join(h.basePath, filename)
    
    // Check file exists
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return fmt.Errorf("fixture file %s does not exist", filename)
    }
    
    // Read and validate JSON
    data, err := os.ReadFile(path)
    if err != nil {
        return fmt.Errorf("failed to read fixture file %s: %w", filename, err)
    }
    
    if !json.Valid(data) {
        return fmt.Errorf("invalid JSON in fixture file %s", filename)
    }
    
    return nil
}

// For backward compatibility
func LoadFixture(filename string, v interface{}) error {
    helper := NewFixtureHelper()
    return helper.LoadFixture(filename, v)
}

func SaveFixture(filename string, v interface{}) error {
    helper := NewFixtureHelper()
    return helper.SaveFixture(filename, v)
}