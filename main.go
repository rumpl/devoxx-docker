package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"github.com/rumpl/devoxx-docker/remote"
)

func main() {
	switch os.Args[1] {
	case "pull":
		if err := pull(os.Args[2]); err != nil {
			log.Fatal(err)
		}
	case "run":
		if err := run(os.Args[2:]); err != nil {
			log.Fatal(err)
		}
	case "child":
		if err := child(os.Args[2], os.Args[3], os.Args[4:]); err != nil {
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

func run(args []string) error {
	imageName := args[0]
	_, err := os.Stat("/fs/" + imageName)
	if err != nil {
		if os.IsNotExist(err) {
			if err := pull(imageName); err != nil {
				return fmt.Errorf("pull %w", err)
			}
		} else {
			return err
		}
	}

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, args...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	return cmd.Run()
}

func child(image string, command string, args []string) error {
	fmt.Printf("Running %s in %s\n", command, image)

	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	volumeDestination := fmt.Sprintf("/fs/%s/volume", image)
	if err := os.Mkdir(volumeDestination, 0755); err != nil {
		return fmt.Errorf("mkdir %w", err)
	}

	if err := syscall.Mount("/workspaces/devoxx-docker/volume", volumeDestination, "", syscall.MS_PRIVATE|syscall.MS_BIND, ""); err != nil {
		return fmt.Errorf("mount volume %w", err)
	}

	if err := syscall.Chroot("/fs/" + image); err != nil {
		return fmt.Errorf("chroot %w", err)
	}

	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("chdir %w", err)
	}

	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("mount proc %w", err)
	}

	if err := syscall.Sethostname([]byte("devoxx-container")); err != nil {
		return fmt.Errorf("set hostname %w", err)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %w", err)
	}

	return syscall.Unmount("proc", 0)
}
