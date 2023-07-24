package repo

// Storage is abstracted from partitions, all logic is handled by the broker.
type Storage interface {

	// Save saves a message to a topic and returns the offset of the message.
	Save(topic string, message []byte) (int64, error)

	// Get gets a message from a topic by reading from the latest offset.
	Get(topic string) ([]byte, error)

	// Explore gets a message from a topic by reading from a specific offset.
	// If the offset is -1, it will read from the latest offset.
	Explore(topic string, offset int64) ([]byte, error)

	// ResetOffset resets the offset of a topic to the latest offset.
	// This is useful when a consumer wants to miss some incorrect messages.
	ResetOffset(topic string) error
}
