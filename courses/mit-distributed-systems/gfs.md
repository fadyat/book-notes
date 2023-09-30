## GFS (The Google File System)

Touches several topics, related to distributed systems:

- parallel computation
- fault tolerance
- replication
- consistency

GFS is a good example of successful real-world system, that have documented design and implementation.

### Context

- Many Google services needed a fast unified storage system

> MapReduce, Crawling, Indexing, Log storage/analysis

- Shared across many applications

> 1000s of clients, 100s of chunkservers, 1 coordinator;

- Huge capacity
- High performance
- Fault tolerance
- Aimed at batch processing, not interactive use

### Capacity story

- big files split into 64MB chunks
- each chunk sharded into 3 replicas
- each chunk is a Linux file

### Throughput story

- many clients, each reading/writing a few large files
- huge parallelism

### Fault tolerance story

- each file chunk replicated on 3 chunkservers
- client writes sends to all of chunk's copies
- read just needs to consult one copy

**MIT notes**: https://pdos.csail.mit.edu/6.824/notes/l-gfs.txt with more details.