package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {

	switch os.Args[1] {
	case "run":
		if err := run(); err != nil {
			panic(err)
		}
	case "child":
		if err := child(); err != nil {
			panic(err)
		}
	default:
		fmt.Println("Hello, World!")
	}
}

func run() error {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}

	return cmd.Run()
}

func child() error {
	fmt.Printf("Running %v \n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := syscall.Sethostname([]byte("container")); err != nil {
		return fmt.Errorf("set hostname %w", err)
	}
	if err := syscall.Chroot("/fs/ubuntu"); err != nil {
		return fmt.Errorf("chroot %w", err)
	}
	if err := os.Chdir("/"); err != nil {
		return fmt.Errorf("chdir %w", err)
	}
	if err := syscall.Mount("proc", "proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("mount proc %w", err)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run %w", err)
	}

	return syscall.Unmount("proc", 0)
}
