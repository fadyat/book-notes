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

Over time, things change in a database:

- the query throughput increase -> want to add more cpu to handle it 
- the dataset increase -> want to add more disk and RAM to store it 
- a machine fail -> others need to take over its work 

All of these changes call for data and requests to be moved from one node to another. 
The process of moving load from one node to another is called **rebalancing**.

No matter which scheme is used, rebalancing is usually expected to meet some min requirements:

- after rebalancing, the load should be shared fairly between the nodes of the cluster 
- while rebalancing is happening, the database should continue accepting reads and writes 
- no more data than necessary should be moved between nodes, to make it fast and minimize network and disk i/o load. 

#### Strategies for Rebalancing 

There are a few different strategies for rebalancing data between nodes:

##### How not to do it: hash mod N

We routing each record to partiontion by making a ranges of hashes:
> e.g 0..10, 10..20 etc. 

Not by mod, because if we add a new node we need to move all previous data to a new 
range.
> e.g was n % 10, 11 -> 1; now n % 11, 11 -> 0 (need to be moved)

We need an approach that doesn’t move data around more than necessary.

##### Fixed number of partitions 

Fortunately, there is a fairly simple solution: create many more partitions than there are nodes, and assign several partitions to each node. For example, a database run‐ ning on a cluster of 10 nodes may be split into 1,000 partitions from the outset so that approximately 100 partitions are assigned to each node.

Now, if a node is added to the cluster, the new node can steal a few partitions from every existing node until partitions are fairly distributed once again.

Only entire partitions are moved between nodes. The number of partitions does not change, nor does the assignment of keys to partitions. The only thing that changes is the assignment of partitions to nodes. 

In principle, you can even account for mismatched hardware in your cluster: by assigning more partitions to nodes that are more powerful, you can force those nodes to take a greater share of the load.

Choosing the right number of partitions is difficult if the total size of the dataset is highly variable

The best performance is achieved when the size of partitions is “just right,” neither too big nor too small, which can be hard to achieve if the number of partitions is fixed but the dataset size varies.

##### Dynamic partitioning 

When a partition grows to exceed a configured size (on HBase it’s 10GB by default)
it splits into two partitions of equal size.

After a large partition have been split it also can be moved to another node to balance the load. 

An advantage of dynamic partitioning is that the number of partitions adapts to the total data volume.

##### Partitioning proportionally to nodes

With dynamic partitioning, the number of partitions is proportional to the size of the dataset, since the splitting and merging processes keep the size of each partition between some fixed minimum and maximum.

On the other hand, with a fixed number of partitions, the size of each partition is proportional to the size of the dataset.

In both of these cases, the number of partitions is independent from the number of nodes.

A third option, used by Cassandra and Ketama, is to make the number of partitions proportional to the number of nodes—in other words, to have a fixed number of partitons per node.

In this case, the size of each partition grows proportion‐ ally to the dataset size while the number of nodes remains unchanged, but when you increase the number of nodes, the partitions become smaller again.

When a new node joins the cluster, it randomly chooses a fixed number of existing partitions to split, and then takes ownership of one half of each of those split parti‐ tions while leaving the other half of each partition in place.

#### Operations: Automatic or Manual Rebalancing 

For example, Couchbase, Riak generate a suggested partition assignment automatically, 
but require an administrator to approve the change before it is applied.

Fully automated rebalancing can be convienient, because there is the less operational 
work to do for formal maintenance, but it can be unpredictable. 

Rebalancing is an expensive operation, because it requires rerouting requests and moving a large amount of data from one node to another.

If it is not done carefully, this process can overload the network or the nodes and harm the performance of other requests while the rebalancing is in progress.

For that reason, it can be a good thing to have a human in the loop for rebalancing. It’s slower than a fully automatic process, but it can help prevent operational surprises.

#### Request Routing 

**Service discovery** is an instance that answers on question: on which IP address and port number 
I need to connect to if I want to read/write the key X?
> Not limited to databases, but also used in microservices.

On high level are a few approaches:

- allow clients to contact any node (e.g via round-robin balancer)
> accepts request or redirect it to another node an return the result to the client 

- send all requests to a routing tier first, which determines the node that should handle the request 
> acts like a partition-aware load balancer 

- require that clients be aware of the partitioning and the assignment to nodes.
> connect directly to the appropriate node 

In all cases, the key problem is: how does the component making the routing decision (which may be one of the nodes, or the routing tier, or the client) learn about changes in the assignment of partitions to nodes?

Many distributed data systems rely on a separate coordination service such as Zoo‐ Keeper to keep track of this cluster metadata

Each node registers itself in ZooKeeper, and ZooKeeper maintains the authoritative mapping of partitions to nodes. 

Other actors, such as the routing tier or the partitioning-aware client, can subscribe to this information in ZooKeeper. Whenever a partition changes ownership, or a node is added or removed, ZooKeeper notifies the routing tier so that it can keep its routing information up to date.

Cassandra and Riak take a different approach: they use a gossip protocol among the nodes to disseminate any changes in cluster state. Requests can be sent to any node, and that node forwards them to the appropriate node for the requested partition.

When using a routing tier or when sending requests to a random node, clients still need to find the IP addresses to connect to. These are not as fast-changing as the assignment of partitions to nodes, so it is often sufficient to use DNS for this purpose.

#### Summary 

We explored different ways of partitioning a large dataset into smaller subsets. Partitioning is necessary when you have so much data that storing and pro‐ cessing it on a single machine is no longer feasible.
The goal of partitioning is to spread the data and query load evenly across multiple machines, avoiding hot spots.



