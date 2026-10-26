# List all tasks
default:
    @just --list

# Format .templ files
templ-fmt:
    go tool templ fmt .

# Generate Go code from .templ files
generate:
    go tool templ generate

# Format Go files (gofumpt and goimports)
fmt:
    golangci-lint fmt

# Run the linters
lint:
    golangci-lint run

# Check for known vulnerabilities
vuln:
    go tool govulncheck ./...

# Run the tests
test:
    go test ./...

# Run the server at http://localhost:8080
[continue]
run:
    # go run ./cmd/server
    go build -o bin/tiebreaker-go ./cmd/server
    ./bin/tiebreaker-go

# Run go tidy to clean up dependencies
tidy:
    go mod tidy

# Format, generate, lint, scan, and test before committing
pre-commit: tidy templ-fmt generate fmt lint vuln test
