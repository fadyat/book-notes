## Container Storage Interface (CSI)

Container Storage Interface (CSI) is an initiative to unify the storage
interface of Container Orchestrator (CO) systems like Kubernetes, Mesos,
Docker, Cloud Foundry, etc. combined with a storage vendors like Ceph,
Portworx, NetApp, etc.

This means, implementing a single CSI for a storage vendor is guaranteed
to work with all COs.

### Internals

Parts of a CSI Driver

- Controller
- Node
- Identity

Communication between these parts is done via gRPC.

CSI Sanity - a tool to test CSI drivers, that your driver implements
correctly all the required methods.

#### CSI Life Cycle

Startup(Identity) -> Create (Controller) -> Use (Node) -> Stop (Node) -> Cleanup (Controller)

- Startup = what's your driver name, what you can do; are you still alive?
- Create = create a volume; attach the volume to the node;
- Use = partition, format, mount the volume; bind mount to the requested container path;
- Stop = unbind; pvc is updated and available for other pods;
- Cleanup = detach the volume from the node; delete the volume;

All methods are idempotent.

### Resources

- https://kubernetes.io/blog/2019/01/15/container-storage-interface-ga/
- https://medium.com/velotio-perspectives/kubernetes-csi-in-action-explained-with-features-and-use-cases-4f966b910774
- https://www.youtube.com/watch?v=FK8Ti39oXEg
- https://www.youtube.com/watch?v=AnfAd6goq-o
