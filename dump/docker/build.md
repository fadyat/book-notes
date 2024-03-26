## Buildkit

Buildkit is an improved version of the Docker build system.

Default builder for users after Docker 23.0.

It provides some additional features like:

- improved caching capabilities
- parallelize building distinct layers
- lazy pulling of images (pull only the layers you need)

When using Buildkit, you quickly notice that the output of the docker
build command looks cleaner and more structured.

```shell
DOCKER_BUILDKIT=1 docker build --platform linux/amd64 . -t some-image:some-tag
DOCKER_BUILDKIT=1 docker push some-image:some-tag
```

## Buildx

Buildx is a plugin for Docker that enables you to use the full
potential of Buildkit.

It was created because Buildkit supports many new configuration options,
that cannot all be integrated into the `docker build` command in a backwards
compatible way.

```shell
docker buildx create --bootstrap --name mybuilder
docker buildx use mybuilder
```

Benefits of Buildx is that it allows to store the build cache in a separate
registry, which can be useful when building images on a CI/CD server.

```shell
docker buildx build --platform linux/amd64,linux/arm64 . \
    -t some-image:some-tag \
    --push \
    --cache-from=type=registry,ref=myregistry.example.com/some-image:cache \
    --cache-to=type=registry,ref=myregistry.example.com/some-image:cache
```

Also, Buildx allows to build images for multiple platforms in parallel,
you can specify the platforms you want to build for using the `--platform`
flag.
