package packemon

import (
	"fmt"
)

// TCProgramManager interface for platform-specific implementations
type TCProgramManagerInterface interface {
	Start() error
	Stop() error
}

// NewTCProgramManager creates a new TCP program manager
// The implementation is platform-specific and is defined in:
// - tc_program_linux.go for Linux
// - tc_program_darwin.go for macOS
func NewTCProgramManager(interfaceName string) (TCProgramManagerInterface, error) {
	return newTCProgramManagerPlatform(interfaceName)
}
