// networkinterface.go is a stub file that imports the platform-specific implementation

package packemon

import (
	"context"
	"net"
)

// NewNetworkInterface creates a new NetworkInterface for the specified interface
// The implementation is platform-specific and is defined in:
// - networkinterface_linux.go for Linux
// - networkinterface_darwin.go for macOS
func NewNetworkInterface(nwInterface string) (*NetworkInterface, error) {
	// Each platform implements this function differently
	// The actual implementation is in the platform-specific files
	return newNetworkInterfacePlatform(nwInterface)
}

// SendEthernetFrame sends an Ethernet frame
func (nwif *NetworkInterface) SendEthernetFrame(ctx context.Context, data []byte) error {
	return nwif.sendEthernetFramePlatform(ctx, data)
}

// ReceiveEthernetFrame receives Ethernet frames
func (nwif *NetworkInterface) ReceiveEthernetFrame(ctx context.Context) {
	nwif.receiveEthernetFramePlatform(ctx)
}

// GetNetworkInfo returns information about the network interface
func (nwif *NetworkInterface) GetNetworkInfo() (macAddr net.HardwareAddr, ipv4Addr net.IP, ipv6Addr net.IP) {
	return nwif.getNetworkInfoPlatform()
}

// Close cleans up resources
func (nwif *NetworkInterface) Close() {
	nwif.closePlatform()
}
