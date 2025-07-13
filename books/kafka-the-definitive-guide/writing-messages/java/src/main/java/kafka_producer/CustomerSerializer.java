package kafka_producer;

import org.apache.kafka.common.serialization.Serializer;

import java.nio.ByteBuffer;

public class CustomerSerializer implements Serializer<Customer> {

    @Override
    public byte[] serialize(String topic, Customer customer) {
        if (customer == null) {
            return null;
        }

        byte[] nameBytes = customer.name().getBytes();
        ByteBuffer buffer = ByteBuffer.allocate(4 + 4 + nameBytes.length);
        buffer.putInt(customer.id());
        buffer.putInt(nameBytes.length);
        buffer.put(nameBytes);

        return buffer.array();
    }
}
