# Part I. Foundations

## Chapter 2. Data Models and Query Languages

Most apps are built by layering one data model on top of another.

For each layer, the key question is: how is it represented in terms of next-lower layer?

Each layer hides his complexity from the next layer by providing a clean interface.

### Relational Model vs Document Model

Best-known data model is relational model, which is based on tables and SQL.

Data organized into relations (tables), where each relation is a set of tuples (rows).

Original scope of relational databases is _business data processing_ and _transaction processing_.

In the 1970's, the network and the hierarchical models were popular, but the relational model won out.

Later, the relational model began to be used in a wide variety of use cases.

#### The Birth of NoSQL

NoSQL is the latest attempt to overthrow the relational model.

Driving forces behind the adoption of NoSQL databases:

- Greater scalability, including large datasets and very high write throughput
- Open source software
- Specialized query languages
- Desire for a more dynamic and expressive data model

**Polyglot persistence** -- using multiple data models within the same application.

#### The Object-Relational Mismatch

**Impedance mismatch** -- the differences between the data model of an application and the data model of the database.

**ORM** (Object-Relational Mapping) libraries where created to solve this problem.

The json representation has better _locality_ than the multi-table relational representation.
If you want to fetch a data in relational model, you need to fetch all data from all tables, then apply joins or some
another hard operations. In NoSQL all relevant data is stored in one place.

#### Many-to-One and Many-to-Many Relationships

Advantages of having special data structures for many-to-one and many-to-many relationships:

- Consistent style and spelling
- Avoiding of ambiguity (like having cities with the same name)
- Ease for updating (for example city changed its name)
- Localization support
- Better search

When you store a text directly, you are duplicating the human-meaningful information in the database.

**Normalization** -- a process of organizing data to minimize redundancy.

#### Are Document Databases Repeating History?

How best to represent such relationships in a database?

NoSQL databases are similar to the hierarchical model, which was popular in the 1970's.

They both work well with one-to-many relationships, but not with many-to-many relationships.

Various solutions were proposed to solve the limitations of the hierarchical model. The relational model, the network
model were created.

##### The network model

The network model is a generalization of the hierarchical model, where each node can have multiple parents.

**Access path** -- a sequence of nodes that you need to traverse to get from one node to another.

It's too hard and long to find a path in a network model.

##### The relational model

No labyrinth based structures, no complicated access paths, just relations which a collection of tuples.

The query optimizer automatically figures out the best way to execute a query.

The relational model made it much easier to add new features to the database.

##### Comparison to document databases

Document databases reverted back to the hierarchical model in one aspect: storing nested records within their parent
record rather than in a separate table.

However, when it comes to representing many-to-one and many-to-many relation‐ ships, relational and document databases
are not fundamentally different: in both cases, the related item is referenced by a unique identifier, which is called a
foreign key in the relational model and a document reference in the document model.

That identifier is resolved at read time by using a join or follow-up queries.

#### Relational vs Document Databases Today

In this chapter, we will concentrate only on the differences in the data model.

The main arguments in favor of the document data model are schema flexibility, better performance due to locality, and
that for some applications it is closer to the data structures used by the application. The relational model counters by
providing better support for joins, and many-to-one and many-to-many relationships.

##### Which data model leads to simpler application code?

If the data in your application has a document-like structure (i.e., a tree of one-to- many relationships, where
typically the entire tree is loaded at once), then it’s
probably a good idea to use a document model.

The relational technique of **shredding** -- splitting a document-like structure into multiple tables can lead to
cumbersome schemas and unnecessarily complicated application code.

The document model has limitations: for example, you cannot refer directly to a nested item within a document.

The poor support for joins in document databases may or may not be a problem, depending on the application. For example,
many-to-many relationships may never be needed in an analytics application that uses a document database to record which
events occurred at which time.

However, if your application does use many-to-many relationships, the document model becomes less appealing.

It’s not possible to say in general which data model leads to simpler application code; it depends on the kinds of
relationships that exist between data items. For highly interconnected data, the document model is awkward, the
relational model is acceptable, and graph models are the most natural.

##### Schema flexibility in the document model

Schema-on-read is similar to dynamic type checking, where schema-on-write is similar to static type checking.

