# Mount namespace and root directory

Mount namespaces are a complex and powerful feature of the Linux kernel, we are
only scratching the surface here. We are going to implement the bare minimum but
feel free to explore and play with them.

Let's briefly talk about mount namespaces before writing the code.

A mount namespace isolates a process's (and its children) view of the
filesystem, it also isolates that process's mounts from other processes in the
system.

Here's an example, for this you will need two terminals, in the first terminal
type:

```console
$ docker run --rm -it --privileged alpine sh
# unshare --mount sh
# mount -t tmpfs tmpfs /mnt
# mount
...
...
tmpfs on /mnt type tmpfs (rw,relatime)
```

> Note: `unshare` is a Linux command that runs a program in a new namespace. In this case, it creates a new mount namespace and then runs the `sh` shell in it. This isolates the mount operations from the rest of the system.

In the second terminal first run `docker ps` and find the ID of the first
running container. Then you can:

```console
$ docker exec -it <container id> sh
# mount
...
...
```

Compare the two `mount` commands, if everything went well you should see that
the `tmpfs` mount that was made inside unshare of the container is not visible
in the container itself. Indeed `unshare --mount` creates a new (private) mount
namespace before running a command (`sh` in our case). Which effectively
separates the views between the unshared process and other processes.

Feel free to play around, see what happens if you create a new mount inside the
container, is it visible in the unshared process? What happens if you open a
third terminal, exec into the container and create a new mount namespace with
`sh`?

Let's now implement this in our container runtime!

# Step 1: add mount namespace

Modify the parent process to include mount namespace capability to the child
process:

```go
func parent() error {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[1:]...)...)

	// TODO:
	// 1. Add the mount namespace flag
}
```

<details>
<summary>Hint</summary>

Use `syscall.CLONE_NEWNS` for mount namespace isolation

</details>

# Step 2: pull the image

We provide the pulling functionality for you, now is the best time to use it.
Indeed in the next step we will already start to have everything we need to run
a command inside something that looks like a container, inside the base rootfs
of an image.

Here is how you should use the pulling functions:

```golang
image := "alpine"  // Note: this should be a string in quotes
// Create a puller for an image
puller := remote.NewImagePuller(image)
err := puller.Pull()
// check the error
```

Once you run this code, the image will be pulled from Docker Hub and its rootfs
extracted to `/fs/<image>`.

Create a function `pull(image string) error` that you can hook up to when the
command is `pull`.

> Note: The image pull command should be run with sudo from the dev container terminal.

Running the program after this step should look like:

```console
$ sudo ./bin/devoxx-docker pull alpine
Pulling alpine
Pulling done
$ sudo ls -l /fs/alpine/rootfs
... contents of the root filesystem of the alpine image ..
```

If you get an error running the pull command a second time, you may need to clean up the existing directories:

```console
$ sudo rm -rf /fs/alpine
```

# Step 3: change the root directory

Now that we have a root filesystem of an image, we can make the child use that
root filesystem as its root directory.

Implement container root directory setup:

```go
func setupContainer() error {
	// TODO:
	// 1. Print the current working directory
	// 2. Change root to "/fs/<image>/rootfs"
	// 3. Change current directory to root ("/")

	return nil
}
```

<details>
<summary>Hint</summary>

Look at `syscall.Chroot` and `os.Chdir` functions

</details>

> [!NOTE]
> Container systems don't use `chroot` because it does not really change
> the real root of the process, it only restricts the view of the process. a
> chroot can also be escaped from. We are using `chroot` because it's simpler, a
> nice extra exercice you can do is to make your program use `pivot_root`
> instead.

# Step 4: run a command in the child

With the current setup we have:

- namespace isolation
- a prepared rootfs

We basically have everything we need to run a command inside a real container
image we downloaded from Docker Hub. We are missing one last piece: running a
command in the child process, the container entrypoint if you will.

Write the needed code in the child process that will take the command to run
passed from the parent. Remember to hook up stdin and stdout just like you did in the parent process.

Once done you should be able to run:

```console
$ sudo ./bin/devoxx-docker run alpine /bin/sh
```

# Step 5: extra

Now that we have our container running, what happens when you type `ps`?
How could we fix that?

<details>
<summary>Hint</summary>

Look at the [default things](https://github.com/moby/moby/blob/6cbca96bfa3a2632e1636fb426ad69f9c38524d2/oci/defaults.go#L67-L110) that Docker defines for all containers, maybe take a couple?

You may need to use `syscall.Mount` to mount the `/proc` filesystem inside your container to make the `ps` command work correctly.

</details>

# Summary

We have now implemented mount namespace isolation and changed the root directory
for the container. This provides a contained filesystem environment for the
container.

# Additional Resources

- [man
  mount_namespaces](https://man7.org/linux/man-pages/man7/mount_namespaces.7.html)
- [man chroot](https://man7.org/linux/man-pages/man2/chroot.2.html)
- [Linux Filesystem Hierarchy
  Standard](https://refspecs.linuxfoundation.org/FHS_3.0/fhs/index.html)

[Previous step](./03-namespace-isolation.md)

Possible nexte steps, in any order:

- [Discovering cgroups](05-cgroups.md)
- [Implementing volume mounts](06-volumes.md)
- [Adding network support](07-network.md)

## Solution

<details>
<summary>Click to see the complete solution</summary>

```go
func main() {
    if len(os.Args) < 2 {
        log.Fatal("Not enough arguments")
    }

    switch os.Args[1] {
    case "child":
        if len(os.Args) < 3 {
            log.Fatal("Missing image name")
        }
        if err := child(os.Args[2]); err != nil {
            log.Fatal(err)
        }
    case "pull":
        if len(os.Args) < 3 {
            log.Fatal("Missing image name")
        }
        if err := pull(os.Args[2]); err != nil {
            log.Fatal(err)
        }
    case "run":
        if len(os.Args) < 4 {
            log.Fatal("Missing image name or command to run")
        }
        if err := run(); err != nil {
            log.Fatal(err)
        }
    default:
        log.Fatal("Unknown command", os.Args[1])
    }
}

func pull(image string) error {
    fmt.Printf("Pulling %s\n", image)
    puller := remote.NewImagePuller(image)
    if err := puller.Pull(); err != nil {
        return fmt.Errorf("pull failed: %w", err)
    }
    fmt.Println("Pulling done")
    return nil
}

func child(image string) error {
    fmt.Printf("CHILD PID: %d\n", os.Getpid())
    
    if err := syscall.Sethostname([]byte("container")); err != nil {
        return fmt.Errorf("sethostname failed: %w", err)
    }

    hostname, err := os.Hostname()
    if err != nil {
        return err
    }
    fmt.Printf("CHILD Hostname: %s\n", hostname)

    // Change root directory
    if err := syscall.Chroot(fmt.Sprintf("/fs/%s/rootfs", image)); err != nil {
        return fmt.Errorf("chroot failed: %w", err)
    }

    if err := syscall.Chdir("/"); err != nil {
        return fmt.Errorf("chdir failed: %w", err)
    }

    // Execute the command
    cmd := exec.Command(os.Args[3], os.Args[4:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    return cmd.Run()
}

func run() error {
    cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
    
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Add mount namespace along with existing namespaces
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
    }

    if err := cmd.Start(); err != nil {
        return fmt.Errorf("start failed: %w", err)
    }

    if err := cmd.Wait(); err != nil {
        return fmt.Errorf("wait failed: %w", err)
    }

    fmt.Printf("Container exited with code %d\n", cmd.ProcessState.ExitCode())
    return nil
}
```
</details>
