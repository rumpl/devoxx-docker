# Creating parent and child processes

## Objective

Technically, at its lowest level a container is a Linux process isolated using
`cgroups` (control groups) and `namespaces` (like PID, net, mount, user, etc.)
to restrict resource usage and provide a separate view of the system.

Let's start small and first only create a new process.

## Steps

### Step 1: create the main function

Set up the basic program structure:

```go
func main() {
   // TODO: Check if we're running the initial command or the child process
   // If args contain "child", call child()
   // Otherwise, continue with parent process creation
}
```

### Step 2: implement the child process

Create the child process handler:

```go
func child() error {
   // TODO:
   // 1. Print the PID of the current process
   // 2. Execute the desired command, printing a simple `Hello from child` is enough for now
}
```

### Step 3: Implement Parent Process Creation

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

</details>

### Step 4: Testing

1. Build and run your program:

```console
# Build the program
go build -o devoxx-container

# Run the program
./devoxx-container
PARENT: Hello from parent, my pid is 1234
CHILD: Hello from child, my pid is 1325
PARENT: Child exited with exit code 0
```

### Summary

We have the basic first step into our journey to creating a container, we have a
parent process that can manage the child process. This child process will soon
become a real container.

[Previous step](./01-intro.md) [Next step](./03-namespace-isolation.md)

## Additional Resources

- [Go os/exec package](https://pkg.go.dev/os/exec)
- [Go os package](https://pkg.go.dev/os)
