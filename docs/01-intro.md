# Introduction

Welcome to this hands-on workshop where you'll learn how to build the basic
features of a container runtime from scratch.  
Through a series of exercises, you'll implement core container features using
Go.

Things that we will cover during this workshop:

- Process isolation and management
- Namespace implementation
- Resource control with cgroups
- Filesystem operations
- Network configuration
- Volume management

This that we _will not_ cover:

- Downloading an image from Docker hub, the code for this is provided for you
- Using overlayfs, if you manage to finish all the exercises before the 3h mark
  this would be a nice next step.

## Prerequisites

- Docker Desktop

## Development Environment

If you're on MacOS or Windows, you can use the provided dev container
environment as the exercises require Linux-specific capabilities. Two options
are available:

1. **VS Code / JetBrains DevContainer**: Configuration provided in
   `.devcontainer/`
2. **Docker Compose**: Run `docker compose run --rm -P --build shell` in the
   `.devcontainer/` directory

## Workshop Structure

The workshop is divided into the following exercises, each building upon the
previous ones:

### 1. Process Management

- [Process Creation Basics](02-process-creation.md)

  - Creating parent and child processes

- [Namespace Isolation](03-namespace-isolation.md)
  - PID namespace implementation
  - UTS namespace for hostname isolation
  - Basic process isolation

### 2. Container Foundation

- [Namespaces and Root Directory](04-namespaces-and-chroot.md)

  - Managing multiple namespaces
  - Implementing chroot
  - Directory structure setup

- [Resource Control with cgroups](05-cgroups.md)
  - CPU limitations
  - Memory constraints
  - Process resource management

### 3. Advanced Features

- [Volume Management](06-volumes.md)

  - Implementing bind mounts
  - Volume persistence
  - Data sharing between host and container

- [Network Configuration](07-network.md)
  - Network namespace setup
  - Virtual ethernet (veth) pairs
  - Basic networking capabilities

## Building and Running

Basic commands to get started:

```bash
# Build the project
make

# Run a basic container
sudo ./bin/devoxx-docker run alpine /bin/sh
```

## Additional Resources

- [Linux Namespaces](https://man7.org/linux/man-pages/man7/namespaces.7.html)
- [Control Groups v2](https://www.kernel.org/doc/Documentation/cgroup-v2.txt)
- [Container Networking](https://docs.docker.com/network/)
- [OCI Runtime Specification](https://github.com/opencontainers/runtime-spec)

## Getting Help

- Use `make help` to see available commands
- Check the documentation in each exercise
- Refer to the hints and command references in each exercise file when stuck

[Start the workshop](02-process-creation.md)
