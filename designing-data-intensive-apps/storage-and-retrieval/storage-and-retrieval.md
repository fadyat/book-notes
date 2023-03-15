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

The most widely used indexing structure in databases is the **B-Tree**.

Like SSTable, B-Trees keep key-value pairs sorted by key.

B-Trees break the database down into fixed-size __blocks__ or __pages__ (
usually 4K or 8K) and read and write one page at a time.

Each page can be identified using an address or location, which allows one
page to refer to another -- simular as a pointer, but on disk instead of
memory. 

Starts with a **root page**. The page contains several keys and references
to a child pages. Each child is responsible for a continious range of keys.

**Leaf pages** are the pages that contain the actual key-value pairs.

**Branching factor** is the number of children that each page can have.

If you want to change the value of existing key, you have to find the leaf
page that contains the key and update the value there.

If you want to insert a new key, you have to find the leaf page that contains
the key and insert the new key there. If the leaf page is full, you have to
split it into two pages and insert the new key into one of them.

Depth of the tree is O(log N) and tree remains balanced.

#### Making B-Trees reliable

Overwrite doesn't change the location of the page.

Some operations require multiple pages to be overwritten. 
Split of the page need change in the parent page and so on.
This is a dangerous operation, because if the database crashes after only some
pages were written you end up with a **corrupted index**.

In order to make the database resilent to crashes, it's common for B-tree impl
to have a **write-head-log** (WAL). This is an append-only file to which every 
modification must be written before it is applied to the pages itself.
If the database crashes, the WAL can be replayed to restore the database to
a consistent state.

Latches (lightweight locks) are used to prevent multiple threads from 
modifying the same page at the same time.

#### B-Tree optimizaitons

- Instead of overwriting pages and maintaining a WAL, some DB like LMDB
use a copy-on-write scheme. Modified pages are copied to a new location
instead of being overwritten.
- We save space in pages not storing the entire key, but
abbreviating it.
- Sequantial order of leaf pages on disk to make it faster.
- B-Tree variants as a **fractal-trees** borrow sime log-strucutured ideas
to reduce disk seeks.

### Comparing B-Trees and LSM-Trees

LSM-Trees are typically faster for writes, whereas B-Trees are typically faster
for reads. 

LSM-Trees are slow for reads, because the have to check several data structures 
and SSTables at different stages of compaction.

#### Advantages of LSM-Trees

A B-tree index must write every piece of data at least twice.
Log-structured index also rewrite data multiple times.

**Write amplification** is an effect when one writing operation to the database
affect miltiple writes to the disk.

In write-heavy apps, the performance bottleneck is often the disk.

Moreover, LSM-Trees are typically able to sustain higher write rates than
B-Trees, because they have lower write amplification and compaction of SSTables
is faster than B-Tree rebalancing.

LSM-Trees can be compressed better and produce smaller files on disk, because
B-Tree save some memory for fragmentation.

#### Downsides of LSM-Trees

Compaction process sometimes can interfere with reads and writes and slow them
down.

With compaction arises a high write throughput: the more data you write, the
more compaction you have to do.

If write throughput is not configured carefilly, it can happen that compaction
can't keep up with the rate of incoming requests.

An advantage of B-Trees is that each key exists in exacly one place in the index.
This aspect make B-Trees more attractive in databases that want to offer strong 
transaction semantics. 

B-Trees are very ingrained in the architecture of databases.
In new datastores log-structured indexes are becoming more popular.

### Other Indexing Structures

**Primary index** uniquely identifies one row or one document or one vertex.
Other records can refer to it.

Secondary index is not unique, and they often crucial for performing joins
efficiently. 

#### Storing values within the index

**Heap file** is a place were rows are stored in no particular order.
The heap file is common because it avoids duplicating data when multiple 
secondary indexes are present: each index just references a location in a 
heap file.

**Clustered index** is an index that stores the actual data in the same 
place as the index itself.

> For example in MySQL, the primary key is clustered by default.

Compromessive between clustered index(stroring all row data with the index) and 
nonclustered index (storing only refrerences to the data with the index) is 
known as **covering index** or **index with included columns**, which 
stores __some__ of the row data with the index. This allows to run some queries
only using the index.

As with any kind of duplication of data, clustered and covering index can speed
up reads, but the require additional storage and can be overhead on writes.

#### Multi-column Indexes

**Concatenated index** combines several fields into one key by appending one 
column to another.

> This index is useless if you want to search by one of the columns.

**Multi-dementional index** is a more general way of querying several
columns at once.

#### Full-text search and fuzzy indexes

All previous indexes don't allow you to search for __simular__ keys, such as 
misspelled words. Such __fuzzy__ querying requires different techniques.

Full-text engines commanly allow a search for one word to be expand to include 
synonyms and ignore grammatical variations.

Levenshtein automation and other techs that goes with ML document classification.

#### Keeping everything in memory

Advantages of disks:
- they are durable 
- lower cost per GByte

If RAM become cheaper and your dataset are simply not big - **in-memory database**.

- Big performance improvements, no disk overhead needed.
- Can avoid the overheads of encoding in-memory data structures in a form 
that can be written to disk.
- Providing data models that a difficult to implement with disk-based indexes.
(Redis offers a database-like interface for priorty queue and sets).

**Anti-caching** approach works by evicting the least recent used data from 
memory to disk when there is not enough memory, and loading back to memory 
when it's accessed again in the future.

### Transaction Processing or Analytics

**Transaction** - group of reads and writes that form a logical unit.

**OLTP (online transaction processing)** is type of data processing that consists
of executing a number of transactions occurring concurrently.

**OLAP (online analytic processing)** is software for performing analysis at
high speed on large volumes of data.

| Property | OLTP | OLAP |
| ---      | ---  | ---  |
| Main read pattern | Small numbers of records per query | Aggregate over large numbers of records |
| Main write pattern | Random-access, low-latency writes from user-input | Bulk import or event stream |
| Primary used for | End user/customer via web app | Internal analyst, for decision support |
| What data represents | Latest state of data | History of events that happend over time |
| Dataset size | GBytes - TBytes | TBytes - PBytes |

At first time databases where used for both processses. 

There was a trend for companies to stop using their OLTP systems for analytics 
purposes and run it on separate database -- **data warehouse**.

#### Data Warehousing

A **data warehouse** is a separate database that analysts can query to their 
hearts content, without affecting OLTP operations. (contains a read-only copy 
of data from OLTP)

OLTP data extracted by cron or using continious stream of updates, transformed
into an analysis-friendly schema, cleaned up, and then loaded to the DWH.
This process of getting data data into DWH is known as **Extract-Transform-Load (ETL)**.

OLAP can be optimized directly for analytics. For example, she don't need indexes.

**The divergence between OLTP databases and data warehouses**

- The data model of DWH is most commonly relational, SQL is good for it.
- There are many graphical data analysis tools that generate SQL queries, 
visualize the results and allow to explore the data (__drill-down__, __slicing__, 
__dicing__).

Both a focused on supporting either transaction processing or analytics 
workflow, but not both.

**Stars and Snowflakes: Schemas for Analytics**

**Star schema (dimentional modeling)** because when the table relationship are 
visualized, the fact table is in the middle and others are connected to her.

**Fact table** - represents an event that occured at a particular time.
**Dimension table** - table which are referenced via foreign key. (represents the 
__who__, __what__, __when__, __how__ and __why__ of the event)

A variation of this template is known as the **snowflake** schema, where dimensions
are further broken down into subdimesions.

In typical DWH tables are often very wide.

### Column-Oriented Storage

