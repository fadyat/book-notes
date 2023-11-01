## Partitioning

RDBMS context: division of a table into distinct independent
tables.

Horizontal partitioning (by row) -- different rows in different
tables.

### Reasons

- Easier to manage
- Performance

## Partitioning in PostgreSQL

- PG 8.1 (2005): inheritance-based
- PG 10 (2017): declarative partitioning

### Declarative partitioning

- partitioning method
- partition key (columns, expressions + value determines data routing)
- partition boundaries (where each table value starts and ends)

```sql
create table cust (id int, signup date)
partion by range (signup);

create table cust_2020
partition of cust for values from
('2020-01-01') to ('2021-01-01');
```

> partitions may be partitioned themselves (sub-partitioning)

### PostgreSQL limits

- database size: unlimited
- tables per database: 1.4 billion
- table size: 32 TB (low for current times)
- default block size: 8192 bytes
- rows per block: depends

### How partitioning can help

- disk size limitations (each partition on different disk)
- performance (partitioning pruning - don't have to scan all rows, only subset + small fast indexes)
- maintenance (deletions, `drop table` + `alter table ... detach partition`)
- [`vacuum`](https://www.postgresql.org/docs/current/sql-vacuum.html)

## Partitioning method

- range
- list
- hash

### Resources

- https://www.youtube.com/watch?v=edQZauVU-ws
- https://www.timescale.com/learn/when-to-consider-postgres-partitioning
- https://supabase.com/blog/postgres-dynamic-table-partitioning