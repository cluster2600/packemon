// +build darwin

package packemon

import (
	"fmt"
	"os/exec"
	"strings"
)

// TCProgramManager manages the TCP RST packet handling on macOS
type TCProgramManager struct {
	interfaceName string
	filterRules   []string
	isActive      bool
}

// NewTCProgramManager creates a new TCP program manager for macOS
func NewTCProgramManager(interfaceName string) (*TCProgramManager, error) {
	return &TCProgramManager{
		interfaceName: interfaceName,
		filterRules:   make([]string, 0),
		isActive:      false,
	}, nil
}

// Start sets up packet filtering rules to drop TCP RST packets on macOS
func (t *TCProgramManager) Start() error {
	if t.isActive {
		return nil // Already active
	}

	// Check if pfctl is available (required for packet filtering on macOS)
	if _, err := exec.LookPath("pfctl"); err != nil {
		return fmt.Errorf("pfctl not found, packet filtering unavailable: %v", err)
	}

	// Create a rule to drop outgoing TCP RST packets
	rule := fmt.Sprintf("block drop out proto tcp from any to any flags R/R on %s", t.interfaceName)
	t.filterRules = append(t.filterRules, rule)

	// Create a temporary pf.conf file with our rules
	tempRules := fmt.Sprintf("# Packemon TCP RST blocking rules\n%s\n", strings.Join(t.filterRules, "\n"))
	
	// Write rules to a temporary file
	tempFile, err := createTempFile("packemon-pf-", ".conf", tempRules)
	if err != nil {
		return fmt.Errorf("failed to create temporary pf rules file: %v", err)
	}

	// Load the rules
	cmd := exec.Command("sudo", "pfctl", "-f", tempFile, "-e")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to load packet filter rules: %v", err)
	}

	t.isActive = true
	return nil
}

// Stop removes the packet filtering rules on macOS
func (t *TCProgramManager) Stop() error {
	if !t.isActive {
		return nil // Not active
	}

	// Disable the packet filter
	cmd := exec.Command("sudo", "pfctl", "-d")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to disable packet filter: %v", err)
	}

	t.isActive = false
	t.filterRules = make([]string, 0)
	return nil
}

// createTempFile creates a temporary file with the given content
func createTempFile(prefix, suffix, content string) (string, error) {
	// Use the mktemp command to create a temporary file
	cmd := exec.Command("mktemp", "-t", prefix, suffix)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	tempFile := strings.TrimSpace(string(output))

	// Write content to the file
	writeCmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' > %s", content, tempFile))
	if err := writeCmd.Run(); err != nil {
		return "", err
	}

	return tempFile, nil
}
