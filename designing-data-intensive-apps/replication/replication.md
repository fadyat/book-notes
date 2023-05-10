# Part 2. Distributed Data

## Chapter 5. Replication

**Replication** means keeping a copy of the same data on multiple machines that are connected via a network.

Reasons for replication:

- Keep data geographically close to your users (reduce latency)
- Allow the system to continue working even if some of its parts have failed (fault tolerance)
- Scale out the number of machines that can serve read queries (increase read throughput)

All difficulties in replication lies in handeling changes to replicated data.

Three popular algos for replicating changes between nodes: [ single-leader, multi-leader, leaderless ]

### Leaders and Followers

**Replica** - node that stores a copy of the database. 

Every write to the database needs to be processed be every replica. 

The most common solution is called **leader-based replication** (or master-slave replication).

It works as follows:

- One replica designed as a leader (or master)

- The other are known as followers (or slaves)
> Whenever leader writes new data to its local storage it also sends the data change 
> to all of his followers as part of **replication log** or **change stream**.

> Each follower takes the log from the leader and updates its local copy of the database accordingly.

- When a client wants to read data it can query any replica (leader or follower)
> Writes are only accepted on the leader, but reads can be handled by any replica.

### Syncronous Versus Asyncronous Replication 

The leader sends the message and wait / don't wait the response from follower.

- If leader suddenly fails we can be sure that the data is still available
on the follower.
- Write to follower is blocking operation. (Need to wait until it will be completed; networks issues)

Good practise: have only one syncronous replica and all others asyncronous;

If syncronous replica becomes unavailable or slow, one of the async is made sync.
(Guarauntees that we have at least two available up-to-date copies on nodes, also 
called **semi-syncronous**)

### Setting Up New Followers 

Process how to set up new follower:
- Take a consistent snapshot of the leader's database at some point in time - 
if possible w/o lock on the entire database. 
- Copy the snapshot to the new follower node.
- The follower connects to the leader an requests all the changes since the snapshot was taken.
(Requires that the snapshot is associated with an exact position in the leader 
replication log).
- When the follower has processed all the data in the replication log since the snapshot,
it has caught up with the leader and is ready.

### Handling Node Outages

How to achieve high availability with leader-based replication?

#### Follower failure: Catch-up recovery

- each follower keeps a log of the data changes it has received from the leader 

#### Leader failure: Failover

Such process is trickier, called **failover**.

One of the followers needs to be prometed to be the new leader, and
all others need to start consuming data changes from new leader. 

Can happen manually or automatically.

Things that can go wrong:

- If async replication is used, new leader may not received all the writes 
from the old leader. What happed if the old leader comes back online?
> The most common solution is for old leader to discard its own writes since it 
> became a follower, which may violate the client durability.

- Discarding writes is especially dangerous if the other systems outside of the database 
need to be coordinated with the db content. (e.g cache based on db index)

- **Split brain** - two nodes both thinking they are the leader and both accepting writes.
> Mechanism for detecting and down one of the nodes.

- Which timeout to use for detecting leader failure?
> If timeout is too short, we may get false positives and failover when the leader is
> still available. If timeout is too long, failover may take a long time.

### Implementation of Replication Logs

#### Statement based replication

In the simplest case, the leader logs every write request that executes and sends that log to all of its followers.
(update, insert, delete).

Break down cases:

- nondeterministic functions (e.g `random()`, `now()`) - different results on different nodes 
> leader can replace the function call with the returned value, but it is not always possible.

- if have autoincrement, or update of existing data - need to be replicated in the same order.
- side effect statements (trigger, stored procedures, user-defined functions).

#### Write-ahead log (WAL) shipping





