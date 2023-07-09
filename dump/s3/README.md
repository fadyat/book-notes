### S3

Simple Storage Service (S3) is a storage service offered by AWS.

Which is a cloud-based object storage service designed to store and retrieve any amount of data from anywhere.

It's highly durable, available and scalable.

AWS stores data in S3 as objects, which consist of data and metadata that describes the data, can range
from a few bytes to terabytes size.

Each object is identified by a unique key, and the data is stored in buckets, which are essentially containers for
storing objects.

Each bucket has a unique DNS name, making it easy to access and retrieve data from anywhere on the web.

Replication is used to ensure high availability and durability of data.

Fragmentation is used to increase the speed of data retrieval, as well as to reduce the cost of storing data.

### Core concepts

- Bucket: container for storing objects
- Object: any data stored in S3
- Multi-part upload: upload large objects in parts
- ACL: access control list, defines who can access the object and what permissions they have
- Lifecycle: automate deleting and archiving objects based on some rules

### S3 storage classes

Amazon S3 offers several storage classes that allow you to choose the right level of availability, throughput, and cost
based on your data requirements. Here are some of the main S3 storage classes:

- S3 Standard: for general purpose storage of frequently accessed data
- S3 Intelligent-Tiering: for data with unknown or changing access patterns
- S3 Standard-Infrequent Access: for long-lived, but less frequently accessed data
- S3 One Zone-Infrequent Access: for long-lived, but less frequently accessed data that doesn't require multiple
  Availability Zone data resilience (by default all data replicated across 3 AZs)
- S3 Glacier: for long-term archive and digital preservation
- S3 Glacier Deep Archive: for long-term archive and digital preservation at the lowest cost, for data that can be
  retrieved within 12 hours

### Sources

- [AWS S3](https://aws.amazon.com/s3/)
- [AWS S3 storage classes](https://aws.amazon.com/s3/storage-classes/)
- [Selectel S3, RU](https://selectel.ru/blog/object-storage-s3/)