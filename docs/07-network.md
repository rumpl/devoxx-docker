# Adding network support

Our container runtime creates "containers", great, and they can even connect to
the internet, you can try this by running:

```console
$ sudo ./bin/devoxx-container run alpine /bin/sh
# ping 1.1.1.1
```

But what happens if we try and ping google.com for example? Of course, it
doesn't work, the base alpine image doesn't contain the `/etc/resolv.conf` file,
we don't know where to look for when we want to find the IP address of a host.

A second issue is that, while the container has networking, it's using the host
network stack, and we don't want that, we want to isolate the processes in our
container as much as we can. This includes the network stack.

There are many ways we could setup networking for containers, in this workshop
we will try and keep it simple (but buckle up, nothing is easy when it comes to
networking). We will create a pair of virtual Ethernet devices (`veth`) and link
them.

> [!NOTE]
> While old, outdated and deprecated, we _will_ be using `iptables` in
> this exercice. If you feel lucky you could also implement this functionnality
> with `nftables`.

# Step 1: create the network namespace

1. Add network namespace isolation in your main container creation code:

```go
cmd.SysProcAttr = &syscall.SysProcAttr{
	// Add CLONE_NEWNET to your existing clone flags
	// This creates a new network namespace for the container
}
```

After adding the new network namespace flag, can we ping anything from inside
the container?

# Step 2: create the veth pair

Of course, when we created a new network namespace for our container, we lost
all connectivity! This is normal and expected, we need to setup everything
manually on the host and inside the container so that we can have network
connectivity.

- The container interface should be in the 10.0.0.0/24 subnet
- Common IP assignments:
  - Host interface (veth0): 10.0.0.1
  - Container interface (veth1): 10.0.0.2
- Required iptables rules should enable NAT for the container subnet

```go
func SetupVeth(vethName string, pid int) error {
	// TODO: Use "ip link" commands to:
	// 1. Create a veth pair (veth0 and veth1)
	// 2. Move veth1 to the container's network namespace
	// 3. Configure veth0 in the host namespace
	// 4. Set up NAT rules using iptables
}
```

<details>
<summary>Hint (veth)</summary>

This command creates a veth pair. A kind of virtual network cable with two ends

```console
ip link add veth0 type veth peer name veth1
```

</details>

<details>
<summary>Hint (move to network namespace)</summary>

This command moves a veth to the network namespace of a process

```console
ip link set veth1 netns <PID>
```

</details>

<details>
<summary>Hint (assign ip address)</summary>

Assign an IP address and subnet to our veth

```console
ip addr add 10.0.0.1/24 dev veth0
```

</details>

<details>
<summary>Hint (NAT rule)</summary>

```console
iptables -t nat -A POSTROUTING -s 10.0.0.0/24 -j MASQUERADE
```

- `-t nat` is the Network Address Translation table, it's used for rewriting
  packed addressses
- `-A POSTROUTING` adds the rule to this chain, which alters the packets just
  before the packet leaves the system
- `-s 10.0.0.0/24` mathes packets with a source IP in this subnet
- `-j MASQUERADE` means it should _masquerade_ the packet: replace its source IP
  with the host's outgoing interface IP

</details>

Create a cleanup function to remove the network configuration:

```go
func CleanupVeth(vethName string) error {
	// TODO: Clean up:
	// 1. Remove NAT rules
	// 2. Delete the veth pair
}
```

# Step 3: configure container networking

We need to do a couple things before getting our networking connection:

- assign an IP address to our veth (`veth1`)
- bring it `up`
- add a default route. This route is the path that packets take when they don't
  match any other routes, in our case the default route should be the IP address
  of the veth that is on the host.

```go
func SetupContainerNetworking(peerName string) error {
	// TODO: Inside the container:
	// 1. Assign IP address to the container interface
	// 2. Bring up the interface
	// 3. Configure the loopback interface
	// 4. Set up the default route
}
```

<details>
<summary>Loopback interface</summary>

While this isn't really needed you can run this command inside the container

```console
ip link set lo up
```

This sets up the loopback interface and makes it possible to `ping 127.0.0.1`
for example

</details>

<details>
<summary>Default gateway</summary>

```console
ip route add default via 10.0.0.1
```

</details>

# Step 3: DNS

It's always DNS, right?

We should hopefully have networking working now but one last little bit remains,
we can't `ping google.com`, we need to create the `/etc/resolv.conf` file inside
the rootfs of the container.

# Step 4: test

1. Test your network implementation:

```console
$ sudo ./bin/devoxx-container run alpine /bin/sh
# ping 1.1.1.1      # Should reach internet
# ping google.com   # Should resolve and reach
```

# Additional Resources

- [man ip-netns](https://man7.org/linux/man-pages/man8/ip-netns.8.html)
- [man veth](https://man7.org/linux/man-pages/man4/veth.4.html)
- [man iptables](https://man7.org/linux/man-pages/man8/iptables.8.html)
- [Linux Network
  Namespaces](https://man7.org/linux/man-pages/man7/network_namespaces.7.html)
- [Container Networking](https://docs.docker.com/network/)
- [Reference Container Network Interface
  plugins](https://github.com/containernetworking/plugins?tab=readme-ov-file#main-interface-creating)

[Previous step](06-volumes.md) [Next step](08-outro.md)
