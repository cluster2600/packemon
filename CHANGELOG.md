# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- eureka: Added macOS support with platform-specific implementations
- Added platform-specific network interface handling for macOS using libpcap
- Added platform-specific TCP RST handling for macOS using Packet Filter (PF)
- Added `networkinterface_darwin.go` for macOS-specific network operations
- Added `tc_program_darwin.go` for macOS-specific packet filtering
- Added `darwin_install.sh` script for native macOS installation
- Added `macos_docker_setup.sh` script for Docker-based installation on macOS
- Added `README_MACOS.md` with macOS-specific instructions
- Added `darwin_support.md` with technical implementation details
- Added PR template for macOS contributions

### Changed
- Refactored network interface code to use platform-specific implementations with build tags
- Refactored TCP program code to use platform-specific implementations with build tags
- Updated `passive.go` with platform-agnostic packet parsing
- Improved error handling for platform-specific operations
- Enhanced packet parsing with better type safety and error checking

### Fixed
- Fixed compatibility issues between Linux and macOS implementations
- Fixed memory management in packet processing for cross-platform support
- Fixed build issues for macOS targets

## [1.0.0] - 2025-01-15

### Added
- Initial release with Linux support
- Packet capture and analysis capabilities
- Protocol support for Ethernet, ARP, IPv4, IPv6, ICMP, ICMPv6, TCP, UDP, DNS, HTTP, TLS
- Terminal-based UI for packet monitoring and generation
- Support for sending custom packets
- Support for filtering packets
