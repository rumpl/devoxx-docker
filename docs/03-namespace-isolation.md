# Implementing Namespace Isolation

## Objective

The next natural step towards something that resembles a real container is
isolating the process from the other processes in the system. In Linux this is
done thanks to namespaces, there are different namespaces provided by the Linux
kernel:

- UTS namespace (for setting a new hostname without it being global)
- PID namespace, the processes inside a PID namespaces can only see the
  processes inside that namespace
- Network namespace ...
- etc.

We will only look at the UTS and PID namespace in this exercise.

## Steps

### Step 1: Prepare the child

To make sure we are really isolated, first add some logs to the child process,
print:

- the hostname
- the pid of the process

<details>
<summary>Hints</summary>

Use the `os` package to get the pid of the current process. pid := os.Getpid()

</details>

### Step 2: Add Namespace Isolation

1.  Modify the parent process creation to include namespace flags:

```go
func run() error {
cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args...)...)

            //TODO:
        // 1. Add namespace flags for PID and UTS namespaces

        if err := cmd.Wait(); err != nil {
            return fmt.Errorf("wait %w", err)
        }

        fmt.Printf("Container exited with exit code %d\n", cmd.ProcessState.ExitCode())
    }
```

    <details>
    <summary>Hint</summary>
    Look at the `SysProcAttr` property of the `exec.Cmd` structure
    </details>

<details>
<summary>Hint 2</summary>
You need to set both `Cloneflags` and `Unshareflags`
</details>

<details>
<summary>Hint 3 / Solution</summary>
cmd.SysProcAttr = &syscall.SysProcAttr {
    Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
    UnshareFlags: syscall.CLONE_NEWNS,
}
</details>

### Step 3: Implement Hostname Changes

Now that the child lives in its own new host and pid namespaces, we can set the
hostname _for that namespace_ and also take a look at our pid, if everything
went well, the child pid should be 1.

1. Add hostname configuration to the child process:

```go
func child() error {
    //TODO:
    // 1. Set the container hostname
    // 2. Print the hostname to verify the change
}
```

<details>
<summary>Hint</summary>
Look at `syscall.Sethostname` function
</details>

### Step 5: Testing

1. Build and run your program:

```bash
# Build the program
make

# Run with sudo (needed for namespace operations)
sudo ./bin/devoxx-container
```

### Summary

We have now implemented PID and UTS namespace isolation, providing process
isolation and custom hostname configuration for containers.  
This is a crucial step towards building a fully functional container runtime.

[Next step](04-namespaces-and-chroot.md)

## Key Points

- PID namespace provides process isolation
- UTS namespace enables custom hostname
- Namespace changes require root privileges
- Child process sees itself as PID 1

## Additional Resources

- [man namespaces](https://man7.org/linux/man-pages/man7/namespaces.7.html)
- [man clone](https://man7.org/linux/man-pages/man2/clone.2.html)
- [Go syscall package](https://pkg.go.dev/syscall)

## Command Reference

### Namespace Operations

```go
// Create new namespaces
cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWPID | syscall.CLONE_NEWUTS,
}

// Set hostname
syscall.Sethostname([]byte("new-hostname"))
```

### Debugging Commands

```bash
# Check process namespaces
ls -l /proc/$$/ns/

# View hostname
hostname

# Check PID in different namespaces
ps aux
```

### Error Handling Examples

```go
// Handle hostname errors
if err := syscall.Sethostname([]byte("container-host")); err != nil {
    if os.IsPermission(err) {
        return fmt.Errorf("permission denied: run with sudo: %w", err)
    }
    return fmt.Errorf("hostname error: %w", err)
}
```
