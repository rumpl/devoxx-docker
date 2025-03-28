# Discovering cgroups Configuration

## Objective

In this exercise, you will learn how to configure cgroups to limit memory and CPU usage for a process. You will focus on setting `memory.max`, `cpu.max`, and adding the process to `cgroup.procs`.

## Steps

### Step 1: Setup cgroups

1. Create a new directory for the cgroup:
    - Add the following code to create the cgroup directory:
      ```go
      cgroupPath := "/sys/fs/cgroup/devoxx-container"
      if err := os.Mkdir(cgroupPath, 0755); err != nil {
          log.Fatalf("Failed to create cgroup directory: %v", err)
      }
      ```

### Step 2: Configure Memory Limit

1. Set the memory limit to 100MB:
    - Add the following code to set the memory limit:
      ```go
      if err := os.WriteFile(cgroupPath+"/memory.max", []byte("104857600"), 0644); err != nil {
          log.Fatalf("Failed to set memory limit: %v", err)
      }
      ```

### Step 3: Configure CPU Limit

1. Set the CPU limit to 50ms per 100ms:
    - Add the following code to set the CPU limit:
      ```go
      if err := os.WriteFile(cgroupPath+"/cpu.max", []byte("50000 100000"), 0644); err != nil {
          log.Fatalf("Failed to set CPU limit: %v", err)
      }
      ```

### Step 4: Add Process to cgroup

1. Add the process to the cgroup:
    - Add the following code to add the process to the cgroup:
      ```go
      if err := os.WriteFile(cgroupPath+"/cgroup.procs", []byte(strconv.Itoa(pid)), 0644); err != nil {
          log.Fatalf("Failed to add process to cgroup: %v", err)
      }
      ```

### Step 5: Run the Program

1. Compile and run the Go program.
    - Save the code in a file named `main.go`.
    - Open a terminal and navigate to the directory containing `main.go`.
    - Run the following commands to compile and execute the program:
      ```sh
      go build -o devoxx-cgroup
      sudo ./devoxx-cgroup
      ```

2. Observe the output to see the cgroup configuration and how it limits the process resources.

## Hints

- Use the `os.Mkdir` function to create the cgroup directory.
- Use the `os.WriteFile` function to configure `memory.max`, `cpu.max`, and `cgroup.procs`.

## Key Points

- Understand how to configure cgroups to limit memory and CPU usage.
- Learn how to use the `os` package in Go to manipulate cgroups.
- Observe the effect of cgroup configuration on process resource usage.

## Additional Resources

- [man cgroups](https://man7.org/linux/man-pages/man7/cgroups.7.html)
- [Go os package documentation](https://pkg.go.dev/os)
