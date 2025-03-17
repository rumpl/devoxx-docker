# Process namespace

In the sample directory, build the program:

```console
$ go build .
```

Run the binary:

```console
$ ./devoxx-sample
My pid is 213123

PIDS
NAME                 PID
docker-init          1
bash                 1001
gopls                14454
...
```

As expected, the output shows all the processes currently running and the pid
of the current process is some random process id.

Let's now create a new process namespace and run our application again

```console
$ sudo unshare --fork --pid --mount-proc ./devoxx-sample
My pid is 1

PIDS
NAME                 PID
devoxx-sample        1
```

It's important that you read `man unshare` but let's explain the command quickly

- `--fork` tells `unshare` to fork before calling the binary
- `--pid` tells `unshare` to create a new pid namespace
- `--mount-proc` means that `unshare` will mount a new proc filesystem

What happened? Since we created a new process namespace and asked `unshare` to
mount that new namespace our application only sees what exists in this new
namespace. Since our application is the only one existing in this namespace, its
PID is 1, and it only sees itself.


Note: for process namespaces you can also use normal bash commands to see the
same thing:

- run `ps aux`, you see all the processes, which is normal
- run `sudo unshare --pid --fork --mount --mount-proc /bin/bash`,you only see
 `bash` and `ps`, which is great.

# Mount namespace

TODO

# UTS namespace

UTS  namespaces provide isolation of two system identifiers: the hostname and
the NIS domain name.

Let's call `unshare` again and tell it to create a new UTS namespace:

```console
$ sudo unshare --pid --fork --mount --uts --mount-proc /bin/sh
```

Once we are in the new shell, we can safely call `hostname foo` and this will
change the host name _in this namespace_. If you open a new terminal inside the
dev container and type `hostname` you will see a different name, i.e. the change
made inside the namespace is visible only in that namespace.

# User namespace

TODO, maybe?
