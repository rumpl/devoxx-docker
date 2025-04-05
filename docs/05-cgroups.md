# Discovering cgroups

## Objective

Control groups (cgroups) are a Linux kernel feature that allows you to group
processes and then limit, prioritize, or monitor their usage of system resources
like CPU, memory, I/O, etc.

In this exercise, you will learn how to configure cgroups to limit memory and
CPU usage for a process. You will focus on setting `memory.max`, `cpu.max`, and
adding the process to `cgroup.procs`.

Let's explore how cgroups work, this is by no means a comprehensive tutorial,
only a brief introduction.

Let's browse the cgroups in general and those created by docker, to do that, we
need to use the `justincormack/nsenter1` image that helps you get inside the
namespace of the process with PID 1, which is where `dockerd` runs.

```console
$ docker run -it --rm --privileged --pid=host justincormack/nsenter1
# ls -l /sys/fs/cgroup
# ls -l /sys/fs/cgroup/docker
```

What do we see here. We see a cgroups (v2) hierarchy, `dockerd` creates one
cgroup for itself and creates, for each container, a new cgroup in
`/sys/fs/cgroup/docker/<container id>`. Let's test this out!

In a new terminal run this:

```console
$ docker run --rm -d --cpus 2 nginx
<container id>
```

In the `nsenter1` terminal, what do we see?

```console
# cat /sys/fs/cgroup/docker/<container id>/cpu.max
200000 100000
# cat /sys/fs/cgroup/docker/cpu.max
max 100000
```

So, by telling docker that we want 2 cpus for this container, it created a new
cgroup for that container and limited the amount of cpus processes in that
cgroup can use.

Neat, now let's implement this in our container runtime.

## Steps

### Step 1: Create a cgroup for the container

Create a new directory for the cgroup:

```go
func child() error {
	// TODO:
	// 1. Create base cgroup directory under the "/sys/fs/cgroup" directory
	// 2. Set appropriate permissions (0755)

	return nil
}
```

### Step 2: Configure memory limit

Set the memory limit to 100MB:

```go
func child() error {
	// TODO:
	// 1. Create the file to set the memory limit
	// 2. Write the limit value (100MB) to the file
	// 3. Handle all potential errors
}
```

### Step 3: Configure CPU limit

Set the CPU limit to 50ms per 100ms:

```go
func child() error {
	// TODO:
	// 1. Create the file to set the CPU limit
	// 2. Write the limit value (50ms per 100ms) to the file
	// 3. Handle all potential errors
}
```

### Step 4: Add process to the cgroup

Add the process to the cgroup:

```go
func child() error {
	// TODO:
	// 1. Get the PID of the current process
	// 2. Create the file to add the process to the cgroup
	// 3. Write the PID to the file
	// 4. Handle all potential errors
}
```

### Summary

We have now implemented cgroup configuration to limit memory and CPU usage for
the container process. This provides resource management capabilities for
containers.

[Previous step](./04-namespace-and-chroot.md) [Next step](06-volumes.md)

## Additional Resources

- [man cgroups](https://man7.org/linux/man-pages/man7/cgroups.7.html)
