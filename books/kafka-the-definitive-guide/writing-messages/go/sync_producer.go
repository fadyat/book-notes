package main

import (
	"log"

	"github.com/IBM/sarama"
)

// Using synchronous send
func withSyncProducer(config *sarama.Config, addr []string) {
	syncProducer, err := sarama.NewSyncProducer(addr, config)
	mayBeDie(err, "Error creating sync producer: ")

	defer func() {
		e := syncProducer.Close()
		mayBeDie(e, "Error closing sync producer: ")
	}()

	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder("hello-key"),
		Value: sarama.StringEncoder("hello from sync producer"),
	}

	partition, offset, err := syncProducer.SendMessage(msg)
	mayBeDie(err, "Error sending message: ")

	log.Printf(
		"Message is stored in topic(%s)/partition(%d)/offset(%d)\n",
		topicName, partition, offset,
	)
}
