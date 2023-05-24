# Part 2. Distributed Data 

## Chapter 7. Transactions

Transaction is the way for an application to group several reads and writes 
together into a logical unit.
> Conceptually, all reads and writes are executed as one operation: either transaction 
succeds (commit) or fails (rollback).

- When fails can be retried
- Error handling is much simpler

Transactions are not the law of nature, they were created with a purpose 
to __simplify the programming model__ for applications to accessing the database.
By using transactions, it's free to ignore some potential problems, because 
database takes care them instead (**safety guarantees**).
> You need to understand where to use transactions and where not to use them.

### The Slippery Concept of a Transaction 

NoSQL hype which aimed to improve scalability and performance,
by abandoning transactions. 

The truth is not that simple: like every other technical design choice, transactions have advantages and limitations.

### The Meaning of ACID

Atomicity, Consistency, Isolation, Durability

Systems which not meet the ACID criteria sometimes called BASE (Basically Available, Soft state, Eventual consistency).

#### Atomicity

In general, atomic refers to something that cannot be broken down into smaller parts.

The word means simular but subtly different things in different branches of computing. 

> In multi-threaded programming, if one thread executes an atomic operation, that 
means there is no way to another thread could see the operation only partially
completed. 

In the context of ACID, atomicity is not about concurrency. It doesn't describe 
what happens if several processes try to access the same data at the same time - 
is covered by **isolation**.

Atomicity describes what happens if a client wants to make several writes, but 
a fault occurs after some of the writes have been processed.

If writes are grouped together into an atomic transaction and she can't be completed,
then the transaction is aborted and the database must discard or undo any writes it has made so far in that transaction.

Without atomicity, if an error occurs partway through making multiple changes, it’s difficult to know which changes have taken effect and which haven’t.

Perhaps **abortability** would have been a better term than atomicity, but we will stick with atomicity since that’s the usual word.

#### Consistency 

This word is terribly overloaded.

- Replica consistency and the issue of eventual consistency that arises in 
async replicated systems.
- Consist hashing is an approach to partitioning that some systems use for rebalancing.
- In CAP theorem, consistency means linearizability.
- In ACID, consistency refers to an application-specific notion of the database being in a “good state”.

If you have certain statements about your data (invariants) that must always be true.

> The database should always change from one valid state to another valid state.

Consistency is a property of the application, when AID is a property of the database.

The letter C doesn't really belong in ACID.

#### Isolation 

Most databases are accessed by several clients at the same time. That is no problem if they are reading and writing different parts of the database, but if they are accessing the same database records, you can run into concurrency problems (race conditions).

Isolation in the sense of ACID means that concurrently executing transactions are isolated from each other: they cannot step on each other’s toes.

Isolation often formalized as a **serilizability**, which means that each transaction 
can pretend that it is the only transaction running on the entire database.

The database ensures that when the transactions have committed, the result is the same as if they had run serially (one after another), even though in reality they may have run concurrently.

#### Durability 

The purpose of a database system is to provide a safe place where data can be stored without fear of losing it. 

Durability is the promise that once a transaction has com‐ mitted successfully, any data it has written will not be forgotten, even if there is a hardware fault or the database crashes.

> In a single-node database, durability typically means that the data has been written to nonvolatile storage such as a hard drive or SSD. 
>
> In a replicated database, durabil‐ ity may mean that the data has been successfully copied to some number of nodes. In order to provide a durability guarantee, a database must wait until these writes or replications are complete before reporting a transaction as successfully committed.

Perfect durability doesn't exist: if all your hard disks and all your backups 
are destroyed at the same time.

### Single-Object and Multi-Object Operations

Recap in ACID:

- Atomicity
> all-or-nothing guarantee

- Isolation
> transactions can run concurrently without messing each other up

Multi-object transactions are often needed if several pieces of data need to be kept in sync.

-  Violating isolation: one transaction reads another transaction's uncommitted writes 
(**dirty read**)

Multi-object transactions require some way of determining which read and write operations belong to the same transaction. In relational databases, that is typically done based on the client’s TCP connection to the database server: on any particular connection, everything between a BEGIN TRANSACTION and a COMMIT statement is considered to be part of the same transaction.

On the other hand, many NoSQL databases don't have such a way of grouping operations together. 

#### Single-object writes 

Atomicity and isolation are also apply when single object is being changed.
> e.g writing a 20 Kb json document to a database, that may failure 
halfway through writing.

