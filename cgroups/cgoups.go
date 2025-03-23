package cgroups

import (
	"os"
	"strconv"
)

func SetupCgroup(pid int) error {
	cgroupPath := "/sys/fs/cgroup/mytoycontainer"
	if err := os.Mkdir(cgroupPath, 0755); err != nil {
		return err
	}

	// Limit memory to 100MB
	if err := os.WriteFile(cgroupPath+"/memory.max", []byte("104857600"), 0644); err != nil {
		return err
	}

	// Limit CPU (50ms per 100ms)
	if err := os.WriteFile(cgroupPath+"/cpu.max", []byte("50000 100000"), 0644); err != nil {
		return err
	}

	// Add process to cgroup
	return os.WriteFile(cgroupPath+"/cgroup.procs", []byte(strconv.Itoa(pid)), 0644)
}

func RemoveCgroup(pid int) error {
	cgroupPath := "/sys/fs/cgroup/mytoycontainer"
	return os.RemoveAll(cgroupPath)
}
