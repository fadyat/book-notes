## Feature expectations

Time: 5 min

- Usecases 
- Scenarios that will not be covered 
- Who will use 
- How many will use 
- Usage patterns 

## Estimations 

Time: 5 min

- Throughput (QPS for read and write queries)
- Latency expected from the system (for read and write queries)
- Read/Write ratio 
- Traffic estimates (write, read - qps, volume of data)
- Storage estimates 
- Memory estimates

> If we are using a cache, what is the kind of data we want to store in cache
>
> How much RAM and how many machines do we need for us to achieve this?
>
> Amount of data you want to store in disk/ssd

## Design goals

Time: 5 min 

- Latency and Throughput requirements
- Consistency vs Availability (weak/strong/eventual, failover/replication)

## High Level Design 

Time: 5-10 min

- APIs for read/write scenarios for crucial components
- database schema
- basic algorithm
- high level design for read/write scenario 

## Deep Dive

Time: 15-20 min 

- Scaling the algorithm
- Scaling individual components 
> Availability, Consistency and Scale story for each component 
>
> Consistency and Availability pattern 

- Think about the following components, how they would fit in and how it would help 

> - DNS
>
> - CDN (push/pull)
>
> - Load balancers (active-passive, active-active, l4, l7)
>
> - Reverse proxy
>
> - Application level scaling (microservices, service delivery)
>
> - DB (relational: sharding, denormalization, replication, sql tuning; nosql: kv, document, graph, column)
>
> - Caches (client, cdn, webserver, db, app)
>
> - Asynchronism (mq, task queue, back pressure)
>
> - Communication (tcp, udp, rest, rpc)

## Justify 

Time: 5 min 

- Throughput of each layer
- Latency caused between each layer
- Overall latency justification

### Resources

- https://leetcode.com/discuss/career/229177/my-system-design-template

