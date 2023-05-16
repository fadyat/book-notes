# Part 2. Distributed Data

## Chapter 6. Partitioning

Normally, partitions are defined in such a way that each piece of data (each record, row, or document) belongs to exactly one partition.

 In effect, each partition is a small database of its own, although the database may support operations that touch multi‐ ple partitions at the same time.

 Main reason - scalability.
 > Different partitions can be placed on different nodes. 

For queries that operate on a single partition, each node can independently execute the queries for its own partition, so query throughput can be scaled by adding more nodes.

### Partitioning and Replication 

Partitioning is usually combined with replication so that copies of each partition are stored 
on multiple nodes.
> This means that, even though each record belongs to exactly one partition, it may be stored on multiple nodes for fault tolerance. 

### Partitioning of Key-Value Data 

How do you decide which records to store on which nodes?

Our goal with partitioning is to spread the data and query load evenly across nodes. 
> If share is fair, then in theory 10 nodes should be able to handle 10 times as much data 
and 10 times the read and write throughput of a single node (ignoring replication).
>
> If share is unfair, one partition have more data or queries than others we can call it 
**skewed**. 

A partition with disproportionately high load is called a **hot spot**. 

The simplest approach is to assign records to partitions randomly.
> This will distribute data equally, but when you're trying to read an item, you 
don't have understanding where it is.

We can do better, let's have a simple key-value data model, which is always access 
a record by its primary key.

#### Partitioning by Key Range 

Continuos range of keys (from min to max) to each partition. Good when you know 
the boundaries of the keys.

- The range of keys are not necessarily contiguous, but they are non-overlapping.
- Boundaries might be chosen manually by admin or the database can choose them automatically. 
- Within partition we can keep keys in sorted order. 
- Certain access patterns can lead to hot spots. 

#### Partitioning by Hash of Key 

Each partition is responsible for a range of hashes, and a record whose key hashes to that range is stored on that partition. 

- Good for distributing keys fairly among the partitions.
- The partition boundaries can be evenly spaced, or they can be chosen randomly.
- Losing the abitily to do efficient range queries.
> Cassandra achieves a compromise between the two strategies by using a **compound pk** consisting
of several columns. Only first part of key is hashed, another used for concatenated index for sorting. 

#### Skewed Workloads and Relieving Hot Spots 

Hash still may lead to hot spots.
> For example we read and write the same key over and over again. 

If the one key is known to be very hot, a simple technique is to add a random number 
at the beginning or end of the key.

However, having split of the writes accross different keys, any reads now have to 
do additional work - read the data from all keys and merge them.

### Partitioning and Secondary Indexes

The situation becomes more complex when we want to use secondary indexes. He doesn't 
identify the record, but speed up the search.

The problem with secondary indexes that they don't map neatly to partitions.

#### Partitioning Secondary Indexes by Document 

If you declared the index the database can perform the indexing automatically. 
Each partition has its own secondary index, and it can be updated independently.
For that reason it's called **local index**.

Reading from a document-partitioned index requires care.
> For example red car appears in both partitions, so we need to send query to both partitions and combine results.

This approach to quering a partitioned database is sometimes called **scatter/gather**, and it can make read 
queries on secondary indexes quity expensive. Even if you make requests in parallel. 

Most database vendors recommend that your structure your partitioning scheme so that secondary 
index queries can be served from a single partition, but it not always possible,
especially when you're using multiple secondary indexes in a single query. 


#### Partitioning Secondary Indexes by Term 

Rather than each partition having its own secondary index, we can 
construct a **global index** that covers all partitions. 

We can't store it in one place, it also needs to be partitioned :) differently from 
primary index.

We can call this kind of index **term-partitioned**, because the term we're looking 
fro determines th partition of the index. 

Advantage over document-partitioned index is that we can make read more efficient.
But writes are slower and complicated, because now we affect multiple partitions of the index. 

However term-partitioned indexes requires a distributed transaction, which is not
available in many databases.

In practice, updates of secondary indexes are often implemented asynchronously, 
may not be immediately consistent with the primary index.

### Rebalancing Partitions 


