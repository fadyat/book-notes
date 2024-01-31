## Affinity

To destribute the pods across different availability zones, we can 
use the `nodeAffinity` property. 

> Availability zones are a way to group resources in a data center.
>
> They are usually used to group resources that are close to each other
> and have a low latency between them.
> 
> Distributing pods across availability zones is a way to make sure that
> if one availability zone goes down, the application will still be
> available in the other availability zones.

Affinity section in the pod definition:

```yaml
...
spec:
  affinity:
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchExpressions:
              - key: app
                operator: In
                values:
                  - my-app
          topologyKey: kubernetes.io/hostname
...
```

The `podAntiAffinity` property is used to make sure that pods with the 
same label are not scheduled on the same node. 

The `requiredDuringSchedulingIgnoredDuringExecution` property is used to
make sure that the pods are not scheduled on the same node.

### Preffered vs Required 

Preffered means that the scheduler will try to schedule the pods on 
different nodes, but if it is not possible, it will schedule them on the 
same node. 

Required means if it can't find a suitable node, it will not schedule the 
pod at all. 

### Node Affinity 

Node affinity is similar to pod affinity, but it is used to schedule pods 
based on the node labels. 

```yaml
...
spec:
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
              - key: disktype
                operator: In
                values:
                  - ssd
...
```

### Inter-pod Affinity 

Inter-pod affinity is used to schedule pods based on the labels of other 
pods. 

### Weighted Affinity 

If scheduler found multiple nodes that match the affinity rules, it will 
use the `weight` property to decide which node to use. 

Nodes with higher weight will be preferred. 

### Resources 

- https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/


