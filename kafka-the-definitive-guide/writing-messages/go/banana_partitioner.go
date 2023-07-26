package main

import (
	"github.com/IBM/sarama"
	"math/rand"
	"time"
)

type bananaPartitioner struct {
	generator *rand.Rand
}

func (b *bananaPartitioner) Partition(msg *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	key, err := msg.Key.Encode()
	if err != nil {
		return 0, err
	}

	if string(key) == "Banana" || numPartitions == 1 {
		return numPartitions - 1, nil
	}

	return int32(b.generator.Intn(int(numPartitions - 1))), nil
}

func (b *bananaPartitioner) RequiresConsistency() bool {
	return false
}

func withBananaPartitioner(topic string) sarama.Partitioner {
	return &bananaPartitioner{
		generator: rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
	}
}