Schema changes are easy in the document model, but they are hard in the relational model.

Schema-on-read approach us advantageous, when items don't have a fixed schema.

##### Data locality for queries

A document is usually stored as a single continuous string, encoded as JSON, BSON, XML, etc.

If your app often needs to access the entire document, there is a performance advantage to this _storage locally_. If
data is split across multiple tables, multiple index lookups are required to retrieve it all, which may require more
disk seeks and take more time.

The locality advantage only applies if you need large parts of the document at the same time.
Generally recommended to keep documents fairly small and avoid writes that increase the size of a document.
These performance limitations significantly reduce the set of situations in which document databases are useful.

It’s worth pointing out that the idea of grouping related data together for locality is not limited to the document
model. For example, Google’s Spanner database offers the same locality properties in a relational data model, by
allowing the schema to declare that a table’s rows should be interleaved (nested) within a parent table

##### Convergence of the document and relational models

It seems that relational and document databases are becoming more similar over time, and that is a good thing: the data
models complement each other. If a database is able to handle document-like data and also perform relational queries on
it, applications can use the combination of features that best fits their needs.

A hybrid of the relational and document models is a good route for databases to take in the future.

### Query Languages for Data

When relational model was introduced, she created a new way for querying data -- SQL, which is a **declarative query**
language. Whereas IMS and CODASYL querying data using **imperative** code.

Many common used programming languages are imperative, like Java, C++, Python, etc.

When SQL was defined, it followed the structure of the relational algebra fairly closely.

An imperative language tells the computer to perform certain operations in a certain order. You can imagine stepping
through the code line by line, evaluating conditions, updating variables, and deciding whether to go around the loop one
more time.

In a declarative query language, like SQL or relational algebra, you just specify the pattern of the data you want —
what conditions the results must meet, and how you want the data to be transformed (e.g., sorted, grouped, and
aggregated) — but not how to achieve that goal. It is up to the database system’s query optimizer to decide which
indexes and which join methods to use, and in which order to execute various parts of the query.

Finally, declarative languages often lend themselves to parallel execution. Today, CPUs are getting faster by adding
more cores, not by running at significantly higher clock speeds than before. Imperative code is very hard to
parallelize across multiple cores and multiple machines, because it specifies instructions that must be performed in
a particular order. Declarative languages have a better chance of getting faster in parallel execution because they
specify only the pattern of the results, not the algorithm that is used to determine the results. The database is free
to use a parallel implementation of the query language, if appropriate.

#### Declarative Queries on the Web

The advantages of declarative query languages are not limited to just databases.

In a web browser, using declarative CSS styling is much better than manipulating styles imperatively in JavaScript.
Similarly, in databases, declarative query languages like SQL turned out to be much better than imperative query APIs.

#### MapReduce Querying

MapReduce is a programming model for processing large amounts of data in bulk across many machines, popularized by
Google.

MapReduce is neither a declarative query language nor a fully imperative query API, but somewhere in between: the logic
of the query is expressed with snippets of code, which are called repeatedly by the processing framework.

It is based on the `map` and `reduce` functions that exist in many functional programming languages.

They must be a _pure_ function, which means that they cannot have side effects, and they must be deterministic, which
means that they must always return the same result when given the same input.

A usability problem with MapReduce is that you have to write two carefully coordinated JavaScript functions, which is
often harder than writing a single query. Moreover, a declarative query language offers more opportunities for a query
optimizer to improve the performance of a query. For these reasons, MongoDB 2.2 added support for a declarative query
language called the **aggregation pipeline**.

The moral of the story is that a NoSQL system may find itself accidentally reinventing SQL, albeit in disguise.

### Graph-Like Data Models

But what if many-to-many relationships are very common in your data? The relational model can handle simple cases of
many-to-many relationships, but as the connections within your data become more complex, it becomes more natural to
start modeling your data as a graph.

A graph consists of two kinds of objects: **nodes** (entities) and **edges** (relationships).

Typical examples of graph:

- Social graphs: nodes are people, edges are friendships
- Web graphs: nodes are web pages, edges are hyperlinks
- Road networks: nodes are intersections, edges are roads

