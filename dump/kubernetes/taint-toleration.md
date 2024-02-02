## Taint 

Taint is a key-value pair that is applied to a node.

It's used to mark a node as unsuitable for a pod to run on.

Taints are used to repel pods from nodes.

## Toleration

Toleration is a key-value pair that is applied to a pod.

It's used to allow a pod to run on a node with a specific taint.

Tolerations are used to attract pods to nodes.

## Taint and Toleration 

Taint and toleration are used together to control the placement of pods on nodes.

For example, if a node has a taint, a pod can be scheduled on the node only if the pod has a matching toleration.

## Effect

Taints have an effect that can be one of the following:

- NoSchedule: The pod will not be scheduled on the node.
- PreferNoSchedule: The pod will be scheduled on the node if there is no other option.
- NoExecute: The pod will be evicted from the node if it's already running on the node.

## Difference between affinity and toleration 

Affinity is used to restrict the placement of pods based on labels,
while toleration is used to attract pods to nodes with specific taints.


