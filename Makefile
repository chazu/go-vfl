.PHONY: help test bench lint coverage build-examples clean install

# Default target
help:
	@echo "Available targets:"
	@echo "  test           - Run all tests"
	@echo "  bench          - Run benchmark tests"
	@echo "  lint           - Run golangci-lint"
	@echo "  coverage       - Generate test coverage report"
	@echo "  build-examples - Build example applications"
	@echo "  clean          - Remove build artifacts"
	@echo "  install        - Install the library"

# Run all tests
test:
	go test -v -race ./...

# Run benchmark tests
bench:
	go test -bench=. -benchmem ./tests/benchmarks/...

# Run linter
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Generate test coverage (Constitution IV requirement: 80% minimum)
COVERAGE_THRESHOLD := 80.0
coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$coverage%"; \
	if (( $$(echo "$$coverage < $(COVERAGE_THRESHOLD)" | bc -l) )); then \
		echo "❌ ERROR: Coverage $$coverage% is below required threshold of $(COVERAGE_THRESHOLD)%"; \
		echo "Constitution IV requires minimum 80% test coverage"; \
		exit 1; \
	else \
		echo "✅ Coverage $$coverage% meets Constitution requirement (>= $(COVERAGE_THRESHOLD)%)"; \
	fi

# CI coverage gate - fails build if below 80%
coverage-gate: coverage
	@echo "Coverage gate passed!"

# Build example applications
build-examples:
	go build -o bin/vfl-tk-example ./examples/tk/

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean -cache

# Install the library
install:
	go install ./...

# Run tests with verbose output and race detection
test-verbose:
	go test -v -race -count=1 ./...

# Run only integration tests
test-integration:
	go test -v -race ./tests/integration/...

# Run only unit tests
test-unit:
	go test -v -race ./vfl/... ./evfl/...

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Check for vulnerabilities
vuln:
	@which govulncheck > /dev/null || go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

# Update dependencies
deps:
	go mod tidy
	go mod download
	go mod verify

# Quick check before committing
pre-commit: fmt lint test

# Full CI check
ci: deps lint test coverage build-examples