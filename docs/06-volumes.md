# Implementing Volume Mounts for Containers

## Objective

Learn how to implement volume mounting functionality for containers using bind mounts. This exercise demonstrates how to share directories between the host and container, enabling data persistence and sharing.

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
        // Create the target directory
        if err := os.MkdirAll(target, 0755); err != nil {
            return fmt.Errorf("mkdir %w", err)
        }

        // Perform the bind mount
        if err := syscall.Mount(source, target, "", syscall.MS_BIND, ""); err != nil {
            return fmt.Errorf("bind mount %w", err)
        }

        return nil
    }
    ```

### Step 3: Add Volume Unmounting

1. Implement clean unmounting of volumes:
    ```go
    func unmountVolume(target string) error {
        // TODO:
        // 1. Unmount the volume using syscall.Unmount
        // 2. Handle any busy mount errors
        // 3. Clean up the mount point directory
        return nil
    }
    ```

### Step 4: Integration with Container Runtime

1. Add volume handling to your container creation flow:
    ```go
    func setupContainerVolumes(containerID string) error {
        volumes := []struct {
            source string
            target string
        }{
            {"/host/path", "/container/path"},
            // Add more volume mappings as needed
        }

        for _, vol := range volumes {
            if err := mountVolume(vol.source, vol.target); err != nil {
                return fmt.Errorf("mount volume %s: %w", vol.source, err)
            }
        }

        return nil
    }
    ```

### Step 5: Testing

1. Test your volume implementation:
    ```bash
    # Create test files in host volume
    echo "test data" > /path/to/host/volume/test.txt

    # Run container with volume
    sudo ./container run -v /path/to/host/volume:/container/volume ubuntu /bin/bash

    # Verify from inside container
    cat /container/volume/test.txt
    touch /container/volume/newfile.txt

    # Verify changes are visible on host
    ls -l /path/to/host/volume/newfile.txt
    ```

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
- [Linux bind mounts](https://man7.org/linux/man-pages/man8/mount.8.html#BIND_MOUNT_OPERATION)
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
```bash
# List mounts
mount | grep container-path

# Check mount points
findmnt

# Debug mount issues
dmesg | tail

# Check mount namespace
ls -l /proc/$PID/ns/mnt
```

## Error Handling Examples

```go
// Handle busy mount point
if err := syscall.Unmount(target, 0); err != nil {
    if err == syscall.EBUSY {
        // Handle busy mount point
        return fmt.Errorf("mount point is busy: %w", err)
    }
    return fmt.Errorf("unmount failed: %w", err)
}

// Handle non-existent source
if _, err := os.Stat(source); os.IsNotExist(err) {
    return fmt.Errorf("source path does not exist: %w", err)
}
```