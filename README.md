# Devoxx Docker

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
