# Part I. Foundations of Data Systems

## Chapter 3. Storage and Retrieval

In fundamental level database doing two things:

- When you give some data, she store this data
- When you ask some data, she give this data

As a developer you need to know how your selected storage work.

### Data Structures That Power Your database

In order to efficiently find the value for a particular key in the database,
we need a different data structure: an **index**.

**Index** is an __additional__ structure that is derived from the primary data.
Many databases allow you to add and remove indexes and doesn't affect the content
of the database. Any kind of index usually slows down writes, because the index
also needs to be updated every time data is written.

#### Hash Indexes

Key-value storages are quite simular to the dictionary type, which is 
implemented as a hash map.

We a have in-memory hash map where every key is mapped to a bye offset in the 
data-file -- the location at which the value can be found.

When the record is updated we don't overwrite previous, just append new.
Compaction -- the process which return only the recient value for each key.
Also we can merge couple segments and continue works with merged.

Reasons why appending is good:
- Merging and appending a sequantial operations, which is much faster than
random access memory.
- Concurrency and crash recovery are much simplier
- Merging old segments avoid the problem of data files getting fragmented
over time.

Limitations:
- The hash table mush be fit in memory.
- Range queries are not efficient. You can scan over all keys between [l, r]
you have to look up each key individually. 

#### SSTables and LSM-Trees

Making a simple change from hash indexes: the key-value pairs are sorted by key.
We call this format **Sorted String Table** or **SSTable** for short.

Advantages over log segments with hash indexes:
- Merging segments is simple and efficient, even if the files are bigger than
the available memory. We reading each input files side by side.

> For each key we also remember the period of time when it was created.

- Search of particular element works faster.
- Opportunity to group records into a block and compress before writing
to disk.

**Constructing and maintaining SSTables**

How to make your data sorted? 

Maintaining on disk is possible, but maintaining in memory is much easier (
AVL tree, red-black tree)

Storage engine props:
- When a write comes in, add it to an in-memory balanced tree data structure.
(memtable)
- When memtable gets bigger than few MBytes -- write it out to disk as
SSTable file.
- In order to serve a read request, first try to find the key in the memtable,
then in the most recent on-disk segment and etc.
- From time to time, run a merging nad compaction process in the background.

Works very well, have problem - if database crashes, the most recent writes
are lost. In order to avoid that problem, we can keep a separate log on disk
to which every write is immediatly appended. Every time when memtable is written
out to an SSTable, the log can be discarded.

**Making a LSM-Tree out of SSTables**

Such algo is used in LevelDB and RocksDB.

Log-Structed Merge-Tree (LSM Tree).

Storage engines that are based in this principle of merging and compacting
sorted files are often called LSM storage engines.

Lucene, an indexing engine for full-text search used by Elasticsearch, uses a 
simular method for storing its __term dictionary__.
Word is a key, value is a list of IDs if all documents that contain the word.

**Performance optimizations**

LSM work slow, when key don't exist.

Bloom filter is a memory-efficient data structure for approximating the contents
of a set. It can tell you if a key doesn't appear in the database.

There are also different storagies to determine the order and timing
of how SSTables are merged and compacted.

The most common are __size-tiered__ and __leveled__.
- Size-tiered: merge the smallest SSTables first.
- Leveled: the key range split up into smaller SSTables and older data is 
moved into separate levels, which allows the compaction to proceed more
incrementally and use less disk space.

### B-Trees

