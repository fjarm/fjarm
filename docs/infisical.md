# Infisical

This doc discusses how to use Infisical to inject secrets into containers.

## Local setup

When running containers that require secrets locally, use the following commands to log into and initialize Infisical:

```bash
infisical login
infisical init
```

`infisical init` will require selecting a project - the Secrets Manager project is usually the best.

After running the commands, secrets can be injected into `podman compose` like this:

```bash
infisical run --path="/helloworld" -- podman compose --file api/internal/helloworld/v1/deployment/docker-compose.yml up --detach
```

To bring down the compose services, use:

```bash
PORT=8000 podman compose down -file api/internal/helloworld/v1/deployment/docker-compose.yml
```

Running the [podman clean up script](../scripts/podman_cleanup.sh) removes the `helloworld:local` image from the local registry.

Running the following command restores it:

```bash
bazel run //api/internal/helloworld/v1/deployment/ghcr:helloworld_load
```

## Local testing

After bringing up the services locally, testing can be done either using the relevant mobile client or using the CLI:

```bash
curl -X POST --data '{}' --header "Content-Type: application/json" --header "request-id: abc123" -o - http://localhost:8000/fjarm.helloworld.v1.HelloWorldService/GetHelloWorld
```

## References and links

* [Podman compose setup docs](https://podman-desktop.io/docs/compose/setting-up-compose)
* [Infisical CLI config docs](https://infisical.com/docs/cli/project-config)
