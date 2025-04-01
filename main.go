package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	"github.com/rumpl/devoxx-docker/cgroups"
	"github.com/rumpl/devoxx-docker/ipc"
	"github.com/rumpl/devoxx-docker/mount"
	"github.com/rumpl/devoxx-docker/net"
	"github.com/rumpl/devoxx-docker/remote"
)

func main() {
	switch os.Args[1] {
	case "pull":
		if err := pull(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	case "run":
		if err := parent(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
	case "child":
		socketFD, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		socket := os.NewFile(uintptr(socketFD), "child-socket")
		if socket == nil {
			log.Fatal("Failed to create socket file from FD")
		}

		if err := child(socket, os.Args[3], os.Args[4], os.Args[5:]); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Unknown command %s\n", os.Args[1])
	}
}

func pull(image string) error {
	fmt.Printf("Pulling %s\n", image)
	puller := remote.NewImagePuller(image)
	err := puller.Pull()
	fmt.Println("Pulled image")
	return err
}

func parent(args []string) error {
	image := args[0]
	_, err := os.Stat("/fs/" + image)
	if err != nil {
		if os.IsNotExist(err) {
			if err := pull(image); err != nil {
				return fmt.Errorf("pull %w", err)
			}
		} else {
			return err
		}
	}

	if err := os.MkdirAll("/fs/"+image+"/rootfs/etc", 0755); err != nil {
		return fmt.Errorf("create etc dir: %w", err)
	}

	if err := os.WriteFile("/fs/"+image+"/rootfs/etc/resolv.conf", []byte("nameserver 1.1.1.1\n"), 0644); err != nil {
		return fmt.Errorf("write resolv.conf: %w", err)
	}

	parentSocket, childSocket, err := ipc.CreateSocketPair()
	if err != nil {
		return fmt.Errorf("create socket pair: %w", err)
	}
	defer parentSocket.Close()

	childFD := childSocket.Fd()
	socketFD := strconv.Itoa(int(childFD))

	cmd := exec.Command("/proc/self/exe", append([]string{"child", socketFD}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.ExtraFiles = []*os.File{childSocket}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start %w", err)
	}

	childSocket.Close()

	vethName := "veth0"
	if err := net.SetupVeth(vethName, cmd.Process.Pid); err != nil {
		return fmt.Errorf("setup veth %w", err)
	}
	defer func() {
		if err := net.CleanupVeth(vethName); err != nil {
			log.Printf("cleanup veth %s", err)
		}
	}()

	if err := cgroups.SetupCgroup(cmd.Process.Pid); err != nil {
		return fmt.Errorf("setup cgroup %w", err)
	}
	defer func() {
		if err := cgroups.RemoveCgroup(cmd.Process.Pid); err != nil {
			log.Printf("remove cgroup %s", err)
		}
	}()

	// We are done setting up things, tell the child to continue
	if err := ipc.SendReady(parentSocket); err != nil {
		return fmt.Errorf("send ready: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("wait %w", err)
	}

	fmt.Printf("Container exited with exit code %d\n", cmd.ProcessState.ExitCode())

	return err
}

func child(socket *os.File, image string, command string, args []string) error {
	fmt.Printf("Running %s in %s\n", command, image)
	defer socket.Close()

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	volumeDestination := fmt.Sprintf("/fs/%s/rootfs/volume", image)
	if err := os.MkdirAll(volumeDestination, 0755); err != nil {
		return fmt.Errorf("mkdir %w", err)
	}

	if err := syscall.Mount("/workspaces/devoxx-docker/volume", volumeDestination, "", syscall.MS_PRIVATE|syscall.MS_BIND, ""); err != nil {
		return fmt.Errorf("mount volume %w", err)
	}

	if err := syscall.Chroot("/fs/" + image + "/rootfs"); err != nil {
		return fmt.Errorf("chroot %w", err)
	}

	// Change to the root directory
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("chdir %w", err)
	}

	if err := mount.Mount(); err != nil {
		return fmt.Errorf("mount %w", err)
	}

	if err := syscall.Sethostname([]byte("devoxx-container")); err != nil {
		return fmt.Errorf("set hostname %w", err)
	}

	// Wait for parent to complete network setup
	if err := ipc.WaitForReady(socket); err != nil {
		return fmt.Errorf("wait for network setup: %w", err)
	}

	peerName := "veth1"
	if err := net.SetupContainerNetworking(peerName); err != nil {
		return fmt.Errorf("setup container networking %w", err)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("start %w", err)
	}

	if err := mount.Unmount(); err != nil {
		return fmt.Errorf("unmount %w", err)
	}

	return nil
}
