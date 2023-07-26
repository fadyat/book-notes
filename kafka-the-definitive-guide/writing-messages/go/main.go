package main

import (
	"encoding/binary"
	"github.com/IBM/sarama"
	"github.com/riferrei/srclient"
	"log"
	"os"
	"strings"
)

const (
	subjectNotFound         = 40401
	topicName               = "test"
	brokerURL               = "localhost:9094"
	schemaRegistryClientURL = "http://localhost:8081"
	customerAvroSchemaPath  = "./avro/customer.avsc"
)

func toSubject(topic string) string {
	if strings.HasSuffix(topic, "-value") {
		return topic
	}

	return topic + "-value"
}

func getAvroSchema() string {
	content, err := os.ReadFile(customerAvroSchemaPath)
	mayBeDie(err, "Error reading avro schema: ")

	return string(content)
}

func toSchemaRegistryFormat(schema *srclient.Schema, customer Customer) []byte {
	schemaInfo := binary.BigEndian.AppendUint32([]byte{0x0}, uint32(schema.ID()))

	customerBytes, err := schema.Codec().BinaryFromNative(schemaInfo, customer.toMap())
	mayBeDie(err, "Error converting customer to bytes: ")

	return customerBytes
}

func withAvroSerializer(config *sarama.Config, addr []string) {
	producer, err := sarama.NewSyncProducer(addr, config)
	mayBeDie(err, "Error creating sync producer: ")

	defer func() {
		e := producer.Close()
		mayBeDie(e, "Error closing sync producer: ")
	}()

	schemaRegistry := srclient.CreateSchemaRegistryClient(schemaRegistryClientURL)
	schema, err := schemaRegistry.GetLatestSchema(toSubject(topicName))
	e, ok := err.(srclient.Error)
	if !ok || e.Code != subjectNotFound {
		mayBeDie(err, "Error retrieving schema: ")
	}

	if schema == nil {
		log.Println("Schema not found, creating a new one")
		schema, err = schemaRegistry.CreateSchema(toSubject(topicName), getAvroSchema(), srclient.Avro)
		mayBeDie(err, "Error creating schema: ")
	}

	customer := Customer{ID: 1, Name: "John Doe"}
	customerBytes := toSchemaRegistryFormat(schema, customer)

	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder("Banana"),
		Value: sarama.ByteEncoder(customerBytes),
	}

	partition, offset, err := producer.SendMessage(msg)
	mayBeDie(err, "Error sending message: ")

	log.Printf(
		"Message is stored in topic(%s)/partition(%d)/offset(%d)\n",
		topicName, partition, offset,
	)
}

func main() {
	config := sarama.NewConfig()
	config.Producer.Partitioner = withBananaPartitioner

	// If we want to send message with fire and forget,
	// we can set `config.Producer.Return.Successes` to false.
	config.Producer.Return.Successes = true
	addr := []string{brokerURL}

	// Using an Avro encoded message with Schema Registry for
	// serialization/deserialization of the message payload.
	withAvroSerializer(config, addr)
}
