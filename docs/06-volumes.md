# Implementing volume mounts

As with cgroups, let's play a bit before implementing basic volumes for our
container.

If you've used docker at all you most certainly know that you can create a
`volume` with docker, it's a neat way to start a database for example and mount
that volume inside a container so that you don't lose your data when the
container stops. But how does Docker implement these volumes?

Let's create a volume first:

```console
$ docker volume create devoxx
```

Let's now go back to our PID 1 namespace as we did in the last exercice

```console
$ docker run -it --rm --privileged --pid=host justincormack/nsenter1
# ls -l /var/lib/docker/volumes/devoxx/_data/
total 0
```

Open a new terminal and run this command

```console
$ docker run --rm -v devoxx:/devoxx alpine sh -c 'echo "world" > /devoxx/hello'
```

And finally in the nsenter1 terminal

```console
# cat /var/lib/docker/volumes/devoxx/_data/hello
world
```

By this time you hopefully understand that docker volumes are nothing other than
special directories that live in a special, managed by docker, directory!

Let's see now how we can create something similar in our container runtime.

# Step 1: create the volume directory

In the parent process, create a directory inside the rootfs of the container, we
will use the `volume` directory present in this repository as our source volume.

```go
func setupVolume(volumePath, containerPath string) error {
	// TODO: Create a directory inside the rootfs of the container
	return nil
}
```

# Step 2: bind mount the volume

Create a function to handle bind mounting, make sure that you are using the
right flags, look at the different mount flags available, which ones should we
use? Where should the mount be made? In the parent or the child process?

```go
func mountVolume(source, target string) error {
	// TODO: Perform bind mount
	return nil
}
```

<details>
<summary>Hint</summary>

Use the `syscall.Mount` function

</details>

<details>
<summary>Hint</summary>

Don't forget to give the mount call the `syscall.MS_PRIVATE` flags, this ensures
that this mount stays private for our current mount namespace.

</details>

<details>
<summary>Hint</summary>

Since this mount is for the container, the mount should be done in the child
process i.e. in the process that lives in a new namespace.

</details>

# Step 3: unmount when done

Let's cleanup after all is done, we don't want to have dangling mounts all over
the place.

```go
func unmountVolume(target string) error {
	// TODO:
	// 1. Unmount the volume
	// 2. Handle any busy mount errors
	// 3. Clean up the mount point directory
	return nil
}
```

<details>
<summary>Hint</summary>

Look at `syscall.Unmount` function

</details>

# Step 4: test

1. Test your volume implementation:

```console
# Build the program
make

# Run with sudo
sudo ./bin/devoxx-container ...

# check the content of the mounted volume
```

# Summary

We have now implemented volume mounting functionality for containers using bind
mounts. This enables data persistence and sharing between the host and
container.

# Additional Resources

- [man mount](https://man7.org/linux/man-pages/man2/mount.2.html)
- [man umount](https://man7.org/linux/man-pages/man2/umount.2.html)
- [Linux bind
  mounts](https://man7.org/linux/man-pages/man8/mount.8.html#BIND_MOUNT_OPERATION)
- [Container volumes](https://docs.docker.com/storage/volumes/)

[Previous step](./05-cgroups.md) [Next step](07-network.md)
