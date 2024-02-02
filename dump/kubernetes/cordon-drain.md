## Cordon 

Cordon is a kubernetes feature that allows you to mark a node
as unschedulable.

This means that no new pods will be scheduled on the node.
This is useful when you want to drain a node for maintenance.

Node will have the `SchedulingDisabled` condition set to `True`.

```bash
kubectl cordon <node-name>
```

## Drain 

Drain is a kubernetes feature that allows you to evict all the
pods from a node. This is useful when you want to perform
maintenance on a node.

When you drain a node, kubernetes will evict all the pods
from the node and reschedule them on other nodes.

```bash
kubectl drain <node-name>
```


