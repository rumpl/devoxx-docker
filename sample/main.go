package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/moby/sys/mountinfo"
)

func pids() error {
	color.Cyan("\nPIDS")
	dirs, err := os.ReadDir("/proc")
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 8, ' ', 0)
	fmt.Fprintln(tw, "NAME\tPID")
	for _, dir := range dirs {
		if pid, err := strconv.Atoi(dir.Name()); err == nil {
			b, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
			if err != nil {
				return err
			}

			fmt.Fprintf(tw, "%s\t%d\n", strings.TrimSpace(string(b)), pid)
		}
	}
	tw.Flush()

	return nil
}

func mounts() error {
	color.Cyan("\nMOUNTS")
	mis, err := mountinfo.GetMounts(nil)
	if err != nil {
		return err
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "MOUNT POINT\tROOT\tFS TYPE")

	for _, mi := range mis {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", mi.Mountpoint, mi.Root, mi.FSType)
	}
	tw.Flush()
	return nil
}

func main() {
	fmt.Println("My pid is", os.Getpid())

	if err := pids(); err != nil {
		log.Fatal(err)
	}

	if err := mounts(); err != nil {
		log.Fatal(err)
	}
}
