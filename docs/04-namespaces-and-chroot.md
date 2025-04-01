# Managing Namespaces and Root Directory

## Objective

Learn how to manage filesystem isolation in a containerized environment by
implementing mount namespaces and changing the root directory using `chroot`.  
This exercise demonstrates how to create a contained filesystem environment.

## Steps

### Step 1: Add Mount Namespace

1. Modify the parent process to include mount namespace capability to the child
   process:

```go
func parent() error {
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[1:]...)...)

    // TODO:
    // 1. Add the mount namespace flag
}
```

### Step 2: Setup Root Directory Structure

1.  Create the function to setup root filesystem:

```go
func child() error {
    // TODO:
    // 1. Create base rootfs directory at "/fs/container/rootfs"
    // 2. Set appropriate permissions (0755)
    // 3. Handle all potential errors

    return nil
}
```

<details>
<summary>Hint</summary>
Look at `os.MkdirAll` function
</details>

### Step 3: Change Root Directory

1.  Implement container root directory setup:

```go
func setupContainer() error {
    // TODO:
    // 1. Print the current working directory
    // 2. Change root to "/fs/container/rootfs"
    // 3. Change current directory to root ("/")
    // 4. Handle all potential errors
    // 5. Implement proper error handling
    // 6. Print the new working directory

        return nil
    }
```

<details>
<summary>Hint</summary>
Look at `syscall.Chroot` and `os.Chdir` functions
</details>

### Step 4: Testing

1. Build and run your program:

```bash
# Build the program
make

# Run with sudo (needed for namespace operations)
sudo ./bin/devoxx-container
```

### Summary

We have now implemented mount namespace isolation and changed the root directory
for the container.
This provides a contained filesystem environment for the container.

[Next step](05-cgroups.md)

## Hints

- Use `syscall.CLONE_NEWNS` for mount namespace isolation
- Root privileges are required for namespace operations
- Use absolute paths when working with directories
- Remember to handle cleanup in case of errors
- Check if directories exist before operations
- Use `defer` for cleanup operations

## Key Points

- Mount namespaces provide filesystem isolation
- `chroot` changes the root directory view
- Proper cleanup is essential to avoid resource leaks
- Namespace operations require careful error handling

## Additional Resources

- [man
  mount_namespaces](https://man7.org/linux/man-pages/man7/mount_namespaces.7.html)
- [man chroot](https://man7.org/linux/man-pages/man2/chroot.2.html)
- [Linux Filesystem Hierarchy
  Standard](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/index.html)

## Command Reference

### Namespace Operations

```go
// Create mount namespace
syscall.CLONE_NEWNS

// Change root
syscall.Chroot(path)
```

### Debugging Commands

```bash
# Check filesystem structure
ls -la /fs/container/rootfs

# View mount namespaces
ls -l /proc/$$/ns/mnt

# View process namespaces
ls -l /proc/$$/ns/
```

### Error Handling Examples

```go
// Handle chroot errors
if err := syscall.Chroot("/fs/container/rootfs"); err != nil {
    if os.IsPermission(err) {
        return fmt.Errorf("chroot permission denied (run with sudo): %w", err)
    }
    return fmt.Errorf("chroot failed: %w", err)
}

// Handle directory operations
if err := os.MkdirAll("/path/to/dir", 0755); err != nil {
    return fmt.Errorf("failed to create directory: %w", err)
}
```

```

```
