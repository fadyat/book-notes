package repo

type Message struct {

	// offset is the current index of the message in the partition.
	offset int64

	// content is the raw bytes of the message.
	content []byte
}

type Partition struct {

	// key is the unique identifier of the partition.
	// Used for messages distribution across partitions.
	key string

	// topic is the name of the topic that the partition belongs to.
	topic string

	// offset is the first available index in the partition.
	offset int64

	// messages is the FIFO queue of messages in the partition.
	// structured from the oldest to the newest.
	messages Queue[Message]
}

type Topic struct {

	// name is the unique user-defined identifier of the topic.
	name string

	// partitions is the separate areas, each of which can be managed independently.
	// in one time, only one consumer can read from a partition.
	// messages are distributed across partitions using a partition key or a round-robin algorithm.
	//
	// by default topic has exactly one partition, can be increased by the user.
	partitions map[string]*Partition
}

type Producer struct {

	// id is the unique identifier of the producer.
	id int64
}

type Consumer struct {

	// id is the unique identifier of the consumer.
	id int64
}
