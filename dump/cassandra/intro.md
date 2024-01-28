### Cassandra

Cassandra is NoSQL database that is designed to handle large amounts of data across many commodity servers, providing
high availability with no single point of failure.

Cassandra is a distributed database, which means that it runs on multiple machines, and these machines communicate
with each other to store and retrieve data.

- No single point of failure
- No master-slave architecture, peer-to-peer architecture (all nodes are the same)
- No strong ACID guarantees, can be configured to provide eventual consistency
- Highly configurable to meet your needs, AP can be tuned to CP
- No joins, no foreign keys, no subqueries - denormalization is the way to go
- Linear scalability
- Wide column store - type of the database in which names and format of the columns can vary from row to row in the same
  table
- Query first approach - data model is designed based on the queries that will be performed

### Resources

- https://medium.com/geekculture/system-design-solutions-when-to-use-cassandra-and-when-not-to-496ba51ef07a
- https://medium.com/@aymannaitcherif/beginners-guide-to-learn-cassandra-part-1-cassandra-overview-bf1634e4ce30