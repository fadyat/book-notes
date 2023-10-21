### Docker Context

Docker contexts - an easy way to store connection information for multiple Docker
engines and switch between them.

> Also, can be used to switch between different Kubernetes clusters, docker
> swarms, Amazon ECS clusters, etc.

Example of environment variables, that will be switched:

```bash
DOCKER_HOST="..."
DOCKER_TLS_VERIFY="..."
DOCKER_CERT_PATH="..."
```

### Resources

- https://www.youtube.com/watch?v=x0Kbj4lEOag
- https://docs.docker.com/engine/context/working-with-contexts/