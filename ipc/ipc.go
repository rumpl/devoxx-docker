package ipc

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

// CreateSocketPair creates a socketpair and returns the file descriptors
func CreateSocketPair() (parent *os.File, child *os.File, err error) {
	// Create a pair of connected Unix domain sockets
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, nil, fmt.Errorf("socketpair: %w", err)
	}

	// Convert the file descriptors to *os.File
	parent = os.NewFile(uintptr(fds[0]), "parent-socket")
	child = os.NewFile(uintptr(fds[1]), "child-socket")
	return parent, child, nil
}

// SendReady sends a ready message through the socket
func SendReady(f *os.File) error {
	_, err := f.Write([]byte("READY"))
	return err
}

// WaitForReady blocks until a ready message is received from the socket
func WaitForReady(f *os.File) error {
	buf := make([]byte, 5) // "READY" is 5 bytes
	_, err := io.ReadFull(f, buf)
	if err != nil {
		return fmt.Errorf("wait for ready: %w", err)
	}
	if string(buf) != "READY" {
		return fmt.Errorf("unexpected message: %s", string(buf))
	}
	return nil
}
