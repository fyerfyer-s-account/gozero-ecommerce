package serializer

// Serializer defines the interface for serializing/deserializing data
type Serializer interface {
    // Marshal serializes data into bytes
    Marshal(v interface{}) ([]byte, error)
    
    // Unmarshal deserializes data from bytes
    Unmarshal(data []byte, v interface{}) error
    
    // ContentType returns the content type for the serialization format
    ContentType() string
}