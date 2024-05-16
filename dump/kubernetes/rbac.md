## RBAC 

Role based access control (RBAC) is a method of regulating access to computer or network resources
based on the roles of individual users within an organization.

### API Objects

| Object | Description |
| --- | --- |
| Role | A role contains rules that represent a set of permissions. Permissions are purely additive (there are no "deny" rules). |
| RoleBinding | A role binding grants the permissions defined in a role to a user or set of users. It holds a list of subjects (users, groups, or service accounts), and a reference to the role being granted. |
| ClusterRole | Like a role, but cluster-wide, not namespaced. |
| ClusterRoleBinding | Like a role binding, but cluster-wide. |

### Role 

Role is a namespaced resource, which represents a set of permissions to perform actions on a group of resources in a namespace. 

```yaml
# following pod-reader role allows to get, watch and list pods in the default namespace

apiVersion: rbac.authorization.k8s.io/v1
kind: Role 
metadata:
  namespace: default
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```

### RoleBinding

RoleBinding is a namespaced resource, which grants the permissions defined in a role to a user or set of users. 

```yaml
# following read-pods role binding grants the permissions
# defined in the pod-reader role to the user "aboba"

apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding 
metadata:
  name: read-pods
  namespace: default
subjects:
- kind: User
  name: aboba
  apiGroup: rbac.authorization.k8s.io 
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

### ClusterRole and ClusterRoleBinding

Works the same way, but cluster-wide, not namespaced.

Allows to grant cluster-wide permissions, like managing nodes.

### Usecases

- Granting different permissions to different users
- Granting permissions to a service account 
- Granting permissions to a group of users 
- Granting permissions to a user in a specific namespace 
- Granting permissions to a user in a specific cluster 


