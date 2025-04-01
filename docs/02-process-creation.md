# Creating Parent and Child Processes

## Objective

Technically, a container is a Linux process isolated using `cgroups` (control
groups) and `namespaces` (like PID, net, mount, user, etc.) to restrict resource
usage and provide a separate view of the system.

Let's start small and first only create a new process.

## Steps

### Step 1: Create the Main Function

1. Set up the basic program structure:

```go
func main() {
   //TODO: Check if we're running the initial command or the child process
   // If args contain "child", call child()
   // Otherwise, continue with parent process creation
}
```

### Step 2: Implement Child Process

1. Create the child process handler:

```go
func child() error {
   //TODO:
   // 1. Print current PID to demonstrate namespace isolation
   // 2. Execute the desired command, a simple `Hello from child` print to the console is enough for now
   // 3. Keep process running to observe isolation
}
```

### Step 3: Implement Parent Process Creation

1. Create a function to handle parent process logic: `go func run() error {
//TODO: // 1. Create a new command using current executable // 2. Set up
    stdin/stdout/stderr // 3. Start the child process // 4. Wait for completion
       and print a message letting us know the child process has exited } `
<details>
<summary>Hints</summary>

- Use `/proc/self/exe` to re-execute the same process
- Use `os.Args` to detect if running as child
- Remember to handle all potential errors
- Use `cmd.Start()` and `cmd.Wait()` for better process control

</details>

### Step 4: Testing

1. Build and run your program:

```bash
# Build the program
go build -o devoxx-container

# Run the program
./devoxx-container
```

### Summary

We have the basic first step into our journey to creating a container, we have a
parent process that can manage the child process. This child process will soon
become a real container.

[Next step](03-namespace-isolation.md)

## Additional Resources

- [Go os/exec package](https://pkg.go.dev/os/exec)
- [Go os package](https://pkg.go.dev/os)

## Command Reference

### Process Information

```bash
# View process tree
ps -ef --forest

# Get current process info
ps -p $$

# View process environment
ps eww -p <pid>
```

### Common Operations

```go
// Get current PID
pid := os.Getpid()

// Get parent PID
ppid := os.Getppid()

// Create command with arguments
cmd := exec.Command("program", "arg1", "arg2")

// Run command and wait for completion
err := cmd.Run()
```

[Previous step](./01-intro.md) [Next step](./03-namespace-isolation.md)
