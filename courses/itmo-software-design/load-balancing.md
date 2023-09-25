## Load balancing

A load balancer is a software or hardware device that keeps any one server from 
getting overloaded with requests by distributing the requests across multiple
servers.

There are two primary approaches to load balancing:

- **Dynamic**

    > uses algorithms that take into account the current state of each server and 

- **Static**

    > distributes traffic w/o making these adjustments, some equal amount of traffic to each server group


### Dynamic load balancing algorithms

- Least connection:

    > Checks which servers have the fewest connections open at the time and sends traffic to those servers.
    >
    > This assumes all connections require roughly equal processing power.

- Weighted least connection:

    > Gives administrators the ability to assign different weights to each server, assuming that some servers can handle more connections than others.

- Weighted response time 

    > Averages the response time of each server, and combines that with the number of connections each server has open to determine where to send traffic.
    >
    > By sending traffic to the servers with the quickest response time, the algorithm ensures faster service for users.

- Resource-based

    > Distributes load based on what resources each server has available at the time. 
    >
    > Specialized software (called an "agent") running on each server measures that server's available CPU and memory, and the load balancer queries the agent before distributing traffic to that server.

### Static load balancing algorithms

- Round robin

    > Sends traffic to each server in turn, regardless of the current load on each server.

- Weighted round robin

    > Gives administrators the ability to assign different weights to each server, assuming that some servers can handle more connections than others.

- IP hash

    > fn(ip) -> hash -> server 
    >
    > Based on the hash, the connection is assigned to a specific server.

## Resources 

- https://samwho.dev/load-balancing/
- https://www.cloudflare.com/learning/performance/types-of-load-balancing-algorithms/


