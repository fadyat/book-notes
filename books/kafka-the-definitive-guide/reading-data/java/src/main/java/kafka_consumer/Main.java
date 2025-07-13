package kafka_consumer;

import io.confluent.kafka.serializers.KafkaAvroDeserializer;
import io.confluent.kafka.serializers.KafkaAvroDeserializerConfig;
import org.apache.avro.generic.GenericRecord;
import org.apache.kafka.clients.consumer.*;
import org.apache.kafka.common.TopicPartition;
import org.apache.kafka.common.errors.WakeupException;
import org.apache.kafka.common.serialization.StringDeserializer;

import java.time.Duration;
import java.util.*;


public class Main {
    private static final String TOPIC = "test";
    private static final String BOOTSTRAP_SERVERS = "localhost:9094";
    private static final String GROUP_ID = "test-group";
    private static final Map<TopicPartition, OffsetAndMetadata> currentOffsets = new HashMap<>();

    private record CustomRebalanceListener(KafkaConsumer<?, ?> consumer) implements ConsumerRebalanceListener {

        @Override
        public void onPartitionsAssigned(Collection<TopicPartition> partitions) {
            System.out.println("Partitions assigned: " + partitions);
        }

        @Override
        public void onPartitionsRevoked(Collection<TopicPartition> partitions) {
            System.out.println(
                    "Lost partitions in rebalance. Committing current offsets: " + currentOffsets
            );

            consumer.commitSync(currentOffsets);
        }
    }

    public static <K, V> KafkaConsumer<K, V> createConsumer(String keyDeserializerClass, String valueDeserializerClass) {
        var props = new Properties() {{
            put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, BOOTSTRAP_SERVERS);
            put(ConsumerConfig.GROUP_ID_CONFIG, GROUP_ID);
            put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, keyDeserializerClass);
            put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, valueDeserializerClass);
            put(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, false);
            put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");
            put(ConsumerConfig.MAX_POLL_RECORDS_CONFIG, 10);
            put("schema.registry.url", "http://localhost:8081");
        }};

        return new KafkaConsumer<>(props);
    }

    public static void registerShutdownHook(
            KafkaConsumer<?, ?> consumer, Thread mainThread
    ) {
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.out.println("Starting exit...");
            consumer.wakeup();

            try {
                mainThread.join();
            } catch (InterruptedException e) {
                System.out.println("Interrupted exception: " + e.getMessage());
            }
        }));
    }

    public static void main(String[] args) {
        KafkaConsumer<String, GenericRecord> consumer = createConsumer(
                StringDeserializer.class.getName(),
                KafkaAvroDeserializer.class.getName()
        );


        var mainThread = Thread.currentThread();
        registerShutdownHook(consumer, mainThread);

        var rebalanceListener = new CustomRebalanceListener(consumer);
        try {
            var timeoutMs = 1000;
            consumer.subscribe(Collections.singletonList(TOPIC), rebalanceListener);

            System.out.println("Subscribed to topic " + TOPIC);
            while (true) {
                var records = consumer.poll(Duration.ofMillis(timeoutMs));

                for (var record : records) {
                    var customer = Customer.fromGenericRecord(record.value());
                    System.out.printf(
                            "topic = %s, partition = %s, offset = %d, key = %s, id = %s, name = %s%n",
                            record.topic(), record.partition(), record.offset(), record.key(),
                            customer.id(), customer.name()
                    );

                    currentOffsets.put(
                            new TopicPartition(record.topic(), record.partition()),
                            new OffsetAndMetadata(record.offset() + 1, "no metadata")
                    );
                }

                if (records.isEmpty()) {
                    continue;
                }

                System.out.printf("Committing %d records%n", records.count());
                consumer.commitAsync(currentOffsets, null);
            }
        } catch (WakeupException ignored) {
        } catch (Exception e) {
            System.out.printf("Unexpected error: %s%n", e.getMessage());
        } finally {
            try {
                consumer.commitSync(currentOffsets);
            } finally {
                consumer.close();
            }
        }
    }
}