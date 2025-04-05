# Implementing namespace isolation

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

# Step 1: prepare the child

To make sure we are really isolated, first add some logs to the child process,
print:

- the hostname
- the pid of the process

<details>
<summary>Hints</summary>

Use the `os` package to get the pid of the current process: `pid := os.Getpid()`

</details>

# Step 2: add namespace isolation

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

You need to set the `Cloneflags` to the `cmd`.

</details>

<details>
<summary>Hint 3 / Solution</summary>

```golang
cmd.SysProcAttr = &syscall.SysProcAttr {
    Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
}
```

</details>

# Step 3: set the hostname

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

# Step 5: test

1. Build and run your program:

```console
# Build the program
make

# Run with sudo
sudo ./bin/devoxx-container
```

# Summary

We have now implemented PID and UTS namespace isolation, providing process
isolation and custom hostname configuration for containers.  
This is a crucial step towards building a fully functional container runtime.

# Additional Resources

- [man namespaces](https://man7.org/linux/man-pages/man7/namespaces.7.html)
- [man clone](https://man7.org/linux/man-pages/man2/clone.2.html)
- [Go syscall package](https://pkg.go.dev/syscall)

[Previous step](./02-process-creation.md) [Next
step](04-namespaces-and-chroot.md)
