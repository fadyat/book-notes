# Part 2. Distributed Data

## Chapter 6. Partitioning

Normally, partitions are defined in such a way that each piece of data (each record, row, or document) belongs to exactly one partition.

 In effect, each partition is a small database of its own, although the database may support operations that touch multi‐ ple partitions at the same time.

 Main reason - scalability.
 > Different partitions can be placed on different nodes. 

For queries that operate on a single partition, each node can independently execute the queries for its own partition, so query throughput can be scaled by adding more nodes.

### Partitioning and Replication 
