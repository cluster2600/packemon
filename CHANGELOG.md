# Changelog

All notable changes to the Packemon project will be documented in this file.

## [Unreleased]

### eureka Added
- macOS native support with platform-specific implementations
  - New `networkinterface_darwin.go` implementation using libpcap instead of raw sockets
  - New `networkinterface_linux.go` with Linux-specific implementation moved from networkinterface.go
  - New `tc_program_darwin.go` for TCP RST packet handling on macOS using Packet Filter (PF)
  - New installation script `darwin_install.sh` for building on macOS
  - New Docker setup script `macos_docker_setup.sh` for running in Docker on macOS
  - Added detailed documentation in `README_MACOS.md` and `darwin_support.md`

### Changed
- Refactored networking code to use platform-specific implementations with build tags
- Updated README to mention macOS support

### Fixed
- None

## [Earlier Changes]
- Various features and improvements prior to implementing the changelog
