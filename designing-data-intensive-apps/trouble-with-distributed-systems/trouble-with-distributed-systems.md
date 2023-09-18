# Part 2. Distributed Data

## Chapter 8. The Trouble with Distributed Systems

Anything that can go wrong will go wrong.

### Faults and Partial Failures

Single computer: if an internal error occurs, we prefer to crush completely
rather that running in an inconsistent state.

Distributed systems: we prefer to keep running in a partially working state
rather than crashing completely.

Nondeterminism and possibility of partial failures is what makes distributed
systems hard to work with.

### Cloud Computations and Supercomputing

Philosophies on how to build large-scale computing systems:

- High-performance computing (HPC)
  > Supercomputers with thousands of CPUs used for computationally intensive
  > scientific tasks

- Cloud computing
  > Multi-tenant datacenters, commodity computers connected with an IP network, elastic/on-demand
  > resource allocation, and metered billing.

- Traditional enterprise datacenters

Depending on the selected philosophy, handling of the faults will be different.

In a supercomputers, a job typically checkpoints the state of its computation to
durable storage from time to time. Supercomputer is more like single-node computer.
> Node fails -> stopping entire cluster -> repair the node -> restart the job from the last checkpoint.

If we want to make distributed systems work, we must accept the possibility of partial failure and build fault-tolerance
mechanisms into the software.
> We need to build a reliable system from unreliable components.

### Unreliable Networks

**Shared-nothing** systems - a bunch of machines connected by a network,
is the only way to communicate, each has its memory and disk and
can't access the memory or disk of another.

Is not the only way of building systems, but it has become the dominant
approach for building internet services, for several reasons:

- cheap, because it doesn't require specialized hardware
- can use commoditized cloud computing services
- can achieve high reliability by replicating data across multiple DC

The internet and most internal networks in datacenters are **asynchronous
packet networks**.
> One node can send a message to another, but it has no guarantee of when
> or whether the message will be delivered.

Usual way to of handling network issues is **timeout**: after some time you
give up waiting and assume that the response isn't going to arrive.
> We still don't know whether the remote node received the message and is
> processing it, or whether the remote node is down and the message was lost.

### Network Faults in Practice

TLDR: network faults are common and unavoidable.

### Detecting Faults

Examples:

- load balancer needs to stop sending requests to a node that is dead
- a database needs to elect a new leader if the current leader fails

Unfortunately, the uncertainty about the network makes it difficult to tell whether a node is working or not.

In some specific circumstances you might get some feedback to explicitly tell you that something is not working:

- can reach the machine o which the node is running, but no process is listening on the expected port,
  the OS will close or refuse the TCP connection be sending RST or FIN packet in reply.
- node process crashed but the node's OS is still running, a script can notify other nodes about the crash, and
  they can take appropriate action.

Conversely, if something has gone wrong, you may get an error response at some level of the stack, but in general you
have to assume that you will get no response at all.

You can retry a few times, wait for a timeout to elapse, and eventually declare the node dead.

### Timeouts and Unbounded Delays

How long timeout should be?

Long timeout -> slow response time; Short timeout -> false positives

When a node is declared dead, its responsibilities need to be transferred to other nodes, which places additional load
on other nodes and the network.
> If the system is already struggling with a high load, declaring nodes dead prematurely can make the problem worse.

Asynchronous networks have **unbounded delays** (no upper limit on how long it can take for a message to be delivered).

#### Network congestion and queueing

