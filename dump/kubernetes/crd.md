## Custom Resource

Custom resources are extensions of the Kubernetes API.

They are not available in a default Kubernetes installation
and have to be installed/created separately.

All requests are coming to the API server as usual.

### Custom Resource Definition (CRD)

CRD is a way to define a custom resource.

To work with CRD, you need to create a Custom Controller to handle
CRD requests, because the API server doesn't know how your CRD
should behave in different situations.

### Controller vs Operator

Controller works on vanilla Kubernetes objects.

Operator is a controller that works with custom resources.

### Resources

- https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/
- https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/

