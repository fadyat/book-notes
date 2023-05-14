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

Every write is appended to the log on disk before the write is applied to the database.
> The log is an append-only sequence of bytes, so it is easy to replicate.

- Main disadvantage: log describes data on very low level, WAL contains details of which bytes were changed in which disk block.
> If the database is upgraded her storage format, it's typically 
> isn't possible to run different versions on the leader and the follower.

#### Logical (row-based) log replication 

An alternative is to use diff log formats for replication and for the storage 
engine, which allows the replication log to be decoupled from the storage engine internals.

- The replication log is a sequence of logical records, each of which describes a
write operation that was executed on a row in a table.

- For an inserted row, new values of all columns
- For a deleted row, some identification of the row (e.g primary key)
- For an updated row, identification of the row and new values of columns 

A transaction that affects multiple rows generate several log records followed by a record 
indicating commit. 

- Easier to be kept a backward compatible, different versions on leader and follower.
- Easier to parse for external systems, (e.g dwh, cache), called **change data capture** (CDC).

#### Trigger-based replication 

In previous cases all replication is done by the db itself and what if
we want only some of the data to be replicated?

**Trigger** - a piece of code that is executed when a data change occurs in a db.

It have an opportunity yo log this change into separate table, from which can 
be read by external process. (e.g Oracle GoldenGate)

### Problems with Replication Lag

**Evantual consistency** - the guarantee that if no new updates are made
to the data, eventually all reads will return the same result.

**Replication lag** - the amount of time by which a replica is behind the leader.

3 problems that are likely to occur when there is replication lag.

#### Reading Your Own Writes

When user write data record goes through the leader, but when he reads it 
may goes through the follower.

We need **read-after-write consistency** - guarantee that if the user reloads the page,
he will see any themselves submitted updates. 

Possible techniques:

- When reading smth that the user may modify - read from leader. 
- Using of others criteria to deside read from leader or not. 
(e.g last update time for 1 minute, replication lag value)
- Remember the timestamp of the most recent write on client. 
- If replicas distributes across multiple dc, any request that needs to be served 
by the leader must be routed to that dc.

**Cross-device read-after-write consistency** - guarantee that if the user writes data on one device,
the data will be immediately visible when accessing the service from another device.

Some additional issues:

- Metadata with last update of the client needs to be centralized (shared between 
all devices of the user).
- No guarantee of routing to the same dc.

#### Monotonic Reads 

Possible to see **moving backward in time** when reading from the async follower.
This may happen if user makes reads from different replicas. 

**Monotonic reads** - guarantee that this kind of anomaly doesn't happen.
> One way of achieving this is to always read from the same replica for 
a particular user (e.g by hash of user_id).

#### Consistent Prefix Reads 

**Consistent prefix reads** - guarantee that if a sequence of writes happens 
in the certain order, then anyone reading those writes will see them appear 
in the same order.
> One solution is to make sure that any writes that are causally related 
to each other are written to the same partition - but in some apps that can't 
be done efficiently. 

#### Solutions for Replication Lag 

Transactions are too expensive in distributed world?

### Multi-Leader Replication

Leader based replication model allows accept writes only on one node. 

#### Use Cases for Multi-Leader Replication 

##### Multi-datacenter operation

You can have a leader on each datacenter.
> Within dc regular leader-follower replication is used. 
> Between dc, each dc leader replicates its changes to the others leaders.

Advantages:

- Performance 
> Each write happens on the local dc, instead of dc with single leader.

- Tolerance of datacenter outages 
> In failure cases each dc can continue to operate independently, instead of
making follower from another dc as a leader. 

- Tolerance of network problems 
> Traffic between dc usually goes over the public internet, which is less 
reliable than the local network.

External tools for multi-leader replication:
- BDR (Bi-Directional Replication) for PostgreSQL
- Tungsten Replicator for MySQL
- GoldenGate for Oracle

Multi-leader replication must be done with care. (e.g conflict resolution,
autoincrementing keys on each dc)

##### Clients with offine operation 

Each device have a local database that acts as a leader and then 
asyncronous share data with followers. In that case each device can 
continue to operate independently when offline. 

Each device is a "datacenter".

##### Collaborative editing

Real-time collaborative editing of a document by multiple users.

When one user edit a document, the changes are apply to his local replica 
and async replicated to the server and other users who are edititing the same document. 

App must obtain a lock on the document before a user can edit it, if you want to 
guarantee that there will be no conflicts.

Perfect case: very small units of change and non-locking.

#### Handling Write Conflicts

This is the biggest problem with multi-leader replication, conflict resolution 
is required. 

