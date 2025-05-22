# Packemon for macOS

This document provides instructions for running Packemon on macOS, either natively or using Docker.

## Native Installation

Packemon now includes native support for macOS! The implementation uses platform-specific code to handle networking operations that are compatible with macOS.

### Prerequisites

- macOS 10.15 (Catalina) or newer
- Go 1.24 or newer
- Homebrew (for installing dependencies)
- Administrator privileges (for packet capture)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ddddddO/packemon.git
   cd packemon
   ```

2. Run the macOS installation script:
   ```bash
   ./darwin_install.sh
   ```

   This script will:
   - Check for required dependencies
   - Install libpcap if needed
   - Build Packemon with macOS-specific code

3. Run Packemon:
   ```bash
   # Monitor mode
   sudo ./packemon
   
   # Generator mode
   sudo ./packemon --send
   ```

### Implementation Details

The macOS implementation differs from the Linux version in several key ways:

1. **Network Interface Handling**: Uses `libpcap` instead of raw sockets for packet capture and injection
2. **TCP RST Packet Handling**: Uses macOS's Packet Filter (PF) instead of eBPF
3. **Platform-Specific Code**: Uses build tags to separate macOS and Linux implementations

## Docker Installation (Alternative)

If you prefer to run Packemon in Docker, you can use the provided Docker setup script:

```bash
./macos_docker_setup.sh
```

This will build a Docker image and provide instructions for running Packemon in a container.

## Limitations

While the macOS implementation provides most of the functionality available in the Linux version, there are some limitations:

1. **Performance**: Packet processing may be slightly slower due to the use of libpcap instead of direct raw socket access
2. **Privileges**: Administrator privileges (sudo) are required to capture and send packets
3. **Network Interface Access**: Some advanced network interface operations may be limited

## Troubleshooting

### Permission Issues

If you encounter permission issues when running Packemon, make sure you're using `sudo` to run the application.

### Packet Capture Issues

If packet capture isn't working:

1. Check that your user has permission to access network devices
2. Verify that libpcap is installed correctly: `brew info libpcap`
3. Try specifying a different network interface with `--interface`

### Build Issues

If you encounter build issues:

1. Make sure Go 1.24+ is installed: `go version`
2. Verify that libpcap is installed: `brew list libpcap`
3. Check for any compiler errors in the output

## Contributing

If you encounter issues with the macOS implementation or have suggestions for improvements, please open an issue or submit a pull request on GitHub.
