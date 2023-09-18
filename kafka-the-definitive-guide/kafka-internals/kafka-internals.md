# Chapter 5

## Kafka Internals 

### Cluster Membership

Kafka uses Apache Zookeeper to maintain the list of brokers that are currently 
members of a cluster. 

Every broker has a unique identifier that is either set in the broker 
configuration file or automatically generated. 

Every time a broker process starts, it registers itself with ID in Zookeeper by 
creating an **ephemeral node**. Different Kafka components subscribe to the 
`/broker/ids` path in Zookeeper where brokers are registered so they get notified 
when brokers are added or removed.

If you try to start another broker with the same ID, you will get an error.

When a broker loses connectivity to Zookeeper, the ephemeral node that the broker 
created when starting will be automatically removed from Zookeeper. Kafka components 
that are watching the list of brokers will be notified that the broker is gone. 

### The Controller 

The controller is one of the brokers that in addition to usual broker functionality, 
is responsible for electing partition leader.

The first broker becomes the controller by creating and ephemeral node in Zoo.
When others brokers start, they also try to create this node, but receive a exception.
They create a Zookeeper watch on the controller node - get notified of changes.

To summarize, Kafka uses Zookeeper’s ephemeral node feature to elect a controller 
and to notify the controller when nodes join and leave the cluster. 

The controller is responsible for electing leaders among the partitions and 
replicas whenever it notices nodes join and leave the cluster. 
The controller uses the epoch number to prevent a "split brain" scenario where two nodes believe each is the current controller.

### Replication 

 Each topic is partitioned, and each partition can have multiple replicas. 

 Those replicas are stored on brokers, and each broker typically stores thousands 
 of replicas belonging to different topics and partitions.

 Types of replicas:

 - leader
   > each partition has a single replica designed as the leader.
   >
   > all produce and consume requests go through the leader, in order to guarantee consistency. 

 - follower
   > their only job is to replicate messages from the leader and stay up-to-date with the most recent messages the leader has.
   >
   > when leader crushed, one of the followers become a new leader.

Another task the leader is responsible for is knowing which of the follower replicas is up-to-date with the leader.

