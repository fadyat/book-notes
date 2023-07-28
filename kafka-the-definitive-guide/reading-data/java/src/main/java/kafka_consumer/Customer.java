package kafka_consumer;

import org.apache.avro.generic.GenericRecord;


public record Customer(Integer id, String name) {

    public static Customer fromGenericRecord(GenericRecord record) {

        return new Customer(
                (Integer) record.get("id"),
                record.get("name").toString()
        );
    }
}
