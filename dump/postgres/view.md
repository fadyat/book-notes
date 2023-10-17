## `View` vs `Materialized View`

| View                                 | Materialized View                                                        |
|--------------------------------------|--------------------------------------------------------------------------|
| A view is a virtual table.           | A materialized view is a physical copy of the table, stored on the disk. |
| "Updated" every time it is accessed. | "Updated" only via triggers or refresh command (manually).               |
| Slow                                 | Fast, because it is precomputed.                                         |

### Advantages of `View`

- No storage space is required for a view.
- Restrict users from accessing some columns.
- Reducing complexity of queries.

View = alias for user query.

### Advantages of `Materialized View`

- Faster than views, because it is precomputed.
- Can be indexed.

Materialized view = cached query result.

### Resources

- https://www.timescale.com/blog/how-postgresql-views-and-materialized-views-work-and-how-they-influenced-timescaledb-continuous-aggregates/