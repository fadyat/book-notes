package kafka_producer;

import io.confluent.common.utils.Utils;
import org.apache.kafka.clients.producer.Partitioner;
import org.apache.kafka.common.Cluster;
import org.apache.kafka.common.InvalidRecordException;
import org.apache.kafka.common.PartitionInfo;

import java.util.List;
import java.util.Map;
import java.util.Objects;

public class BananaPartitioner implements Partitioner {

    @Override
    public int partition(
            String topic,
            Object key,
            byte[] keyBytes,
            Object value,
            byte[] valueBytes,
            Cluster cluster
    ) throws InvalidRecordException {
        List<PartitionInfo> partitions = cluster.partitionsForTopic(topic);
        int numPartitions = partitions.size();

        if (Objects.isNull(key) || !(key instanceof String)) {
            throw new InvalidRecordException("All messages must have customer name as key");
        }

        if (key.equals("Banana") || numPartitions == 1) {
            return numPartitions - 1;
        }

        return Math.abs(Utils.murmur2(keyBytes)) % (numPartitions - 1);
    }

    @Override
    public void close() {
    }

    @Override
    public void configure(Map<String, ?> map) {
    }
}