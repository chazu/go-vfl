
# Implementation Plan: VFL Parser Library with UI Integration

**Branch**: `001-create-a-golang` | **Date**: 2025-09-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-create-a-golang/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from file system structure or context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Fill the Constitution Check section based on the content of the constitution document.
4. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
5. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
6. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, `GEMINI.md` for Gemini CLI, `QWEN.md` for Qwen Code or `AGENTS.md` for opencode).
7. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
8. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
9. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Create a Go library that parses Visual Format Language (VFL) and Extended VFL strings into Abstract Syntax Trees, supporting composition of multiple programs with named views. The library will provide comprehensive error handling with partial AST generation, circular reference detection, and an example integration demonstrating UI rendering with modernc.org/tk.

## Technical Context
**Language/Version**: Go 1.21.4
**Primary Dependencies**: participle/v2 (parsing), modernc.org/tk (UI example), testify (testing)
**Storage**: N/A (library operates on in-memory strings/AST)
**Testing**: go test with testify/assert, table-driven tests, benchmark tests
**Target Platform**: Cross-platform Go library (Linux, macOS, Windows)
**Project Type**: single (library project with example)
**Performance Goals**: Best-effort optimization for typical UI layouts (no specific target per clarifications)
**Constraints**: Memory-bound only, no specific performance constraints
**Scale/Scope**: Unlimited program complexity (memory-bound), 10+ comprehensive examples

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Go Best Practices Gates**:
- [x] Idiomatic Go patterns followed (naming, errors, interfaces)
- [x] Package design single-purpose with minimal public API
- [x] Error handling explicit at all levels
- [x] Test coverage strategy defined (80% minimum)
- [x] Godoc documentation planned for all exports
- [x] Concurrency safety documented for shared state
- [x] Resource management with proper defer patterns

## Project Structure

### Documentation (this feature)
```
specs/001-create-a-golang/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
vfl/
├── parser.go           # Main parser interface and implementation
├── ast.go             # AST node definitions
├── errors.go          # Custom error types with location info
├── validator.go       # Syntax validation logic
├── composer.go        # Program composition and named view linking
└── doc.go            # Package documentation

evfl/
├── parser.go          # Extended VFL parser
├── ast.go            # Extended AST nodes
└── doc.go            # Package documentation

examples/
└── tk/
    ├── main.go        # Example UI application
    ├── layouts/       # 10+ example VFL programs
    └── renderer.go    # VFL to TK constraint conversion

internal/
├── parser/           # Existing parser implementation (refactored)
│   ├── lexer.go     # Tokenization
│   ├── grammar.go   # Participle grammar definitions
│   └── types.go     # Internal parsing types
└── layout/          # Existing layout implementation

tests/
├── vfl/
│   ├── parser_test.go          # Parser unit tests
│   ├── validator_test.go       # Validation tests
│   ├── composer_test.go        # Composition tests
│   └── testdata/              # Test VFL programs
├── evfl/
│   ├── parser_test.go         # Extended VFL tests
│   └── testdata/             # Test EVFL programs
├── integration/
│   ├── circular_ref_test.go   # Circular reference detection
│   ├── error_handling_test.go # Error reporting tests
│   └── composition_test.go    # Multi-program composition
└── benchmarks/
    └── parser_bench.go        # Performance benchmarks
```

**Structure Decision**: Single Go library project with separate packages for VFL and Extended VFL parsing, internal implementation details hidden, and comprehensive examples using modernc.org/tk. The existing internal packages will be refactored and reused where appropriate.

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `.specify/scripts/bash/update-agent-context.sh claude`
     **IMPORTANT**: Execute it exactly as specified above. Do not add or remove any arguments.
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `.specify/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each API contract file → implementation and test tasks
- Each entity in data model → struct definition task
- Each test scenario from quickstart → test implementation task
- Example implementations for all 10 required layouts

**Specific Task Categories**:
1. **Setup Tasks** (3 tasks):
   - Initialize Go module structure
   - Set up package directories
   - Configure linting and testing

2. **Core VFL Parser** (8-10 tasks):
   - Parser interface and types
   - AST node definitions
   - Lexer implementation
   - Grammar rules
   - Error handling with partial AST
   - Validation logic

3. **Extended VFL Support** (6-8 tasks):
   - Percentage parsing
   - Expression evaluation
   - Attribute references
   - Special operators (~ and →)
   - View stacks and ranges

4. **Composition & Named Views** (4-5 tasks):
   - Composer implementation
   - Named view registry
   - Circular reference detection
   - View hierarchy builder

5. **Testing** (8-10 tasks):
   - Unit tests for each component
   - Integration tests for composition
   - Error handling tests
   - Benchmark tests
   - Test data fixtures

6. **TK Integration Example** (5-6 tasks):
   - TK renderer implementation
   - Constraint conversion
   - 10 example layouts
   - Example application

**Ordering Strategy**:
- TDD order: Tests before implementation
- Bottom-up: AST types → Parser → Composer → Examples
- Parallel execution for independent packages
- Sequential for dependent components

**Estimated Output**: 35-40 numbered tasks with clear dependencies

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (none required)

---
*Based on Constitution v1.0.0 - See `.specify/memory/constitution.md`*
