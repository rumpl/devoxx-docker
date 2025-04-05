# Implementing Volume Mounts for Containers

## Objective

Learn how to implement volume mounting functionality for containers using bind
mounts. This exercise demonstrates how to share directories between the host and
container, enabling data persistence and sharing.

## Steps

### Step 1: Create Volume Directory Structure

1. Set up the volume directories:

```go
func setupVolume(volumePath, containerPath string) error {
	// TODO:
	// 1. Create the source volume directory on host if it doesn't exist
	// 2. Create the target mount point in container
	// 3. Ensure proper permissions (0755)
	return nil
}
```

### Step 2: Implement Bind Mount

1. Create a function to handle bind mounting:

```go
func mountVolume(source, target string) error {
	// TODO:
	// 1. Check if source and target paths exist
	// 2. Create target directory if it doesn't exist
	// 3. Perform bind mount
	// 4. Handle any errors
	return nil
}
```

<details>
<summary>Hint</summary>
Look at `syscall.Mount` function
</details>

### Step 3: Add Volume Unmounting

1. Implement clean unmounting of volumes:

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

### Step 4: Integration with Container Runtime

1. Add volume handling to your container creation flow:

```go
func setupContainerVolumes(containerID string) error {
	// TODO:
	// 1. Define volume mappings
	// 2. Mount each volume
	// 3. List content of mounted volume
	

	return nil
}
```
<details>
<summary>Hint</summary>
You can create an array of volume mappings and iterate over them to mount each volume.
</details>

### Step 5: Testing

1. Test your volume implementation:

```console
# Build the program
make

# Run with sudo (needed for namespace operations)
sudo ./bin/devoxx-container

# check the content of the mounted volumes
```

### Summary

We have now implemented volume mounting functionality for containers using bind
mounts. This enables data persistence and sharing between the host and container.

[Previous step](./05-cgroups.md) [Next step](07-ipc.md)

## Hints

- Use `syscall.Mount()` with `MS_BIND` flag for bind mounts
- Always create target directories before mounting
- Remember to handle unmounting during container cleanup
- Use `defer` for cleanup operations
- Check for existing mounts before mounting
- Ensure proper error handling and cleanup on failures

## Key Points

- Bind mounts create a view of a host directory in the container
- Proper cleanup is essential to avoid orphaned mounts
- Volume paths must exist before mounting
- Changes in mounted volumes are immediately visible in both host and container
- Mount flags affect the behavior of the mounted volume

## Additional Resources

- [man mount](https://man7.org/linux/man-pages/man2/mount.2.html)
- [man umount](https://man7.org/linux/man-pages/man2/umount.2.html)
- [Linux bind
  mounts](https://man7.org/linux/man-pages/man8/mount.8.html#BIND_MOUNT_OPERATION)
- [Container volumes](https://docs.docker.com/storage/volumes/)

## Command Reference

### Mount Operations

```go
// Basic bind mount
syscall.Mount(source, target, "", syscall.MS_BIND, "")

// Bind mount with additional flags
syscall.Mount(source, target, "", syscall.MS_BIND|syscall.MS_REC, "")

// Unmount
syscall.Unmount(target, 0)
```

### Directory Operations

```go
// Create mount point
os.MkdirAll(path, 0755)

// Check if directory exists
if _, err := os.Stat(path); os.IsNotExist(err) {
    // Directory doesn't exist
}

// Remove mount point
os.RemoveAll(path)
```

### Debugging Commands

```console
# List mounts
mount | grep container-path

# Check mount points
findmnt

# Debug mount issues
dmesg | tail

# Check mount namespace
ls -l /proc/$PID/ns/mnt
```
