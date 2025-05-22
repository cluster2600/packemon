# Packemon: macOS Support Implementation Guide

## Current Status

After analyzing the codebase, I can confirm that Packemon currently does not support macOS natively. As mentioned in the README:

> This tool is not available for Windows and macOS. I have confirmed that it works on Linux (Debian and Ubuntu on WSL2).

## Implementation Challenges

Packemon uses Linux-specific networking features that are not directly compatible with macOS:

1. Raw socket implementation using `golang.org/x/sys/unix`
2. Linux-specific network interface handling
3. eBPF programs for TCP RST packet handling

## Options for Running on macOS

### Option 1: Run via Docker (Recommended)

This is the simplest approach to get Packemon running on macOS without code modifications:

1. Install Docker Desktop for Mac
2. Build the Docker image using the provided Dockerfile
3. Run the container with appropriate networking privileges

```bash
# Build the Docker image
docker build -t packemon .

# Run Packemon in Docker
docker run --rm -it --privileged --network host packemon
```

Note: The `--network host` option works differently on macOS than on Linux. You might experience limitations with network interface access.

### Option 2: Run in a Linux VM

1. Install a Linux VM using VirtualBox, Parallels, or VMware
2. Build Packemon inside the VM
3. Run with sudo privileges in the VM

### Option 3: Implement macOS Support (Advanced)

To add native macOS support, you would need to:

1. Create platform-specific implementation files:
   - `networkinterface_darwin.go` - Implement macOS-specific network interface handling
   - `tc_program_darwin.go` - Alternative approach for TCP RST handling on macOS

2. Modify core components to be platform-aware:
   - Use platform-specific compilation tags
   - Implement BPF packet filtering using libpcap instead of eBPF on macOS
   - Use macOS-compatible raw socket alternatives

## Implementation Steps (for Option 3)

1. Create a new file `networkinterface_darwin.go` with macOS-specific socket implementation
2. Create a new file `tc_program_darwin.go` with alternative TCP RST handling
3. Add build constraints to Linux-specific files
4. Modify packet capture implementation to use libpcap on macOS
5. Update Makefile and build scripts to support macOS

## Recommended Approach

For immediate use, I recommend Option 1 (Docker) or Option 2 (Linux VM) since the code would require significant refactoring to work natively on macOS.

If you're interested in contributing macOS support to the project, Option 3 provides a starting point, but would require deep knowledge of macOS networking APIs and Go's platform-specific compilation.
