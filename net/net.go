package net

import (
	"fmt"
	"os/exec"
)

func SetupVeth(vethName string, pid int) error {
	// Create veth pair
	peerName := "veth1"
	if err := exec.Command("ip", "link", "add", vethName, "type", "veth", "peer", "name", peerName).Run(); err != nil {
		return fmt.Errorf("create veth pair %w", err)
	}

	// Move peer end into container network namespace
	if err := exec.Command("ip", "link", "set", peerName, "netns", fmt.Sprintf("%d", pid)).Run(); err != nil {
		return fmt.Errorf("move veth to netns %w", err)
	}

	// Setup host end
	if err := exec.Command("ip", "addr", "add", "10.0.0.1/24", "dev", vethName).Run(); err != nil {
		return fmt.Errorf("add ip to host veth %w", err)
	}
	if err := exec.Command("ip", "link", "set", vethName, "up").Run(); err != nil {
		return fmt.Errorf("set host veth up %w", err)
	}

	// Setup NAT
	if err := exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-s", "10.0.0.0/24", "-j", "MASQUERADE").Run(); err != nil {
		return fmt.Errorf("setup NAT %w", err)
	}

	return nil
}

func CleanupVeth(vethName string) error {
	// Cleanup NAT rule
	if err := exec.Command("iptables", "-t", "nat", "-D", "POSTROUTING", "-s", "10.0.0.0/24", "-j", "MASQUERADE").Run(); err != nil {
		return fmt.Errorf("cleanup NAT rule %w", err)
	}

	// Cleanup veth pair
	if err := exec.Command("ip", "link", "delete", vethName).Run(); err != nil {
		return fmt.Errorf("cleanup veth pair %w", err)
	}

	return nil
}

func SetupContainerNetworking(peerName string) error {
	// This command assigns the IP address 10.0.0.2 with a subnet mask of /24 (255.255.255.0)
	// to the network interface named in peerName (which is "veth1").
	// This sets up the container's network interface with an IP address in the 10.0.0.0/24 subnet,
	// allowing it to communicate with the host (which has 10.0.0.1) and potentially the outside world.
	// The command fails if the IP address cannot be assigned, which could happen if the interface
	// doesn't exist or if there's already an IP address conflict.
	if err := exec.Command("ip", "addr", "add", "10.0.0.2/24", "dev", peerName).Run(); err != nil {
		return fmt.Errorf("add ip to peer %w", err)
	}

	// This command activates the network interface named in peerName (which is "veth1").
	// The "ip link set <interface> up" command brings the interface to the "up" state,
	// making it operational so it can send and receive network traffic.
	// Without this step, the interface would remain in the "down" state and wouldn't function.
	// The command fails if the interface can't be activated, which could happen if
	// the interface doesn't exist or if there are permission issues.
	if err := exec.Command("ip", "link", "set", peerName, "up").Run(); err != nil {
		return fmt.Errorf("set peer up %w", err)
	}

	// This command adds a default route to the network interface named in peerName (which is "veth1").
	// The default route is the path that packets take when they don't match any other routes.
	// In this case, it's the route to the host's IP address (10.0.0.1),
	// which allows the container to communicate with the host.
	// The command fails if the route can't be added, which could happen if
	// the interface doesn't exist or if there are permission issues.
	if err := exec.Command("ip", "route", "add", "default", "via", "10.0.0.1").Run(); err != nil {
		return fmt.Errorf("add route %w", err)
	}

	return nil
}
