# Contributing to Packemon

Thank you for your interest in contributing to Packemon! This document provides guidelines and instructions for contributing to the project.

## Project Structure

```
packemon/
├── assets/                 # Images and documentation assets
├── cmd/                    # Application entry points
│   ├── packemon/           # Main TUI application
│   ├── packemon-api/       # Web API server
│   └── debugging/          # Debugging utilities
├── internal/               # Internal packages
│   ├── tui/                # TUI components
│   │   ├── generator/      # Packet generator UI
│   │   ├── monitor/        # Packet monitor UI
│   └── debugging/          # Debugging utilities
├── tc_program/             # eBPF program for TCP RST packets
└── (root)                  # Core protocol implementations
```

## Development Environment Setup

1. Ensure you have Go 1.24+ installed
2. Clone the repository: `git clone https://github.com/ddddddO/packemon.git`
3. Install dependencies:
   ```
   cd packemon
   go mod download
   ```
4. Build the eBPF components:
   ```
   cd tc_program/
   go generate
   cd -
   ```
5. Build the application:
   ```
   go build -o packemon cmd/packemon/*.go
   ```

## Development Workflow

1. Create a new branch for your feature/fix
2. Make your changes
3. Write/update tests as needed
4. Ensure all tests pass: `go test ./...`
5. Submit a pull request

## Coding Standards

- Follow Go best practices and style guidelines
- Document all exported functions and types
- Include comments explaining complex protocol-specific logic
- Handle errors appropriately
- Add unit tests for new functionality

## Protocol Implementation Guidelines

When implementing or improving protocol support:

1. Create a new file named after the protocol (e.g., `icmpv6.go`)
2. Define protocol constants at the top of the file
3. Create struct(s) to represent the protocol header and data
4. Implement the following functions:
   - `New<Protocol>()` - Constructor
   - `Bytes()` - Serialization method
   - `Parsed<Protocol>()` - Parsing function
5. Update `networkinterface.go` to handle the new protocol
6. Add appropriate TUI components in `internal/tui/generator/` and `internal/tui/monitor/`

## Pull Request Process

1. Update the README.md with details of your changes if applicable
2. Update the CHANGELOG.md following the existing format
3. Your PR will be reviewed by a maintainer
4. Address any feedback or requested changes
5. Once approved, your PR will be merged

## Testing

- Write unit tests for protocol implementations
- Test packet generation and parsing
- Verify TUI functionality

## License

By contributing to Packemon, you agree that your contributions will be licensed under the project's license.
