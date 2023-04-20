# Part I. Foundations of Data Systems

## Chapter 4. Encoding and Evolution 

In most cases changes to an application's features also requires change to 
data that it stores.

Relational databases have one schema, that can be changed through migrations.
Schemaless databases don't enforces to schema, so the database can contain a mixture
of older and newer data. 

Old and new versions of the code, and old and new data formats may potentially 
all coexist in the system all the time. In order for the system to continue 
running smoothly, we need to maintain compatibility in both directions:

- **Backward compatibility**
> Newer code can read data that was written by older code

- **Forward compatibility**
> Older code can read data that was written by newer code 

### Formats for Encoding Data

Programs usually work with data in two different representations:

- In memory, data is kept in objects, structs, lists, arrays, maps etc. Optimized 
for efficient access and manipulation by the CPU (via pointers)
- When you want to write data to a file or send it over network, you have to encode it 
as some kind of self-contained sequence of bytes (ex: JSON document)

**Encoding (serialization, marshalling)** is the translation from the in-memory representation to a byte of 
sequence, and reverse is called **decoding**

#### Language-Specific Formats 

Many programming languages come with built-in support for encoding in-memory objects 
into byte sequences.

These encoding libraries are very convenient, because they allow to do it with 
minimal additional code. However, they also have a number of deep problems:

- The encoding is often tied to a particular programming language and reading data 
in another language is very difficult. 
- In order to restore data in the same object types, the decoding process needs to
be able to instantiate arbitrary classes. (Security problems)
- Versioning data is often an afterthought.
- Efficiency is also often afterthought.

It's generally a bad idea to use your language built-in encoding for anything other 
than very transient purposes.

#### JSON, XML and Binary Variants

- Types problems. 
- Unicode support, but they don't support binary strings.
- Optional Schemaless

Despite these flaws they are good enough for many purposes. 

**Binary encoding**

JSON by textual encoding and JSON by binary encoding have a small size difference.
It’s not clear whether such a small space reduction is worth the loss of human-readability.

#### Thift and Protocol Buffers

Binary encoding libraries.

```Protocol
message Person {
    required string user_name = 1;
    optional int64 number     = 2;
    repeated string interests = 3;
}
```

- instead of field names the encoded data contains **encoded tags** (numbers 1, 2 in example)
- `requires` and `optional` makes no difference to how the field is encoded; checks in runtime

**Field tags and schema evolution**

- encoded record is just a concatenation of encoded fields. Each field is identified by tag number, 
and annotated with datatype. You can change the name, but cannot tag = invalid.
- old code can read the records written by new code.
- new code can always read new data, because tag is constant. If you add new field 
and make it required the check will failed.
- remove a field is just like adding a field, but reversed 


**Datatypes and schema evolution**

Proto:

- int32 and int64 works fine
- don't have `list`, have `repeated` word
- `required`, `optional`

#### Avro

Another binary encoding format, different from Protocol Buffers.

```Avro
record Person {
    string               userName;
    union { null, long } favoriteNumber = null;
    array<string>        interests; 
}
```

- No tag numbers, no types in schema;  only values

To parse the binary data, you go through the fields in the order that they 
appear in the schema and use the schema to tell you the datatype of each field.
Reading data must use the exact same schema as the code that wrote the data. 

**The writer's schema and the reader's schema**

Encode - writer's schema, decode - reader's schema 

- Schemas don't have to be the same, only compatible. Avro library resolves by looking
at the both schemas side by side.

**Schema evolution rules**

- for backward: add/remove field that has a default value.
- changing the datatype is possible
- changing the name is backward compatible, but not forward compatible

**But what is the writer's schema?**

How does the reader know the writer’s schema with which a particular piece of
data was encoded?

The answer depends on the context in which Avro us being used:
- include writer schema ones in the beginning
- include version number at the beginning of every record and keep list of schemas 
- negotiate the schema version on connection setup

**Dynamically generated schemas**

Avro is friendlier to dynamically generated schemas. (generating new Avro schema solve the problem)

**Code generation and dynamically typed languages**

After a schema has been defined, you can generate code that implements this schema 
in programming language of your choice. 

#### The Merits of Schemas 


Schema evolution allows the same kind of flexibility as a schemaless JSON databases,
while also providing better guarantees about your data and better tooling.

### Modes of Dataflow 

Who encodes the data, and who decodes it?

#### Dataflow Through Databases 

In a database the process that writes to the database encodes the data, and the 
process that reads from the database decodes it. There may just be a single 
process accessing database, in which case the reader is simply a later version of the same process 
== storing something in the database as sending a message to your future self.

Backward compatibility is necessary - can't read old data. 

Forward compatibility can be omitted. 

**Different values written at different times**

**Migrations** is a process that rewrite data into a new schema.

- Expensive
- Values for previous columns filling with `null`

**Archival storage**

Dump using latest schema version.

Also it's a good opportunity to encode the data in an analytics-friendly format 
as a Parquet. 

#### Dataflow Through Services: REST and RPC

- Client, server, exposed server API = service
- Service oriented architecture = microservices architecture 

**Web services**

- Service with underlying HTTP protocol for talking - web service 
- REST is not a protocol, but rather a design philosophy that builds upon the 
principles of HTTP.
- An API designed according to the principles of REST is called RESTful.
- SOAP
- A definition such as OpenAPI (Swagger) can be used to describe RESTful API and 
produce documentation. 

**The problems with remote procedure calls (RPCs)**

The RPC model tries to make a request to a remote network service look the same 
as calling a function in your programming language, within the same process. (the 
abstraction is called **location transparency**)

- Local functions are predictable, RPCs because of network are not.
- Local function call either succeeds or failds, RPCs have another possible 
outcome - timeout. In that case you simply don't know what happened.
- If you retry a failed network request, it could happen that the request are actually getting 
through, but responses are getting lost. You will retry the same request multiple times,
unless you build a **idempotent** into your service.
- In local function you can effectively pass pointers to objects in local memory,
but in RPC you have to pass the data itself.
- The client and the service may be implemented in different programming languages,
so the RPC framework must translate datatypes between the two languages.

All of these factors means that there's no point in trying to make RPCs look like
local function calls.

REST doesn't try to hide the fact that this is a network protocol. 

**Current direction of RPC**

The new generation of RPC frameworks is more explicit about the fact that they are 
network protocols. 

Finagle and Rest.li use futures (promises) to encapsulate the asynchronous actions 
that may fail. There also good for making parallel requests and compare results.

gRPC supports streams, where call consists of a series of requests and responses
over time.

Some of the provide **service discovery** - a way to find the address of a service 
from a client.

RPC protocols with a binary encoding format can achieve better performance than
RESTful APIs with a JSON encoding format.

But REST API is better for experimentation and debugging.

RPC often used for internal communication between services, while REST is used for 
external communication with clients.

** Data encoding and evolution for RPC**

For evolvability is important that RPC clients and servers can be changed 
and deployed independently. 

We only need a backward compatibility on requests (previous clients can talk to
new servers), and a forward compatibility on responses (new server response 
will be readable by old clients).

Service compatibility makes harder when your service is used by other services 
that are not under your control (for example in your company).

For REST API you can just use version numbers in the URL.

#### Message-Passing Dataflow 



