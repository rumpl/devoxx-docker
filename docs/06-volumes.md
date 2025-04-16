# Implementing volume mounts

As with cgroups, let's play a bit before implementing basic volumes for our
container.

If you've used docker at all you most certainly know that you can create a
`volume` with docker, it's a neat way to start a database for example and mount
that volume inside a container so that you don't lose your data when the
container stops. But how does Docker implement these volumes?

Let's create a volume first:

```console
$ docker volume create devoxx
```

Let's now go back to our PID 1 namespace as we did in the last exercise

```console
$ docker run -it --rm --privileged --pid=host justincormack/nsenter1
# ls -l /var/lib/docker/volumes/devoxx/_data/
total 0
```

Open a new terminal and run this command

```console
$ docker run --rm -v devoxx:/devoxx alpine sh -c 'echo "world" > /devoxx/hello'
```

And finally in the nsenter1 terminal

```console
# cat /var/lib/docker/volumes/devoxx/_data/hello
world
```

By this time you hopefully understand that docker volumes are nothing other than
special directories that live in a special, managed by docker, directory!

Let's see now how we can create something similar in our container runtime.

# Step 1: create the volume directory

In the parent process, create a directory inside the rootfs of the container, we
will use the `volume` directory present in this repository as our source volume.

```go
func setupVolume(containerPath string) error {
	// TODO: Create a directory inside the rootfs of the container
	return nil
}
```

# Step 2: bind mount the volume

Create a function to handle bind mounting, make sure that you are using the
right flags, look at the different mount flags available, which ones should we
use? Where should the mount be made? In the parent or the child process?

```go
func mountVolume(source, target string) error {
	// TODO: Perform bind mount
	return nil
}
```

<details>
<summary>Hint</summary>

Use the `syscall.Mount` function

</details>

<details>
<summary>Hint</summary>

Don't forget to give the mount call the `syscall.MS_PRIVATE` flags, this ensures
that this mount stays private for our current mount namespace.

</details>

<details>
<summary>Hint</summary>

Since this mount is for the container, the mount should be done in the child
process i.e. in the process that lives in a new namespace.

</details>

# Step 3: unmount when done

Let's cleanup after all is done, we don't want to have dangling mounts all over
the place.

```go
func unmountVolume(target string) error {
	// TODO:
	// 1. Unmount the volume
	// 2. Handle any busy mount errors
	// 3. Clean up the mount point directory
	return nil
}
```

<details>
<summary>Hint</summary>

Look at `syscall.Unmount` function

</details>

# Step 4: test

1. Test your volume implementation:

```console
# Build the program
make

# Run with sudo
sudo ./bin/devoxx-docker ...

# check the content of the mounted volume
```

If everything works correctly, the file you created in the container should be visible in the volume directory on the host.

# Summary

We have now implemented volume mounting functionality for containers using bind
mounts. This enables data persistence and sharing between the host and
container.

# Additional Resources

- [man mount](https://man7.org/linux/man-pages/man2/mount.2.html)
- [man umount](https://man7.org/linux/man-pages/man2/umount.2.html)
- [Linux bind
  mounts](https://man7.org/linux/man-pages/man8/mount.8.html#BIND_MOUNT_OPERATION)
- [Container volumes](https://docs.docker.com/storage/volumes/)

[Previous step](./05-cgroups.md) [Next step](07-network.md)

## Solution

<details>
<summary>Click to see the complete solution</summary>

```go
const (
	CGROUP_ROOT = "/sys/fs/cgroup"
	MEMORY_MAX  = "104857600"    // 100MB memory limit
	CPU_MAX     = "50000 100000" // 50ms per 100ms period
	VOLUME_ROOT = "/volumes"     // Base directory for volumes
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

func setupVolume(containerPath string) error {
	// Create the volume directory if it doesn't exist
	if err := os.MkdirAll(containerPath, 0755); err != nil {
		return fmt.Errorf("failed to create volume directory: %w", err)
	}

	fmt.Printf("Created volume directory at %s\n", containerPath)
	return nil
}

func mountVolume(source, target string) error {
	// Ensure target directory exists
	if err := os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %w", err)
	}

	// Perform bind mount
	flags := syscall.MS_BIND | syscall.MS_REC | syscall.MS_PRIVATE
	if err := syscall.Mount(source, target, "", uintptr(flags), ""); err != nil {
		return fmt.Errorf("failed to bind mount volume: %w", err)
	}

	fmt.Printf("Mounted volume from %s to %s\n", source, target)
	return nil
}

func unmountVolume(target string) error {
	// Try to unmount
	if err := syscall.Unmount(target, syscall.MNT_DETACH); err != nil {
		if err == syscall.EBUSY {
			// If mount is busy, retry with force unmount
			fmt.Printf("Mount point busy, attempting force unmount of %s\n", target)
			if err := syscall.Unmount(target, syscall.MNT_FORCE); err != nil {
				return fmt.Errorf("failed to force unmount volume: %w", err)
			}
		} else {
			return fmt.Errorf("failed to unmount volume: %w", err)
		}
	}

	// Clean up the mount point directory
	if err := os.RemoveAll(target); err != nil {
		return fmt.Errorf("failed to remove mount point directory: %w", err)
	}

	fmt.Printf("Unmounted and cleaned up volume at %s\n", target)
	return nil
}

func child(image string) error {
	// Print the PID of the current process
	fmt.Println("CHILD: Hello from child, my pid is", os.Getpid())

	// Print a simple message
	fmt.Println("Hello from child")

	// Set container hostname
	if err := syscall.Sethostname([]byte("container")); err != nil {
		return err
	}

	// Print new hostname to verify the change
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	fmt.Printf("CHILD Hostname: %s\n", hostname)

	// Set up volume mounts if specified
	if len(os.Args) > 4 && os.Args[4] == "-v" {
		volumeSpec := os.Args[5]
		parts := strings.Split(volumeSpec, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid volume specification: %s", volumeSpec)
		}

		source := filepath.Join(VOLUME_ROOT, parts[0])
		target := filepath.Join("/", parts[1])

		if err := mountVolume(source, target); err != nil {
			return err
		}

		// Register cleanup handler
		defer unmountVolume(target)
	}

	// Execute the command
	cmd := exec.Command(os.Args[3], os.Args[4:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func run() error {
	// Create a unique volume path for this container
	volumePath := filepath.Join(VOLUME_ROOT, fmt.Sprintf("vol-%d", time.Now().UnixNano()))

	if err := setupVolume(volumePath); err != nil {
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

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("wait failed: %w", err)
	}

	fmt.Printf("Container exited with code %d\n", cmd.ProcessState.ExitCode())
	return nil
}
```
</details>
