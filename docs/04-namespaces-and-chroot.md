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
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[	:]...)...)

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
image := alpine
// Create a puller for an image
puller := remote.NewImagePuller(image)
err := puller.Pull()
// check the error
```

Once you run this code, the image will be pulled from Docker Hub and its rootfs
extracted to `/fs/<image>`.

Create a function `pull(image string) error` that you can hook up to when the
command is `pull`.

Running the program after this step should look like:

```console
$ sudo ./bin/devoxx-docker pull alpine
Pulling alpine
Pulling done
$ ls -l /fs/alpine/rootfs
... contents of the root filesystem of the alpine image ..
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
passed from the parent.

Once done you should be able to run:

```console
$ sudo ./bin/devoxx-container run alpine /bin/sh
```

# Step 4: extra

Now that we have our container running, what happens when you type `ps`?
How could we fix that?

<details>
<summary>Hint</summary>

Look at the [default things](https://github.com/moby/moby/blob/6cbca96bfa3a2632e1636fb426ad69f9c38524d2/oci/defaults.go#L67-L110) that Docker defines for all containers, maybe take a couple?

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

[Previous step](./03-namespace-isolation.md) [Next step](05-cgroups.md)
