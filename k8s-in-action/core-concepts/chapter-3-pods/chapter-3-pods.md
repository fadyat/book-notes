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

Processes running in the same pod are like processes running on the same physical or virtual machine, except that each
process is encapsulated in a container.

Pod is also the basic unit of scaling.
> Kubernetes can’t horizontally scale individual containers; instead, it scales whole pods.

The main reason to put multiple containers into a single pod is when the application consists of one main process and
one or more complementary processes.
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

Structure:

- k8s api version
- type of k8s object
- metadata (name, namespace, labels, other info)
- specification (list of containers, volumes, etc.)
- status (info about the running pod, condition, internals)

`kubectl explain pods`

`kubectl logs <pod-name> -c <explicit-container-name>`

Forwarding + connecting: `kubectl port-forward` + `curl`

## 3.3 Organizing pods with labels

As the number of pods increases, the need for categorizing them
into subsets becomes more and more evident.

Organizing pods and all other Kubernetes objects is done through labels.

A label is an arbitrary key-value pair you attach to a resource, which is then utilized when selecting resources using
label selectors
(resources are filtered based on whether they include the label specified in the selector)
A resource can have more than one label, as long as the keys of those labels are unique within that resource.

## 3.4 Listing subsets of pods through label selectors

> Canary release - deploying a new version of an application to a small subset of users to test it before rolling it out
> to the entire user base.

```shell
# allows to show the labels of selected resources
kubectk get pods --show-labels

# allows to filter resources which have a specific label
kubectl get pods -L [<label-key>, ...]
```

Label selector allows you to select resources based on some criteria:

- contains/doesn't contain a label with a specific key
- contains/doesn't contain a label with a specific key and equal/not equal to a specific value

```shell
# selecting resources based on label selectors
kubectl get pods -l <label-key>=<label-value>

# selecting resources which don't have a specific label
kubectl get pods -l '!<label-key>'
```

## 3.5 Using labels and selectors to constrain pod scheduling

Pods aren't the only resources that can be labeled; nodes can be labeled as well.

Basically, pods are scheduled on every node that meets the criteria (CPU, memory, etc.),
but you can also use labels to constrain pod scheduling.

To do this, you need to create a label on a node and then use a node selector in the pod specification.

```shell
# label a node
kubectl label nodes <node-name> <label-key>=<label-value>
```

Node selector tells to the scheduler to deploy the pod only on nodes that have the specified label.

```yaml
spec:
  nodeSelector:
    <label-key>: <label-value>
```

The importance and usefulness of label selectors will become more evident when we talk about
Replication Controllers and Services.

## 3.6 Annotating pods

Annotations are also key-value pairs, so in essence, they’re similar to labels, but they aren’t meant to hold
identifying information.

They can’t be used to group objects the way labels can. While objects can be selected through label selectors, there’s
no such thing as an annotation selector.

On the other hand, annotations can hold much larger pieces of information and are primarily meant to be used by tools.
Certain annotations are automatically added to objects by Kubernetes, but others are added by users manually.

Annotations are also commonly used when introducing new features to Kubernetes. Usually, alpha and beta versions of new
features don’t introduce any new fields to API objects. Annotations are used instead of fields, and then once the
required API changes have become clear and been agreed upon by the Kubernetes developers, new fields are introduced and
the related annotations deprecated.

A great use of annotations is adding descriptions for each pod or other API object, so that everyone using the cluster
can quickly look up information about each individual object. For example, an annotation used to specify the name of
the person who created the object can make collaboration between everyone working on the cluster much easier.

Example of default k8s annotations is `kubernetes.io/created-by` - contains information about the creator of the object.

```shell
# adding an annotation to a pod
kubectl annotate pod <pod-name> <annotation-key>=<annotation-value>
```

It's better to use unique prefixes for annotations to avoid conflicts with other tools.

## 3.7 Using namespaces to group resources

Namespaces provide the scope for object names, so the names of resources only need to be unique within the namespace.

By default, Kubernetes has three namespaces: [`default`, `kube-system`, `kube-public`].

Namespaces are good for isolating resources, but they can also be used to divide resources between different teams or
projects.

```shell
# listing all the namespaces
kubectl get namespaces

# creating a new namespace
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Namespace
metadata:
  name: <namespace-name>
EOF

# creating a pod in a specific namespace
kubectl apply -f <pod-manifest> -n <namespace-name>

# listing all the pods in a specific namespace
kubectl get pods -n <namespace-name>
```

Although namespaces allow you to isolate objects into distinct groups,
which allows you to operate only on those belonging to the specified
namespace, they don’t provide any kind of isolation of running objects.

For example, pods in different namespaces can communicate with each other
if they have the necessary network access (depending on the network
solution you’re using).

## 3.8 Stopping and removing pods

When deleting a pod, Kubernetes sends a `SIGTERM` signal to the main process in the container,
which allows the process to shut down gracefully (by default, the process has 30 seconds to shut down).

If the process doesn't shut down within that time, Kubernetes sends a `SIGKILL` signal, which forces the process to shut
down immediately.

```shell
# simple way to delete a pod
kubectl delete pod <pod-name>

# deleting a pod by label selector
kubectl delete pods -l <label-key>=<label-value>

# deleting pods (not only pods) in a specific namespace
kubectl delete namespace <namespace-name>

# deleting pods with keeping the namespace
# be sure, that pods aren't controlled by a replication controller
kubectl delete pods --all -n <namespace-name>
```

## 3.9 Summary

- Pods are the basic building blocks in Kubernetes.
- Pods are groups of containers that share the same network and UTS namespaces.
- Pods are the basic unit of scaling in Kubernetes.
- Pods are organized using labels and selectors.
- Pods can be annotated with additional information.
- Pods can be grouped into namespaces to provide scope for object names.


