# K8S infra
- Node, API, k8s master, CLI, Storage

## Pod
- Pod is a minimal object in k8s, containers + network namespace container
- `kcl apply`
- `kcl create` 
- `kcl edit` - edit in cluster [bad practise]
- `kcl describe` - info
- `kcl delete`

## ReplicaSet
- ReplicaSet for pods manipulation, can pass pods number
- `selector` and `labels` for pods detecting, template for pods configuration like labels, pods specification
- `kcl scale`
- When something is crushed resplicaset will scale app back

> Cascade deletion. When you delete parent (for example replicaset) all his child will be delete too (pods).
> Can be disabled with flag

## Deployment
- Deployment for updating apps inplace
- `strategy` field: how app will be updated? [recreate, rollingUpdate(good choice) - one by one]
- Creates new replica set
- `kcl rollout undo` - roll back

## Probes
- Liveness probe: lifecycle control
- Readiness probe: is app ready for getting traffic?
- Startup probe: is app initialized and started successfully?
- `readiness probe`, `liveness probe` in `spec` part
- `kcl explain` doc about specific field

## Resources 
- Limits, number of resources that pod can use (upper bound)
- Requests, number of resources that reserved for pod on node 
- When using more memory will be killed
- When using more cpu will be throttling

### QoS classes
- which resources are setted for app
- [ Guaranteed, Burstable, BestEffort ] left > right

## Config map
- section `data`, `volume` on pod level, `volumeMounts` on container level
> `kcl exec` enter volume, like `docker exec`.
- `kcl port-forward`

## Secrets
- [ generic, docker-registry, tls ]
- `kcl create secret`
- `env` section in containers scecification

## Service
- `clusterIP`, traffic balancing ...

## Ingress
- Routing rules for services, services to pods
- For working ingress needed ingress controller (need to be downloaded manually)

## PV/PVC
- Data storage for pod
- Persistent volume (where, size, id)
- Persistent volume claim

# Cluster components

- etcd
- api server
- controllerr-manager
- scheduler
- kubelet
- kube-proxy

> containers, network, dns

## Etcd

- Key-value storage, storing all cluster info
- `etcdctl` - util for managing etcd cluster
- Needs fast disks

## API server

- central component of k8s
- only he can talk to etcd
- rest api
- authentication and authorization

## Controller-manager

- control loop that watches the shared state of the cluster through the apiserver and makes changes attempting to
move the current state towards the desired state. 
- node controller
- replicaset controller
- endpoints controller
- garbage collector

## Scheduler

- binding pods to specific node by QoS, (anti)-affinity, requested resources, priority class
> priority classes good for staging and production separation in one cluster

## Kubelet

- working agent on each node
- make commands to docker daemon
- creates pods

## Kube-proxy

- on each server
- managing network rules on nodes
- actually implements Service (ipvs, iptables)

Service:
- static IP
- DNS name in kube-dns
- not proxy
- NAT problems in Linux

When creating new service:
- he get ip
- dns record

iptables:
- random destination of traffic
- chains

ipvs:
- support some balancing techs


# Network

- main goal: connect nodes and podes
- network interfaces
- network plugins (flannel, calico): network for nodes, pods, network policies

## Ingress


