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
