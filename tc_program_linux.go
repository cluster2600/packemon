// +build linux

package packemon

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cilium/ebpf"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -Wall -Werror" tc_program ./tc_program.bpf.c

// TCProgramManager manages the TCP RST packet handling on Linux
type TCProgramManager struct {
	interfaceName string
	qdisc         netlink.Qdisc
	filter        netlink.Filter
	objs          tc_programObjects
	isActive      bool
}

// newTCProgramManagerPlatform creates a new TCP program manager for Linux
func newTCProgramManagerPlatform(interfaceName string) (TCProgramManagerInterface, error) {
	return &TCProgramManager{
		interfaceName: interfaceName,
		isActive:      false,
	}, nil
}

// Start sets up eBPF program to drop TCP RST packets on Linux
func (t *TCProgramManager) Start() error {
	if t.isActive {
		return nil // Already active
	}

	// Load pre-compiled eBPF program
	if err := loadTc_programObjects(&t.objs, nil); err != nil {
		return fmt.Errorf("loading objects: %w", err)
	}

	// Get network interface
	link, err := netlink.LinkByName(t.interfaceName)
	if err != nil {
		return fmt.Errorf("getting interface %s: %w", t.interfaceName, err)
	}

	// Add clsact qdisc
	qdisc := &netlink.GenericQdisc{
		QdiscAttrs: netlink.QdiscAttrs{
			LinkIndex: link.Attrs().Index,
			Handle:    netlink.MakeHandle(0xffff, 0),
			Parent:    netlink.HANDLE_CLSACT,
		},
		QdiscType: "clsact",
	}

	if err := netlink.QdiscAdd(qdisc); err != nil {
		return fmt.Errorf("adding clsact qdisc: %w", err)
	}
	t.qdisc = qdisc

	// Add filter for egress
	filterAttrs := netlink.FilterAttrs{
		LinkIndex: link.Attrs().Index,
		Parent:    netlink.HANDLE_MIN_EGRESS,
		Handle:    netlink.MakeHandle(0, 1),
		Protocol:  unix.ETH_P_ALL,
		Priority:  1,
	}

	filter := &netlink.BpfFilter{
		FilterAttrs:  filterAttrs,
		Fd:           t.objs.TcDropRst.FD(),
		Name:         "tc_drop_rst",
		DirectAction: true,
	}

	if err := netlink.FilterAdd(filter); err != nil {
		netlink.QdiscDel(qdisc)
		return fmt.Errorf("adding eBPF filter: %w", err)
	}
	t.filter = filter
	t.isActive = true

	return nil
}

// Stop removes the eBPF program on Linux
func (t *TCProgramManager) Stop() error {
	if !t.isActive {
		return nil // Not active
	}

	// Remove filter
	if err := netlink.FilterDel(t.filter); err != nil {
		return fmt.Errorf("deleting filter: %w", err)
	}

	// Remove qdisc
	if err := netlink.QdiscDel(t.qdisc); err != nil {
		return fmt.Errorf("deleting qdisc: %w", err)
	}

	// Close eBPF objects
	if err := t.objs.Close(); err != nil {
		return fmt.Errorf("closing objects: %w", err)
	}

	t.isActive = false
	return nil
}