Atomicity can be implemented by using a log for crash recovery, and isolation 
can be implemented by using a lock on each object (allowing only one transaction 
to access the object at a time).

Some databases also provide more complex atomic operations, such as incrementing
operation, which removes the need for a read-modify-write cycle.

Simularly popular is a compare-and-set operation, which allows a write to 
happen only if the value has not been concurrently changed by someone else.

These single-object operations are useful, as they can prevent lost updates when sev‐ eral clients try to write to the same object concurrently.

However, they are not transactions in the usual sense of the word. Compare-and-set and other single-object operations have been dubbed “light‐ weight transactions”.

A transaction is usually understood as a mechanism for grouping multiple operations on multiple objects into one unit of execution.

#### The need for multi-object transactions 

Not implemented in multiple databases, because it's hard to implement it across 
partitions and in some cases high availability or performance is more important.

There are some use cases in which single-object inserts, updates, and deletes are sufficient:

- Row of one table have a reference to a row in another table.
> Multi-object transactions allow you to ensure that these refer‐ ences remain valid: when inserting several records that refer to one another, the foreign keys have to be correct and up to date, or the data becomes nonsensical.

- Updating a denormalized information (update several objects in one go)

- Maintaining indexes
> Update of column triggers an update of secondary index.

Such applications can still be implemented without transactions. However, error han‐ dling becomes much more complicated without atomicity, and the lack of isolation can cause concurrency problems.

#### Handling errors and aborts

A key feature - abort of a transaction.

Not all systems follow that phlosophy, e.g datastores with leaderless replication work 
much more on a "best effort" basis - "database will do much as it can, and if it 
runs into an error, it won't undo something it has already done". So that's an 
application responsibility to recover from errors.

Errors will inevitably happen, but many software developers prefer to think only about the happy path rather than the intricacies of error handling. 

Retrying an aborted transaction is a simple and effective error handling mechanism, it isn’t perfect:

- If the transaction actually succeeded, but the network failed while the server tried to acknowledge the successful commit to the client (so the client thinks it failed), then retrying the transaction causes it to be performed twice—unless you have an additional application-level deduplication mechanism in place.

- If the error is due to overload, retrying the transaction will make the problem worse, not better. To avoid such feedback cycles, you can limit the number of retries, use exponential backoff, and handle overload-related errors differently from other errors (if possible).

-  It is only worth retrying after transient errors (for example due to deadlock, 
isolation violation); after a permanent error (such as constraint violation) a 
retry will simply fail again.

- If the transaction also has side effects outside of the database, those side effects may happen even if the transaction is aborted.

- If the client process fails while retrying, any data it was trying to write to the database is lost.

### Weak Isolation Levels

If two transactions don't touch the same data, they can safety execute in parallel.

Concurrency is hard and databases a trying to hide concurrency issues from application 
developers by providing __transaction isolation__.

**Serializable isolation** means that the databases guarantees that transactions 
have the same effect as if the run serially (one after another) without any
concurrency.

In practice, isolation is unfortunately not that simple. Serializable isolation has a 
performance cost, and many databases don’t want to pay that price.

#### Read Committed 

- No dirty reads
> When reading you will only see data that has been committed.
- No dirty writes
> When writing you will only overwrite data that has been committed.

##### No dirty reads

Can another transation see data of another transaction that has not been committed yet?
If yes, then it's called **dirty read**.

When to prevent dirty reads?

- If a transaction needs to update several objects
- If a transaction aborts, any writes it has made need to be rolled back 

##### No dirty writes

What happens if two transactions concurrently try to update the same object in a database? We don’t know in which order the writes will happen, but we normally assume that the later write overwrites the earlier write.

By preventing dirty writes, this isolation level avoids some kinds of concurrency problems:

- If transactions update multiple objects, dirty writes can lead to a bad outcome.

- However, read committed does not prevent the race condition between two counter increments

##### Implementing read committed 

Most commonly, databases prevent dirty writes by using row-level locks: when a transaction wants to modify a particular object (row or document), it must first acquire a lock on that object.

How do we prevent dirty reads? One option would be to use the same lock, and to require any transaction that wants to read an object to briefly acquire the lock and then release it again immediately after reading.

However, the approach of requiring read locks does not work well in practice, because one long-running write transaction can force many read-only transactions to wait until the long-running transaction has completed.

Better solution: for every object that is written, the database remembers both the old com‐ mitted value and the new value set by the transaction that currently holds the write lock. 

