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


