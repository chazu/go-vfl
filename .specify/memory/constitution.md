<!-- Sync Impact Report
Version Change: 0.0.0 → 1.0.0 (initial adoption)
Added Sections:
- Core Principles (7 principles for Go library development)
- Development Standards
- Quality Requirements
- Governance
Templates Requiring Updates: ⚠ pending review
- .specify/templates/plan-template.md
- .specify/templates/spec-template.md
- .specify/templates/tasks-template.md
Follow-up TODOs:
- RATIFICATION_DATE: Set to today (2025-09-27) as initial adoption
-->

# Go-VFL Constitution

## Core Principles

### I. Idiomatic Go First
All code MUST follow Go idioms and conventions as defined in Effective Go and the Go Code Review Comments. This includes proper naming conventions (exported vs unexported), interface design, error handling patterns, and package organization. Every deviation from idiomatic Go requires explicit justification in code comments.

### II. Package Design and Organization
Packages MUST have a single, well-defined purpose. Internal packages MUST be used for implementation details not part of the public API. The public API surface MUST be minimal and carefully designed. Package dependencies MUST form an acyclic graph with clear layering.

### III. Error Handling Excellence
Errors MUST be handled explicitly at every level. Custom error types MUST implement the error interface and support error wrapping with errors.Is/As. Sentinel errors MUST be defined as package-level variables. Panics are only acceptable for truly unrecoverable situations during initialization.

### IV. Test Coverage and Quality
Unit tests MUST achieve minimum 80% coverage for all packages. Table-driven tests are required for functions with multiple test cases. Integration tests MUST verify package boundaries and interactions. Benchmark tests are required for performance-critical code paths. The testing package and testify/assert are the only approved testing dependencies.

### V. Documentation and Examples
All exported types, functions, and packages MUST have godoc comments starting with the element name. Example functions (Example_* format) MUST be provided for non-trivial APIs. Package documentation MUST include a clear overview of purpose and usage. README MUST include installation, quick start, and API reference sections.

### VI. Concurrency Safety
All exported types MUST document their concurrency safety guarantees. Shared state MUST be protected by appropriate synchronization primitives (mutex, channels, atomic operations). Race conditions MUST be tested with go test -race. Goroutine lifecycles MUST be explicitly managed with proper cleanup.

### VII. Performance and Resource Management
Resources (files, connections, goroutines) MUST be explicitly closed/terminated. The defer pattern MUST be used for cleanup operations. Memory allocations in hot paths MUST be minimized and benchmarked. Context MUST be used for cancellation and timeout propagation.

## Development Standards

### Module and Dependency Management
- Go modules (go.mod) MUST be used for dependency management
- Module path MUST match the repository URL
- Direct dependencies MUST be minimal and justified
- go mod tidy MUST be run before every commit
- Vendoring is discouraged unless required for reproducible builds

### Build and CI Requirements
- Code MUST compile without warnings on latest stable Go version
- go fmt MUST be run on all code (enforced by CI)
- go vet MUST pass without errors
- golangci-lint or similar linter MUST be configured and pass
- CI MUST run on pull requests with all tests passing

### API Stability and Versioning
- Semantic versioning (v1.2.3 format) MUST be used for releases
- v1+ APIs MUST maintain backward compatibility
- Breaking changes require major version bump with migration guide
- Deprecated features MUST be marked with // Deprecated: comments
- Experimental features MUST be clearly marked in documentation

## Quality Requirements

### Code Review Standards
- All code changes require peer review before merging
- Reviews MUST verify adherence to this constitution
- Reviews MUST check for idiomatic Go patterns
- Security implications MUST be considered for all changes
- Performance impact MUST be assessed for critical paths

### Security Considerations
- Input validation MUST be performed at API boundaries
- Cryptographic operations MUST use standard library or vetted packages
- Secrets MUST never be logged or included in error messages
- SQL queries MUST use parameterized statements
- User input MUST be sanitized before use in system operations

## Governance

### Constitution Authority
This constitution supersedes all other coding practices and standards within the project. All contributors MUST read and understand these principles before submitting code. Violations of constitution principles MUST be addressed before code can be merged.

### Amendment Process
- Proposed amendments MUST be documented in a pull request
- Amendments require consensus from core maintainers
- Major principle changes require migration plan for existing code
- Version bump follows semantic versioning based on impact
- All dependent artifacts MUST be updated to reflect amendments

### Compliance Verification
- Pull request templates MUST include constitution compliance checklist
- Automated checks SHOULD verify as many principles as possible
- Regular audits of codebase for constitution adherence
- New contributors MUST acknowledge understanding of constitution

**Version**: 1.0.0 | **Ratified**: 2025-09-27 | **Last Amended**: 2025-09-27