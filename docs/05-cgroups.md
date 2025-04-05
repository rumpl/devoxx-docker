# Discovering cgroups

## Objective

In this exercise, you will learn how to configure cgroups to limit memory and
CPU usage for a process. You will focus on setting `memory.max`, `cpu.max`, and
adding the process to `cgroup.procs`.

## Steps

### Step 1: Setup cgroups

1. Create a new directory for the cgroup:

```go
func child() error {
	// TODO:
	// 1. Create base cgroup directory under the "/fs/cgroup" directory
	// 2. Set appropriate permissions (0755)
	// 3. Handle all potential errors

	return nil
}
```

### Step 2: Configure Memory Limit

1. Set the memory limit to 100MB:

```go
func child() error {
	// TODO:
	// 1. Create the file to set the memory limit
	// 2. Write the limit value (100MB) to the file
	// 3. Handle all potential errors
}
```

### Step 3: Configure CPU Limit

1. Set the CPU limit to 50ms per 100ms:

```go
func child() error {
	// TODO:
	// 1. Create the file to set the CPU limit
	// 2. Write the limit value (50ms per 100ms) to the file
	// 3. Handle all potential errors
}
```

### Step 4: Add Process to cgroup

1. Add the process to the cgroup:

```go
func child() error {
	// TODO:
	// 1. Get the PID of the current process
	// 2. Create the file to add the process to the cgroup
	// 3. Write the PID to the file
	// 4. Handle all potential errors
}
```

### Step 5: Testing

1. Build and run your program:

```console
# Build the program
make

# Run with sudo (needed for namespace operations)
sudo ./bin/devoxx-container
```

### Summary

We have now implemented cgroup configuration to limit memory and CPU usage for
the container process. This provides resource management capabilities for
containers.

[Previous step](./04-namespace-and-chroot.md) [Next step](06-volumes.md)

## Additional Resources

- [man cgroups](https://man7.org/linux/man-pages/man7/cgroups.7.html)
- [Go os package documentation](https://pkg.go.dev/os)
