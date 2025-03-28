# Managing Namespaces and Mounts

## Objective

In this exercise, you will learn how to manage namespaces and mounts in a containerized environment. You will focus on adding a mount, changing the root directory using `chroot`, and changing to the root directory using `Chdir`.

## Steps

### Step 1: Add a Mount

1. Use the `syscall.Mount` function to add a mount.
    - Add the following code to mount the volume:
      ```go
      volumeDestination := "/fs/container/rootfs/volume"
      if err := os.MkdirAll(volumeDestination, 0755); err != nil {
          log.Fatalf("Failed to create volume directory: %v", err)
      }
      if err := syscall.Mount("/workspaces/devoxx-docker/volume", volumeDestination, "", syscall.MS_PRIVATE|syscall.MS_BIND, ""); err != nil {
          log.Fatalf("Failed to mount volume: %v", err)
      }
      ```

### Step 2: Change Root Directory

1. Use the `syscall.Chroot` function to change the root directory.
    - Add the following code to change the root directory:
      ```go
      if err := syscall.Chroot("/fs/container/rootfs"); err != nil {
          log.Fatalf("Failed to change root directory: %v", err)
      }
      ```

### Step 3: Change to Root Directory

1. Use the `os.Chdir` function to change to the root directory.
    - Add the following code to change to the root directory:
      ```go
      if err := os.Chdir("/"); err != nil {
          log.Fatalf("Failed to change to root directory: %v", err)
      }
      ```

### Step 4: Run the Program

1. Compile and run the Go program.
    - Save the code in a file named `main.go`.
    - Open a terminal and navigate to the directory containing `main.go`.
    - Run the following commands to compile and execute the program:
      ```sh
      go build -o devoxx-container
      sudo ./devoxx-container
      ```

2. Observe the output to see the process IDs and how they are isolated in the new namespace.

## Hints

- Use the `syscall.Mount` function to add a mount.
- Use the `syscall.Chroot` function to change the root directory.
- Use the `os.Chdir` function to change to the root directory.

## Key Points

- Understand how to create and manage process namespaces.
- Learn how to use the `syscall` package in Go to manipulate namespaces.
- Observe the isolation of processes in a new namespace.
- Learn how to add mounts and change the root directory in a containerized environment.

## Additional Resources

- [man unshare](https://man7.org/linux/man-pages/man1/unshare.1.html)
- [Go syscall package documentation](https://pkg.go.dev/syscall)
- [Go exec package documentation](https://pkg.go.dev/os/exec)