However, graphs are not limited to such homogeneous data: an equally powerful use of graphs is to provide a consistent
way of storing completely different types of objects in a single datastore.

#### Property Graphs

Each vertex consists of:

- Unique identifier
- Set of outgoing edges
- Set of incoming edges
- Set of properties (key-value pairs)

Each edge consists of:

- Unique identifier
- Tail (source) vertex
- Head (destination) vertex
- Label, which is a string that describes the type of relationship
- Set of properties (key-value pairs)

You can think of a graph store as consisting of two relational tables, one for vertices and one for edges.

Some important aspects:

- Any vertex cah have an edge connecting it with any other vertex.
- Given any vertex, you can efficiently find both its incoming and its outgoing edges and thus **traverse** the graph.
- By using different labels for different kinds of relationships, you can store several kinds of info in a single graph,
  while still maintaining a clean data model.

Those features give graphs a great deal of flexibility for data modeling.

Graphs are good for evolvability: as you add features to your application, a graph can easily be extended to accommodate
changes in your application’s data structures.

#### The Cypher Query Language

**Cypher** is a declarative query language for property graphs, created for the _Neo4j_ graph database.

Syntax:

- `()` for nodes
- `[]` for edges
- `:` for labels
- `-->` for directed edges
- `-[]->` for directed edges with labels
- `--` for undirected edges
- `>` can be replaced by `<` to indicate the direction of the edge

Example:

```cypher
MATCH (person:Person)-[:ACTED_IN]->(movie:Movie)
WHERE person.name = "Tom Hanks"
RETURN movie.title
```

#### Graph queries in SQL

SQL is not a good fit for graph data, but it is possible to use SQL to query graphs.

In a relational database, you usually know in advance which joins you need in your query. In a graph query, you may need
to traverse a variable number of edges before you find the vertex you’re looking for that is, the number of joins is not
fixed in advance.

The idea of variable-length traversal paths in a query can be expressed using something called a **recursive common
table expression** (the `WITH RECURSIVE` clause in SQL).

The same query will have much more lines of code in SQL than in Cypher.

#### Triple-Stores and SPARQL

The triple-store model is mostly equivalent to the property graph model, using different words to describe the same
ideas. It is nevertheless worth discussing, because there are various tools and languages for triple-stores that can be
valuable additions to your toolbox for building applications.

In a triple-store all info is stored in the form of triples, which are three-element tuples of the form `(subject,
predicate, object)`. For example, the triple `(Tom Hanks, acted_in, Forrest Gump)` means that Tom Hanks acted in the
movie Forrest Gump.

The object is one of two things:

- A value in a primitive data type, such as a string or a number. For example, the triple `(Tom Hanks, age, 56)` means
  that Tom Hanks is 56 years old.
- Another vertex in a graph. For example, the triple `(Tom Hanks, acted_in, Forrest Gump)` means that Tom Hanks acted in
  the movie Forrest Gump.

##### The semantic web

The semantic web is fundamentally a simple and reasonable idea: websites already publish information as text and
pictures for humans to read, so why don’t they also publish information as machine-readable data for computers to read?

The Resource Description Framework (RDF) was intended as a mechanism for different websites to publish data in
consistent formats, so that computers could read data from different websites and combine it into a single coherent view
of the world.

Was overhyped in the 2000s.

##### The SPARQL query language

SPARQL is a query language for triple-stores using the RDF data model. It predates Cypher, and since Cypher’s pattern
matching is borrowed from SPARQL, they look quite similar.

SPARQL is a nice query language—even if the semantic web never happens, it can be a powerful tool for applications to
use internally.

#### The Foundation: Datalog

Datalog is much older language than Cypher or SPARQL. It provides the foundation that later query languages build upon.

Datalog’s data model is similar to the triple-store model, generalized a bit. Instead of writing a triple
`(subject, predicate, object)`, you write a fact `predicate(subject, object)`. For example, the fact `acted_in(Tom
Hanks, Forrest Gump)` means that Tom Hanks acted in the movie Forrest Gump.

The Datalog approach requires a different kind of thinking to the other query languages discussed in this chapter, but
it’s a very powerful approach, because rules can be combined and reused in different queries. It’s less convenient for
simple one-off queries, but it can cope better if your data is complex.
