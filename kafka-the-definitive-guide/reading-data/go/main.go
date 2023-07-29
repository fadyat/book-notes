package main

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

const (
	topic             = "test"
	schemaRegistryUrl = "http://localhost:8081"
	timeoutMs         = 100
)

func deserializeCustomer(msg *kafka.Message, d *avro.GenericDeserializer) {
	var customer Customer

	err := d.DeserializeInto(*msg.TopicPartition.Topic, msg.Value, &customer)
	if err != nil {
		log.Printf("Failed to deserialize customer: %s\n", err)
		return
	}

	log.Printf(
		"Message on %s[%d-%d]: %v\n",
		*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset,
		customer,
	)
}

func safeRebalance(c *kafka.Consumer, e kafka.Event) error {
	switch ev := e.(type) {
	case kafka.AssignedPartitions:
		log.Println(ev)
		if err := c.Assign(ev.Partitions); err != nil {
			return err
		}
	case kafka.RevokedPartitions:
		log.Println(ev)
		if c.AssignmentLost() {
			log.Printf("Lost partitions\n")
		}

		if err := syncCommit(c); err != nil {
			log.Printf("Failed to commit offsets: %v\n", err)
			return err
		}

	default:
		log.Printf("Ignored rebalance event: %s\n", ev)
	}

	return nil
}

func syncCommit(c *kafka.Consumer) error {
	offsets, err := c.Commit()
	if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {
		return err
	}

	log.Printf("Sync committed offsets: %v\n", offsets)
	return nil
}

func asyncCommit(c *kafka.Consumer, errChan chan<- error) {
	// emulate some network delay -> next commit can arrive before this one
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)

	offsets, err := c.Commit()
	if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {
		errChan <- err
		return
	}

	log.Printf("Async committed offsets: %v\n", offsets)
}

func main() {
	// confluent-kafka-go don't have any special API for asynchronous commit,
	// and batch-processing of messages.
	//
	// So, we have to implement it by ourselves.
	//
	// > batch-processing isn't done in this example

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9094",
		"group.id":           "myGroup",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	mbDie(err, "Failed to create consumer")
	defer func() {
		log.Println("Closing consumer")
		e := c.Close()
		mbDie(e, "Failed to close consumer")
	}()

	err = c.SubscribeTopics([]string{topic}, safeRebalance)
	mbDie(err, "Failed to subscribe to topics")

	schemaRegistry, err := schemaregistry.NewClient(schemaregistry.NewConfig(schemaRegistryUrl))
	mbDie(err, "Failed to create schema registry client")

	deserializer, err := avro.NewGenericDeserializer(schemaRegistry, serde.ValueSerde, avro.NewDeserializerConfig())
	mbDie(err, "Failed to create deserializer")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	commitErrChan := make(chan error)
	defer close(commitErrChan)

	go func() {
		for cerr := range commitErrChan {
			log.Printf("Failed to commit offsets: %v\n", cerr)
		}
	}()

	log.Println("Starting consumer loop")

poll:
	for {
		select {
		case sig := <-sigChan:
			log.Printf("Caught signal %v: terminating\n", sig)
			break poll
		default:
			switch e := c.Poll(timeoutMs).(type) {
			case *kafka.Message:
				deserializeCustomer(e, deserializer)
				go asyncCommit(c, commitErrChan)
			case *kafka.Error:
				log.Printf("Error: %v\n", e)
			case nil:
			default:
				log.Printf("Ignored %v\n", e)
			}
		}
	}

	log.Println("Trying to commit missed offsets")
	if err = syncCommit(c); err != nil {
		log.Printf("Failed to commit offsets: %v\n", err)
	}
}
