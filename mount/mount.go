package mount

import (
	"fmt"
	"syscall"
)

func Mount() error {
	if err := syscall.Mount("proc", "/proc", "proc", 0, ""); err != nil {
		return fmt.Errorf("mount proc %w", err)
	}

	if err := syscall.Mount("sysfs", "/sys", "sysfs", 0, ""); err != nil {
		return fmt.Errorf("mount sys %w", err)
	}
	if err := syscall.Mount("cgroup2", "/sys/fs/cgroup", "cgroup2", 0, ""); err != nil {
		return fmt.Errorf("mount cgroup2 %w", err)
	}

	if err := syscall.Mount("dev", "/dev", "devtmpfs", 0, ""); err != nil {
		return fmt.Errorf("mount dev %w", err)
	}

	return nil
}

func Unmount() error {
	if err := syscall.Unmount("/proc", 0); err != nil {
		return fmt.Errorf("unmount proc %w", err)
	}

	if err := syscall.Unmount("/sys/fs/cgroup", 0); err != nil {
		return fmt.Errorf("unmount cgroup %w", err)
	}

	if err := syscall.Unmount("/sys", 0); err != nil {
		return fmt.Errorf("unmount sys %w", err)
	}

	if err := syscall.Unmount("/dev", 0); err != nil {
		return fmt.Errorf("unmount dev %w", err)
	}

	return nil
}
