## Load balancing

A load balancer is a software or hardware device that keeps any one server from
getting overloaded with requests by distributing the requests across multiple
servers.

Used to increase the capacity and reliability of applications. Session management
can improve the performance of applications by reducing the load
on individual servers.

> Also, can provide increased scalability, security and improved user experience.

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

  > Gives administrators the ability to assign different weights to each server, assuming that some servers can handle
  more connections than others.

- Weighted response time

  > Averages the response time of each server, and combines that with the number of connections each server has open to
  determine where to send traffic.
  >
  > By sending traffic to the servers with the quickest response time, the algorithm ensures faster service for users.

- Resource-based

  > Distributes load based on what resources each server has available at the time.
  >
  > Specialized software (called an "agent") running on each server measures that server's available CPU and memory, and
  the load balancer queries the agent before distributing traffic to that server.

### Static load balancing algorithms

- Round robin

  > Sends traffic to each server in turn, regardless of the current load on each server.

- Weighted round robin

  > Gives administrators the ability to assign different weights to each server, assuming that some servers can handle
  more connections than others.

- IP hash

  > fn(ip) -> hash -> server
  >
  > Based on the hash, the connection is assigned to a specific server.

### Layers

Typically, load balancers are implemented in the following layers:

- **L4** - transport

> Transport layer, distributes requests based on the network variables like IP address and destination port.
>
> Performing network addressing translations w/o inspecting the content of the packets.
>
> Used when you don't need to make decisions based on the content of the packets.
>
> E.g: video streaming, voice calls, DNS, etc.

- **L7** - application

> Application layer, distributes requests based on data found in application layer protocols such as HTTP.
>
> Can further distribute requests based on URLs, cookies, headers, etc.
>
> Used when you need to make decisions based on the content of the packets.
>
> Reacher set of algorithms, but more expensive (slower than L4).

- **GSLB** - global server load balancing

> Distributes requests across multiple data centers, which may be in different geographic locations.
>
> GSLB can be used to route requests to the closest data center to the user, or to balance the load across multiple
> data centers.

## Resources

- https://samwho.dev/load-balancing/
- https://www.cloudflare.com/learning/performance/types-of-load-balancing-algorithms/
- https://www.appviewx.com/education-center/load-balancer-and-types/
- https://kemptechnologies.com/blog/layer-4-vs.-layer-7-load-balancing-what's-the-difference-and-which-one-do-i-need


