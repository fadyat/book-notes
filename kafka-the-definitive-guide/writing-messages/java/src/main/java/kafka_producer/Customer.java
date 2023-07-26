package kafka_producer;

import org.apache.avro.Schema;
import org.apache.avro.generic.GenericData;
import org.apache.avro.generic.GenericRecord;

public record Customer(Integer id, String name) {

    // hardcoded schema for simplicity, can be read from schema-registry
    // check golang code for complete example
    public GenericRecord toGenericRecord() {
        return new GenericData.Record(
                new Schema.Parser().parse("""
                            {
                                "type": "record",
                                "namespace": "kafka_producer",
                                "name": "Customer",
                                "fields": [
                                    {
                                        "name": "id",
                                        "type": "int"
                                    },
                                    {
                                        "name": "name",
                                        "type": "string"
                                    }
                                ]
                            }
                        """)) {{
            put("id", id);
            put("name", name);
        }};
    }
}
