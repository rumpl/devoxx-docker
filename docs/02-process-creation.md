# Create a processes

Technically, at its lowest level a container is a Linux process isolated using
`cgroups` (control groups) and `namespaces` (like PID, net, mount, user, etc.)
to restrict resource usage and provide a separate view of the system.

Let's start small and first only create a new process.

# Step 1: open the main.go file

First, open the main.go file in your editor:

```console
# Make sure you're in the dev container terminal
cd /workspaces/devoxx-docker
code main.go
```

Now set up the basic program structure:

```go
func main() {
	// TODO: Check if we're running the initial command or the child process
	// If args contain "child", call child()
	// Otherwise, continue with parent process creation
}
```

# Step 2: the child process

Create the child process handler:

```go
func child() error {
	// TODO:
	// 1. Print the PID of the current process
	// 2. Execute the desired command, printing a simple `Hello from child` is enough for now
}
```

# Step 3: the parent process

Create a function to handle parent process logic:

```golang
func run() error {
	// TODO:
	// 1. Print the PID of the current process
	// 2. Create a new process using current executable
	// 3. Set up stdin/stdout/stderr
	// 4. Start the child process
	// 5. Wait for completion and print a message letting us know the child process has exited
}
```

<details>
<summary>Hints</summary>

- Use `os.Getpid()` to get the pid of the current process
- Use `/proc/self/exe` to re-execute the same process
- Use `os.Args` to detect if running as child
- Use `cmd.Start()` and `cmd.Wait()` for better process control
- To hook up stdin/stdout/stderr correctly, use:
  - `cmd.Stdin = os.Stdin`
  - `cmd.Stdout = os.Stdout`
  - `cmd.Stderr = os.Stderr`
    This allows interactive commands to work properly.

</details>

# Step 4: implement the code in the comments

Look at the TODOs in the code snippets above and implement each function. Remember to follow the hints in the previous sections.

# Step 5: test

1. Build and run your program in the dev container terminal:

```console
# Make sure you're in the dev container terminal
# Build the program
make

# Run the program (still in the dev container terminal)
./bin/devoxx-docker
PARENT: Hello from parent, my pid is 1234
CHILD: Hello from child, my pid is 1325
PARENT: Child exited with exit code 0
```

# Summary

We have the basic first step into our journey to creating a container, we have a
parent process that can manage the child process. This child process will soon
become a real container.

# Additional Resources

- [Go os/exec package](https://pkg.go.dev/os/exec)
- [Go os package](https://pkg.go.dev/os)

[Previous step](./01-intro.md) [Next step](./03-namespace-isolation.md)


## Solution

<details>
<summary>Click to see the complete solution</summary>

```go
func main() {
    if len(os.Args) < 2 {
        if err := run(); err != nil {
            log.Fatal(err)
        }
        os.Exit(0)
    }
    switch os.Args[1] {
    case "child":
		if err := child(); err != nil {
			log.Fatal(err)
		}
    default:
        log.Fatal("Unknown command", os.Args[1])
    }
}

func child() error {
    // Print the PID of the current process
    fmt.Println("CHILD: Hello from child, my pid is", os.Getpid())
    
    // Print a simple message
    fmt.Println("Hello from child")
    return nil
}

func run() error {
    // Print the PID of the current process
    fmt.Println("PARENT: Hello from parent, my pid is", os.Getpid())
    
    // Create a new process using current executable
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[1:]...)...)
    
    // Set up stdin/stdout/stderr
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Add namespace flags
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWPID,
    }

    // Start the child process
    if err := cmd.Start(); err != nil {
        return err
    }

    // Wait for completion and print exit status
    if err := cmd.Wait(); err != nil {
        return err
    }
    fmt.Println("PARENT: Child exited with exit code", cmd.ProcessState.ExitCode())
    
    return nil
}
```
</details>
