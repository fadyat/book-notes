## Argo CD 

### Problems with CI/CD 

Jenkins/Gitlab CI flow:

Build image -> Push image to registry -> Deploy to Kubernetes (using kubectl or helm)

Challenges:

- install and setup tools, like kubectl, helm, etc.
- access to Kubernetes cluster
- configure access to cloud platform, registry
- security challenge (credentials to only particular part)
- no visibility of deployment status 

CD part can be improved by using Argo CD.

### Argo CD 

Argo CD is a declarative, GitOps continuous delivery tool for Kubernetes.

- Part of the cluster 
- ArgoCD pulls changes from Git repository and deploys them to Kubernetes cluster
- Configuration is stored in Git repository (YAML files)
- Separate repositories for environment, application, etc.
- GitLab CI will update files in **separate** git repository
- Separations of concerns, different roles can manage different repositories

Commit image to repository <- ArgoCD detects changes -> Deploy to Kubernetes

Supports:
- k8s manifests
- Helm charts
- Kustomize

### Benefits

GitOps:

- Common interface via Git and ArgoCD (all changes are visible)
- All manual changes are overwritten by ArgoCD (model current + desired state)
- Easy rollback (revert commit in Git)
- Cluster disaster recovery (can create new cluster and point to Git repository) 

k8s:

- Access control with Git (without creating k8s users)
- No external access for automation tools (Jenkins, Gitlab CI)
- ArgoCD as k8s extension (etcd for storing data, k8s operator for managing resources using CRD)
- One instance of ArgoCD can manage multiple clusters (also can one-to-one for environments via git branches or overlays in kustomize)

Git (desired) -> ArgoCD (updating state) -> Kubernetes (current)

https://argo-cd.readthedocs.io/en/stable/

### Resources 

- [TechWorld with Nana - ArgoCD Tutorial](https://youtu.be/MeU5_k9ssrs?si=ho7Dammwhjd9Ba1h)