##### Synchronous Versus Asynchronous Conflict Detection

In single leader database the second write will block and wait until the first 
write is committed or aborted.

It's late to detect conflict when both writes are successfully committed.

Conflict detection is syncronous - wait for the write to be replicated to all replicas 
before telling the user that the write was successful. By doing this 
we may lose the main advantage of multi-leader replication - independent writes. 

##### Conflict Avoidance

Just avoid the conflict by design. :)

All writes to a particular record are sent to the same leader.
> Failure, relocation cases are possible.

##### Converging toward a consistent state 

There is no defined ordering of writes, so it's not clear what the final value should be.
> When they will share data between leaders, which value is correct?

Replicas must resolve conflict in **convergent** way - all replicas must arrive 
to a final value when all changes have been replicated.

Ways of achieving convergent conflict resolution:

- Each write have a unique ID (e.g UUID) and the replica with the highest ID wins.
(if timestamps are used it's called **last write wins** - LWW)
> Dataloss is possible
- Each replica have a unique ID and let writes that originated at a highest numbered 
replica win.
> Dataloss is possible 
- Merge the values together 
- Record the conflict in explicit data structure the preserves all info, 
write application code that resolver the conflict later.

##### Custom conflict resolution logic 

On write/read conflict resolution logic can be implemented in application code.

Conflict resolution usually applies to a document/row, not for entire transaction.

#### Multi-Leader Replication Topologies

Describes the communication paths alogns which writes are propagated from 
one node to another.

- Circular topology 
  ```
  [ ] -> [ ]
   ^      |
   |      v
  [ ] <- [ ]
  ```

- Star topology 
  ```
  [ ]    [ ]
   \      / 
    \    / 
     [  ]
    /    \
   /      \
  [ ]    [ ]
  ```

- All-to-all topology
  ```
  [ ] -- [ ]
   | \__/ |
   | /  \ |
  [ ] -- [ ]
  ```

> To prevent an infinite loop each node have a unique ID, in replication log 
also. 

The problem in circular and star topology it's if the node is down,
the replication can't continue. (must have custom resolution logic)

All-to-all topology also have problems - some network links may be faster 
than the others.
> As a result some messages overtake anothers. (e.g update of the value, that 
not exists on current node yet, but exists on another leader node)

This is problem is simular to consistent prefix reads. 
> We need sure that all nodes process inserts first. 
> Simply attaching a timestamp to each message is not enough.
> To order this events correctly a technique called **vector clocks** is used.

### Leaderless Replication

Dynamo-style databases use leaderless replication.
(e.g Cassandra, Riak, Voldemort)

#### Writing to the Database When a Node Is Down 

Client do writes to all replicas -> One of the replicas is down ->
Client ignores the error and continue to write to the other replicas -> 
Replica that was down comes back online -> Data on that replica is now stale (outdated).

> To solve this problem client make a read from multiple replicas (as in a write case) 
and compare the versions of values, if they different. 

##### Read repair and anti-entropy

How does catch up the stale data?

Two popular mechanisms used in Dynamo-style databases:

- Read repair 
> When a client reads from replica it can detect any stale responses.

- Anti-entropy process 
> Periodically compare data on all replicas and copies any missing values from one replica to another.
Order of writes is not important, may be a significant delay before data is copied. 

Not all systems implement both mechanisms.

##### Quorums for reading and writing 

As long as `w + r > n` we expect to get an up-to-date value when reading.
`w` - number of nodes to write to,
`r` - number of nodes to read from,
`n` - total number of replicas

> If `w + r > n` then at least one node that is both in the write set and the read set.

#### Limitations of Quorum Consistency 

`w` and `r` it's not necessary need to be majority of nodes.
> They require just single overlap.

You may also set `w + r <= n`.
> In this case, reads and writes are still be sent to `n` nodes, but a 
smaller number of successful responses is required to consider the operation succeed. 

With a smaller values you more likely to read stale values, on the other 
side it allows lower latency and higher availability.

However, even with `w + r > n` you can still have edge cases with stale values:

- If a **sloppy quorum** is used, the `w` writes may end up on different nodes that `r` reads. 

- If two writes occur concurrently - conflicts, and the conflict resolution.

- If a write happend concurrently with a read, some replicas may have the new value. 

- If some writes failed, successful writes are not rolled back. (quorum value will be lower)

- Even if everything is working correctly, there are edge cases in which 
you can get unlucky timing. 

Returning of latest write it's not simple. 

##### Monitoring staleness

In leader-based replication is possible to monitor a replication lag, 
because all writes go through the leader and in the same order. 

#### Sloppy Quorums and Hinted Handoff


