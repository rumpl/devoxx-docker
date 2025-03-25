# Devoxx Docker

## Running a dev container
If you are on MacOs or Windows you'll need to develop inside a container to be able to run and test the code as it uses Linux specific capabilities.
To do so we propose you 2 options:

### Devcontainer
We provide a devcontainer configuration that you can use directly from VsCode or JetBrains IDEs.
The `devcontainer.json` file is located in the `.devcontainer` directory.

### Compose
You can also use the `docker-compose.yml` file to run a fully configured container, in the `.devcontainer` directory run the following command:
```console
$ docker compose run --rm -P --build shell
```
it will open a shell inside the container where you will be able to run the all the commands from the [#build](#build) and [#run](#run) sections below.

You will need to build the container the first time you run it, but after that you can remove the `--build` flag from the command above.

## Build

```console
$ make
```

This should create a binary `bin/devoxx-docker`

## Run

```console
$ sudo ./bin/devoxx-docker run alpine /bin/sh
```

This will pull the image if it's not present and run the given command.
