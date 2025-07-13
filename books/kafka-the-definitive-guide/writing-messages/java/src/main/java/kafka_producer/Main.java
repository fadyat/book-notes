package kafka_producer;

import org.apache.avro.generic.GenericRecord;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.clients.producer.RecordMetadata;

import java.util.Properties;

public class Main {
    private static final String TOPIC = "test";

    public static void main(String[] args) {
        var producerProps = new Properties() {{
            put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, "localhost:9094");
            put("schema.registry.url", "http://localhost:8081");
            put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, "org.apache.kafka.common.serialization.StringSerializer");
            put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, "io.confluent.kafka.serializers.KafkaAvroSerializer");
            put(ProducerConfig.PARTITIONER_CLASS_CONFIG, "kafka_producer.BananaPartitioner");
        }};

        var customer = new Customer(1, "Ryan Gosling");

        RecordMetadata metadata;
        try (var producer = new KafkaProducer<String, GenericRecord>(producerProps)) {
            var record = new ProducerRecord<>(TOPIC, "Banana", customer.toGenericRecord());
            metadata = producer.send(record).get();
        } catch (Exception e) {
            e.printStackTrace();
            return;
        }

        System.out.println("Record sent to partition " + metadata.partition() + " with offset " + metadata.offset());
    }
}
