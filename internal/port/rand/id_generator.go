package rand

// IDGenerator is a port for generating unique identifiers.
type IDGenerator interface {
	// Generate creates and returns a new unique identifier as a string.
	// Returns an error if the ID generation fails.
	Generate() (string, error)

	// GenerateWithPrefix creates and returns a new unique identifier with a prefix.
	// Format: prefix_uniqueID (e.g., "my-namespace_01HQRS4K2N3P4Q5R6S7T8V9WXY")
	// Returns an error if the ID generation fails.
	GenerateWithPrefix(prefix string) (string, error)
}
