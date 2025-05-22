# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- [eureka] Implemented full ICMPv6 support for both generator and monitor components
- [eureka] Created new CONTRIBUTING.md with comprehensive guidelines for contributors
- [eureka] Added comprehensive testing infrastructure with unit and integration tests
- [eureka] Implemented BGP and OSPF protocol support with RFC-compliant implementations
- [eureka] Added buffer pooling system to reduce GC pressure and improve performance
- [eureka] Created configuration file support for saving and loading packet templates
- [eureka] Added packet statistics dashboard with real-time visualization
- [eureka] Set up CI/CD pipeline with GitHub Actions for automated builds and testing
- [eureka] Created Docker container for easy deployment
- Added proper protocol documentation with RFC references in ICMPv6 implementation
- [eureka] Successfully fixed Fletcher checksum calculation in OSPF implementation
- [eureka] Successfully implemented and tested BGP protocol with all message types
- [eureka] Successfully implemented and tested buffer pooling system for improved performance

### Enhanced 
- Extended the TUI interface with ICMPv6 packet generation capability
- Improved protocol parsing with better error handling and validation
- Added unit tests for ICMPv6 protocol implementation
- Created integration tests for packet generation and monitoring
- Optimized packet parsing with pre-allocation and buffer reuse
- Added keyboard shortcut documentation within the application
- Fixed OSPF Fletcher checksum calculation to comply with RFC 1008
- Improved BGP protocol implementation with proper message parsing and serialization
- Enhanced buffer pooling system with small, medium, and large packet size support

## [1.0.0] - 2025-05-22

### Added
- Initial creation of CHANGELOG.md. [eureka] Successfully read and analyzed README.md for project overview.
- [eureka] Successfully read internal/tui/tui.go, internal/tui/generator/generator.go, internal/tui/monitor/monitor.go, and internal/debugging/networkinterface.go for code analysis.
