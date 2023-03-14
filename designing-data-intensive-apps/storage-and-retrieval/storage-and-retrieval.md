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





