# Chapter 3. Pods: running containers in Kubernetes

## 3.1 Introducing pods

A **pod** is a co-located group of containers represents the basic building block in
Kubernetes.

The key thing about pods is that when a pod does contain multiple containers, all
of them are always run on a single worker node - it never spans multiple worker nodes.

### 3.1.1 Understanding why we need pods

- Why do we even need pods?
- Why can't we run containers directly?
- Why would we even need to run multiple containers together?
- Can't we put all our processes into a single container?

Containers are designed to run only a single process per container (unless the process
itself spawns child processes). If you run multiple processes in a single container,
it's your responsibility to keep all this processes running, manage their logs etc.

### 3.1.2 Understanding pods

Because you're not supposed to group multiple processes into a single container, it's
obvious you need another high-level abstraction that allow you to bind containers together
and manage them as a single unit.

You can take advantage of all the features containers provide, while at the same time giving
the processes the illusion of running together.

You want containers inside each group to share certain resources, although not all, so
that they're not fully isolated. Kubernetes achieves this by configuring Docker to have
all containers of a pod share the same set of Linux namespaces instead of each container
having its own set.

But when it comes to the filesystem, things are a little different. Because most of the
container's filesystem comes from the container image, by default, the filesystem
of each container is fully isolated from other containers.
However, it's possible to have them share file directories using a Kubernetes concept called a **volume**.

One thing to stress here is that because containers in a pod run in the same Network namespace,
they share the same IP address and port space. Containers of different pods can never run into port
conflicts, because each pod has a separate port space. All the containers in a pod also have the same
loop-back network interface, so a container can communicate with other containers in the same pod through localhost.



