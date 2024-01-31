## Pod Disruption Budget (PDB) 

A PDB limits the number of pods of a replicated application that are 
down simultaneously from voluntary disruptions.

> A voluntary disruption is a disruption that is caused by the user.
>
> For example, a user might want to drain a node to perform
> maintenance on it or to remove it from the cluster.

A PDB can be used to ensure that a certain number of pods are available
at all times.

### PDB is not configured 

For example, we have a deployment with 2 replicas and a PDB is not configured.

```txt
-------------------------
| node1 | node2 | node3 |
-------------------------
| pod1  | pod2  |       |
-------------------------
```

We start a drain on node1. The scheduler will move pod1 to node3.

```txt
-------------------------
| node1 | node2 | node3 |
-------------------------
|       | pod2  | pod1  |
-------------------------
```

Let's say that pod2 is in pending state.

And now we start a drain on node2. The scheduler will move pod2 to node3.

And we have no pods running.

### PDB is configured

```yaml
apiVersion: policy/v1beta1
kind: PodDisruptionBudget
metadata:
  name: your-microservice
  namespace: your-namespace
spec:
  minAvailable: 50% # 50% is just to keep the example simple, it is not the recommended value
  selector:
# ...
```

In now case pod on node2 won't be moved to node3, because the PDB is configured.

### Best practices 

- Use % instead of absolute number of pods 
- Don't use non-resolvable combinations of replicas and maxUnavailable
> For example, you have 1 replica and maxUnavailable=20% - it will never work.

### References

- https://github.com/mercari/production-readiness-checklist/blob/master/docs/concepts/pod-disruption-budget.md 
- https://kubernetes.io/docs/tasks/run-application/configure-pdb/


