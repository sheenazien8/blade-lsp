# AGENTS.md - Development Guide for Agentic Coding Agents

## Build/Test Commands
- **Build**: `go build .` (builds main executable)
- **Test**: `go test ./...` (runs all tests)
- **Single test**: `go test ./rpc` (test specific package)
- **Run**: `go run .` (runs the LSP server)

## Project Structure
- `main.go` - Entry point and message handling
- `lsp/` - LSP protocol types and responses
- `rpc/` - RPC message encoding/decoding
- `analysis/` - Document state management and analysis

## Code Style Guidelines
- **Imports**: Standard library first, then local packages (e.g., `github.com/sheenazien8/blade-lsp/lsp`)
- **Naming**: PascalCase for exported types/functions, camelCase for unexported
- **Types**: Use struct tags for JSON serialization (`json:"fieldName"`)
- **Error handling**: Always check and log errors with descriptive messages
- **Comments**: Minimal comments, focus on why not what
- **Formatting**: Use `go fmt` (standard Go formatting)

## Key Patterns
- Message handling via switch statements on method names
- State management through `analysis.State` struct
- JSON unmarshaling for LSP requests/responses
- Logging to file for debugging (`getLogger` function)