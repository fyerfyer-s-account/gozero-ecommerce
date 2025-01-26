package json

import (
    "encoding/json"
    "fmt"
)

const JsonContentType = "application/json"

type JsonSerializer struct{}

// New creates a new JSON serializer
func New() *JsonSerializer {
    return &JsonSerializer{}
}

func (s *JsonSerializer) Marshal(v interface{}) ([]byte, error) {
    if v == nil {
        return nil, fmt.Errorf("cannot marshal nil value")
    }
    
    data, err := json.Marshal(v)
    if err != nil {
        return nil, fmt.Errorf("json marshal error: %w", err)
    }
    
    return data, nil
}

func (s *JsonSerializer) Unmarshal(data []byte, v interface{}) error {
    if len(data) == 0 {
        return fmt.Errorf("cannot unmarshal empty data")
    }
    
    if v == nil {
        return fmt.Errorf("cannot unmarshal into nil value")
    }
    
    if err := json.Unmarshal(data, v); err != nil {
        return fmt.Errorf("json unmarshal error: %w", err)
    }
    
    return nil
}

func (s *JsonSerializer) ContentType() string {
    return JsonContentType
}

// Helper methods for common operations

func Marshal(v interface{}) ([]byte, error) {
    return New().Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
    return New().Unmarshal(data, v)
}

func MustMarshal(v interface{}) []byte {
    data, err := Marshal(v)
    if err != nil {
        panic(err)
    }
    return data
}

func MustUnmarshal(data []byte, v interface{}) {
    if err := Unmarshal(data, v); err != nil {
        panic(err)
    }
}