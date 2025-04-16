# Introduction

Welcome to this hands-on workshop where you'll learn how to build the basic
features of a container runtime from scratch. Through a series of exercises,
you'll implement core container features using Go.


⚠️ **WARNING: This is a workshop project only!**
This container runtime is designed for educational purposes to demonstrate container concepts.   
It is NOT suitable for production use or any real-world deployments.   
For production environments, please use established container runtimes like Docker, containerd, or CRI-O.


## First Step - Clone the Repository

Before we begin, you need to clone this repository:

```console
git clone https://github.com/rumpl/devoxx-docker.git
cd devoxx-docker
```

Go ahead and run these commands now.

Things that we will cover during this workshop:

- Process isolation and management
- Namespace implementation
- Resource control with cgroups
- Filesystem operations
- Network configuration
- Volume management

Things that we _will not_ cover:

- Downloading an image from Docker hub, the code for this is provided for you
- Using overlayfs, if you manage to finish all the exercises before the 3h mark
  this would be a nice next step.

If all goes well, at the end of this workshop you will be able to run an alpine
(or any other) container downloaded from Docker hub.

Here is a sneak peak of how this would look like:

```console
vscode ➜ /workspaces/devoxx-docker (main) $ sudo ./bin/devoxx-docker run alpine /bin/sh
Running /bin/sh in alpine
/ # ping google.com
PING google.com (216.58.214.174): 56 data bytes
64 bytes from 216.58.214.174: seq=0 ttl=62 time=10.580 ms
^C
--- google.com ping statistics ---
2 packets transmitted, 1 packets received, 50% packet loss
round-trip min/avg/max = 10.580/10.580/10.580 ms
/ #
```

# Prerequisites

If you are on Windows or Mac, all you need is Docker Desktop and an IDE that
knows how to run a devcontainer.

If you are on Linux please use a VM, we will be calling things that require root
privileges and could potentially damage your system.

# Development Environment

If you're on MacOS or Windows, you can use the provided dev container
environment as the exercises require Linux-specific capabilities. We recommend
using VS Code or JetBrains DevContainer.

## Important: Open in Dev Container

After cloning the repository:

1. Open the project in VS Code
2. When prompted, click "Reopen in Container" or use the command palette (F1) and select "Remote-Containers: Reopen in Container"
3. Wait for the dev container to build and start

All terminal commands should be run in the VS Code terminal inside the dev container, not on your host machine.

# The code

This repository serves as a starter for this workshop, we already provide the
code for pulling an image from Docker Hub, pulling is rather involved and we
wanted you to be able to concentrate only on the runtime part of the container.

## Building and Running

Basic commands to get you started:

```console
# Build the project
make

# Run the project
sudo ./bin/devoxx-docker <commands>
```

Go ahead and run these commands now in your dev container terminal to make sure everything is set up correctly.

# Workshop Structure

The workshop is divided into 6 different parts. Parts 1, 2 and 3 are each
building upon the previous ones. Once you completed the first 3 parts you are
free to chose the last 3 in any order you like.

- 1. [Process Creation Basics](02-process-creation.md)
  - Creating a processes
- 2. [Namespace Isolation](03-namespace-isolation.md)
  - PID namespace isolation
  - UTS namespace for hostname isolation
- 3.  [Namespaces and Root Directory](04-namespaces-and-chroot.md)
  - Managing multiple namespaces
  - Implementing chroot
  - Directory structure setup
- 4. [Resource Control with cgroups](05-cgroups.md)
  - CPU limitations
  - Memory constraints
  - Process resource management
- 5. [Volume Management](06-volumes.md)
  - Implementing bind mounts
  - Volume persistence
  - Data sharing between host and container
- 6. [Network Configuration](07-network.md)
  - Network namespace setup
  - Virtual ethernet (veth) pairs
  - Basic networking capabilities

# Getting Help

- Use `make help` to see available commands
- Check the documentation in each exercise
- Refer to the hints and command references in each exercise file when stuck

[Start the workshop](02-process-creation.md)
