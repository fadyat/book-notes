# Chapter 3. Pods: running containers in Kubernetes

## 3.1 Introducing pods 

Pod is a co-located group of containers and represents the basic building block in Kubernetes.

### Why do we even need pods?

> Containers are designed to run only a single process per container.
>
> Multiple unrelated processes in single container = managing logs, processes, restarting, etc.

Each process in its own container.

### Understanding pods 

Because you're not supposed to group multiple processes into a single container, 
welcome - pod.

You can take advantage of all the features containers provide, while at the 
same time giving the processes the illusion of running together.

Filesystems on each containers are different, all run under the same network 
and UTS namespaces (they share the same IP address and port space).

All pods in k8s cluster reside in a single flat, shared, network-address
space, each pod can access other by IP address.

**Summary:**

Pods are logical hosts and behave much like physical hosts or VMs in the non-container world.

Processes running in the same pod are like processes running on the same physical or virtual machine, except that each process is encapsulated in a container.

Pod is also the basic unit of scaling.
> Kubernetes can’t horizontally scale individual containers; instead, it scales whole pods.

The main reason to put multiple containers into a single pod is when the application consists of one main process and one or more complementary processes.
> For example, the main container in a pod could be a web server that serves 
> files from a certain directory, while a sidecar container periodically 
> downloads content from an external source and stores it.

To recap how containers should be grouped into pods - when deciding
whether to put two containers into a single pod or into two
separate pods, you always need to ask yourself the following questions:

- Do they need to be run together or can they run on different hosts?
- Do they represent a single whole or are they independent components?
- Must they be scaled together or individually?

## 3.2 Creating pods from YAML/JSON descriptors

Defining YAML manifests -> storing them in vcs -> using kubectl -> calling 
k8s api.

`kubectl get pod <pod-name> -o yaml`

Structere:

- k8s api version
- type of k8s object
- metadata (name, namespace, labels, other info)
- specification (list of containers, volumes, etc.)
- status (info about the running pod, condition, internals)

`kubectl explain pods`

`kubectl logs <pod-name> -c <explicit-container-name>`

Forwarding + connecting: `kubectl port-forward` + `curl`

## 3.3 Organizing pods with labels


