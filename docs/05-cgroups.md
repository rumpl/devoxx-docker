# Discovering cgroups

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

# Step 1: Create a cgroup for the container

This should be done in the parent process before starting the container:

```go
func setupCgroups(childPid string) error {
	// TODO:
	// 1. Create base cgroup directory under the "/sys/fs/cgroup" directory
    // 	For example: "/sys/fs/cgroup/devoxx-docker/<childPid>"
	// 2. Set appropriate permissions (0755)

	return nil
}
```

# Step 2: Configure memory limit

Set the memory limit to 100MB:

```go
func setupCgroups(childPid string) error {
	// TODO:
	// 1. Create the file to set the memory limit
	// 2. Write the limit value (100MB) to the file

        return nil
}
```

# Step 3: Configure CPU limit

Set the CPU limit to 50ms per 100ms:

```go
func setupCgroups() error {
	// TODO:
	// 1. Create the file to set the CPU limit
	// 2. Write the limit value (50ms per 100ms) to the file

        return nil
}
```

# Step 4: Add process to the cgroup

Add the process to the cgroup. This must be done in the parent process after starting the child but before waiting for it to complete:

```go
func addProcessToCgroup(pid int) error {
	// TODO:
	// 1. Get the PID of the child process
	// 2. Create the file to add the process to the cgroup
    // 	The file is: "<cgroup_path>/cgroup.procs"

	return nil
}
```

# Step 5: Testing cgroups

To verify your cgroup implementation works correctly:

1. Add some logging to show the memory and CPU limits you've set
2. Try running a memory-intensive workload in your container:

```console
# Make sure you're in the dev container terminal
$ sudo ./bin/devoxx-docker run alpine /bin/sh
# dd if=/dev/zero of=/dev/null bs=1M count=200
```

If your cgroup memory limit is working, this should either run slower or fail with an out-of-memory error.

# Summary

We have now implemented cgroup configuration to limit memory and CPU usage for
the container process. This provides resource management capabilities for
containers.

# Additional Resources

- [man cgroups](https://man7.org/linux/man-pages/man7/cgroups.7.html)

[Previous step](./04-namespaces-and-chroot.md) [Next step](06-volumes.md)

## Solution

<details>
<summary>Click to see the complete solution</summary>

```go
const (
	CGROUP_ROOT = "/sys/fs/cgroup"
	MEMORY_MAX  = "104857600"    // 100MB memory limit
	CPU_MAX     = "50000 100000" // 50ms per 100ms period
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
	}

	switch os.Args[1] {
	case "child":
		if len(os.Args) < 3 {
			log.Fatal("Missing image name")
		}
		if err := child(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	case "pull":
		if len(os.Args) < 3 {
			log.Fatal("Missing image name")
		}
		if err := pull(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	case "run":
		if len(os.Args) < 4 {
			log.Fatal("Missing image name or command")
		}
		if err := run(); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Unknown command", os.Args[1])
	}
}

func pull(image string) error {
	fmt.Printf("Pulling %s\n", image)
	puller := remote.NewImagePuller(image)
	if err := puller.Pull(); err != nil {
		return fmt.Errorf("pull failed: %w", err)
	}
	fmt.Println("Pulling done")
	return nil
}

func child(image string) error {
	fmt.Printf("CHILD PID: %d\n", os.Getpid())

	if err := syscall.Sethostname([]byte("container")); err != nil {
		return fmt.Errorf("sethostname failed: %w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	fmt.Printf("CHILD Hostname: %s\n", hostname)

	// Change root directory
	if err := syscall.Chroot(fmt.Sprintf("/fs/%s/rootfs", image)); err != nil {
		return fmt.Errorf("chroot failed: %w", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir failed: %w", err)
	}

	// Execute the command
	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func setupCgroups() error {
	// Create base cgroup directory
	cgroupPath := filepath.Join(CGROUP_ROOT, "devoxx-docker")
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		return fmt.Errorf("failed to create cgroup directory: %w", err)
	}

	// Set memory limit
	memoryMaxPath := filepath.Join(cgroupPath, "memory.max")
	if err := os.WriteFile(memoryMaxPath, []byte(MEMORY_MAX), 0644); err != nil {
		return fmt.Errorf("failed to set memory limit: %w", err)
	}

	// Set CPU limit
	cpuMaxPath := filepath.Join(cgroupPath, "cpu.max")
	if err := os.WriteFile(cpuMaxPath, []byte(CPU_MAX), 0644); err != nil {
		return fmt.Errorf("failed to set CPU limit: %w", err)
	}

	fmt.Printf("Created cgroup at %s with memory limit %s and CPU limit %s\n",
		cgroupPath, MEMORY_MAX, CPU_MAX)

	return nil
}

func addProcessToCgroup(pid int) error {
	cgroupPath := filepath.Join(CGROUP_ROOT, "devoxx-docker")
	procsPath := filepath.Join(cgroupPath, "cgroup.procs")

	// Write PID to cgroup.procs
	if err := os.WriteFile(procsPath, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		return fmt.Errorf("failed to add process to cgroup: %w", err)
	}

	fmt.Printf("Added process %d to cgroup %s\n", pid, cgroupPath)
	return nil
}

func run() error {
	// Set up cgroups before starting the container
	if err := setupCgroups(); err != nil {
		return err
	}

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start failed: %w", err)
	}

	// Add the process to cgroup after starting but before waiting
	if err := addProcessToCgroup(cmd.Process.Pid); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("wait failed: %w", err)
	}

	fmt.Printf("Container exited with code %d\n", cmd.ProcessState.ExitCode())
	return nil
}
```
</details>
