package main

import (
	"log"

	"github.com/IBM/sarama"
)

// Using asynchronous send, may be it's not a Java idiomatic way
// because we're using channels here.
func withAsyncProducer(config *sarama.Config, addr []string) {
	asyncProducer, err := sarama.NewAsyncProducer(addr, config)
	mayBeDie(err, "Error creating async producer: ")

	defer func() {
		e := asyncProducer.Close()
		mayBeDie(e, "Error closing async producer: ")
	}()

	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder("hello-key"),
		Value: sarama.StringEncoder("hello from async producer"),
	}

	asyncProducer.Input() <- msg

	select {
	case suc := <-asyncProducer.Successes():
		log.Printf(
			"Message is stored in topic(%s)/partition(%d)/offset(%d)\n",
			suc.Topic, suc.Partition, suc.Offset,
		)
	case fail := <-asyncProducer.Errors():
		log.Println(fail)
	}
}
