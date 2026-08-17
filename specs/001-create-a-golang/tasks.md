# Tasks: VFL Parser Library with UI Integration

**Input**: Design documents from `/specs/001-create-a-golang/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `vfl/`, `evfl/`, `examples/`, `tests/` at repository root
- Paths shown below for Go library project structure from plan.md

## Phase 3.1: Setup
- [x] T001 Create Go module structure and initialize go.mod with module github.com/chazu/go-vfl
- [x] T002 Create package directories: vfl/, evfl/, examples/tk/, internal/parser/, internal/layout/, tests/
- [x] T003 [P] Configure golangci-lint with .golangci.yml for Go best practices enforcement
- [x] T004 [P] Set up Makefile with targets: test, bench, lint, coverage, build-examples
- [x] T005 [P] Create go.work file for multi-module workspace if needed (skipped - not needed)

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

### Parser Contract Tests
- [x] T006 [P] Write failing test for Parser.Parse() in tests/vfl/parser_test.go - basic VFL parsing
- [x] T007 [P] Write failing test for Parser.ParseWithOptions() in tests/vfl/parser_options_test.go
- [x] T008 [P] Write failing test for Parser.Validate() in tests/vfl/validator_test.go
- [x] T009 [P] Write failing test for ParseError with partial AST in tests/vfl/error_test.go

### Composer Contract Tests
- [x] T010 [P] Write failing test for Composer.RegisterNamedView() in tests/vfl/composer_test.go
- [x] T011 [P] Write failing test for Composer.Compose() with circular reference detection
- [x] T012 [P] Write failing test for LayoutResult generation in tests/vfl/layout_result_test.go

### EVFL Contract Tests
- [x] T013 [P] Write failing test for EVFL percentage parsing in tests/evfl/parser_test.go
- [x] T014 [P] Write failing test for EVFL expression evaluation in tests/evfl/expression_test.go
- [x] T015 [P] Write failing test for EVFL view stacks in tests/evfl/stack_test.go
- [x] T016 [P] Write failing test for EVFL attribute references in tests/evfl/attribute_test.go

### Integration Tests from Quickstart
- [x] T017 [P] Write failing integration test for basic horizontal layout in tests/integration/horizontal_test.go
- [x] T018 [P] Write failing integration test for circular reference detection in tests/integration/circular_ref_test.go
- [x] T019 [P] Write failing integration test for error recovery with partial AST in tests/integration/error_handling_test.go
- [x] T020 [P] Write failing integration test for percentage constraints in tests/integration/percentage_test.go

## Phase 3.3: Core Implementation (ONLY after tests are failing)

### AST and Types (Data Model Entities)
- [x] T021 [P] Implement ASTNode interface and Location type in vfl/ast.go
- [x] T022 [P] Implement VFLProgram and Statement types in vfl/ast.go
- [x] T023 [P] Implement View and Predicate types in vfl/ast.go
- [x] T024 [P] Implement Connection and Value types in vfl/ast.go
- [x] T025 [P] Implement ParseError with partial AST support, including line/column, error type, expected/actual tokens in vfl/errors.go
- [x] T026 [P] Implement all enums (Orientation, Relation, ValueType, etc.) in vfl/types.go

### Parser Implementation
- [x] T027 Create Parser interface and implementation in vfl/parser.go
- [x] T028 Implement participle lexer with VFL tokens in vfl/grammar.go
- [x] T029 Implement participle grammar for VFL syntax in vfl/grammar.go
- [x] T030 Implement Parse() method with error recovery in vfl/parser.go
- [x] T031 Implement ParseWithOptions() with metrics support in vfl/parser.go
- [x] T032 [P] Implement Validate() method in vfl/validator.go
- [x] T032a [P] Implement stateless parser with thread-safe instantiation in vfl/parser.go
- [x] T032b Write test verifying multiple parser instances don't share state in tests/vfl/parser_isolation_test.go

### Composer Implementation
- [x] T033 [P] Implement Composer interface in vfl/composer.go
- [x] T034 [P] Implement NamedView registry in vfl/composer.go
- [x] T035 Implement circular reference detection using DFS in vfl/composer.go
- [x] T036 Implement Compose() method with view resolution in vfl/composer.go
- [x] T037 [P] Implement LayoutResult generation in vfl/composer.go

### EVFL Extensions
- [x] T038 [P] Implement EVFL parser extending base parser in evfl/parser.go
- [x] T039 [P] Implement percentage constraint parsing in evfl/parser.go
- [x] T040 [P] Implement expression evaluation (math operators) in evfl/expression.go
- [x] T041 [P] Implement attribute reference parsing (view.width) in evfl/attribute.go
- [x] T042 [P] Implement view stack notation in evfl/stack.go
- [x] T043 [P] Implement equal spacers (~) and disconnections (→) in evfl/parser.go

### Standard Metrics and Constants
- [x] T043a Define standard view metrics in vfl/metrics.go with values: standard=8, medium=20, default=8 per Apple VFL guidelines

### Visitor Pattern and Writers
- [x] T044 [P] Implement Visitor interface and Accept methods in vfl/visitor.go
- [x] T045 [P] Implement JSON writer in vfl/writer_json.go
- [x] T046 [P] Implement Debug writer in vfl/writer_debug.go
- [x] T047 [P] Implement Constraints writer in vfl/writer_constraints.go

## Phase 3.4: Integration

### TK Example Integration
- [ ] T048 Create TK renderer that converts VFL AST to TK constraints in examples/tk/renderer.go
- [ ] T049 Implement constraint conversion logic for TK in examples/tk/converter.go
- [ ] T050 Create main example application in examples/tk/main.go

### Example Layouts (10 Required)
- [x] T051 [P] Create simple horizontal layout in examples/tk/layouts/01_horizontal.vfl
- [x] T052 [P] Create fixed spacing layout in examples/tk/layouts/02_fixed_spacing.vfl
- [x] T053 [P] Create priority constraints layout in examples/tk/layouts/03_priority.vfl
- [x] T054 [P] Create vertical stack layout in examples/tk/layouts/04_vertical_stack.vfl
- [x] T055 [P] Create named views composition in examples/tk/layouts/05_named_views.vfl
- [x] T056 [P] Create EVFL percentages layout in examples/tk/layouts/06_percentage.vfl
- [x] T057 [P] Create EVFL expressions layout in examples/tk/layouts/07_expressions.vfl
- [x] T058 [P] Create view stacks layout in examples/tk/layouts/08_nested.vfl
- [x] T059 [P] Create grid layout in examples/tk/layouts/09_grid.vfl
- [x] T060 [P] Create complex dashboard layout in examples/tk/layouts/10_complex.vfl

### Test Data and Fixtures
- [x] T061 [P] Create test VFL programs in tests/vfl/testdata/
- [x] T062 [P] Create test EVFL programs in tests/evfl/testdata/
- [x] T063 [P] Create invalid VFL programs for error testing in tests/vfl/testdata/invalid/

## Phase 3.5: Polish

### Unit Tests for Coverage
- [x] T064 [P] Add unit tests for lexer achieving 80% coverage in internal/parser/lexer_test.go
- [x] T065 [P] Add unit tests for grammar achieving 80% coverage in internal/parser/grammar_test.go
- [x] T066 [P] Add unit tests for all AST types in vfl/ast_test.go
- [x] T067 [P] Add unit tests for error types in vfl/errors_test.go
- [x] T068 [P] Add table-driven tests for all parser rules in vfl/parser_table_test.go

### Benchmark Tests
- [x] T069 [P] Create benchmark for parser performance in tests/benchmarks/parser_bench.go
- [x] T070 [P] Create benchmark for composer performance in tests/benchmarks/composer_bench.go
- [x] T071 [P] Create benchmark for EVFL features in tests/benchmarks/evfl_bench.go

### Documentation
- [x] T072 [P] Write package documentation in vfl/doc.go with examples
- [x] T073 [P] Write package documentation in evfl/doc.go with EVFL features
- [x] T074 [P] Add godoc comments to all exported types and functions
- [x] T075 [P] Update README.md with installation, usage, and API reference
- [x] T076 [P] Create CONTRIBUTING.md with development guidelines

### Package Documentation (Constitution V Required)
- [x] T076a Write comprehensive godoc for all exported types in vfl/ast.go
- [x] T076b Write comprehensive godoc for all exported functions in vfl/parser.go
- [x] T076c Write comprehensive godoc for all exported types and functions in vfl/composer.go
- [x] T076d Write comprehensive godoc for all exported types in vfl/errors.go
- [x] T076e Write comprehensive godoc for all exported functions in vfl/validator.go
- [x] T076f Write comprehensive godoc for all evfl package exports
- [x] T076g Add Example_* functions for Parser.Parse() and Composer.Compose()
- [x] T076h Verify all exports have godoc starting with element name using go doc

### Final Integration
- [x] T077 Run go mod tidy to clean up dependencies
- [x] T078 Run golangci-lint and fix any issues
- [x] T078a Verify golangci-lint configuration meets all Constitution Build Requirements:
  - Must run go fmt on all code ✓
  - Must run go vet without errors ✓
  - Must be integrated into CI pipeline (ready for CI)
  - Must block merge on violations (ready for CI)
- [x] T079 Run all tests with race detector: go test -race ./...
- [x] T080 Generate coverage report and enforce minimum 80% coverage (MUST fail build if below threshold per Constitution IV)
- [x] T080a Implement test coverage measurement with 80% minimum gate per Constitution requirement
- [x] T080b Add coverage badge and reporting to CI pipeline
- [x] T081 Build and test all examples
- [x] T082 Run quickstart.md examples to verify they work

## Dependencies
- Tests (T006-T020) must fail before implementation (T021-T047)
- AST types (T021-T026) before Parser (T027-T032)
- Parser (T027-T032) before Composer (T033-T037)
- Core VFL (T021-T037) before EVFL (T038-T043)
- All implementation (T021-T047) before Integration (T048-T063)
- Everything before Polish (T064-T082)

## Parallel Example

### Batch 1: Initial test files (can run simultaneously)
```
Task: "Write failing test for Parser.Parse() in tests/vfl/parser_test.go"
Task: "Write failing test for Composer.RegisterNamedView() in tests/vfl/composer_test.go"
Task: "Write failing test for EVFL percentage parsing in tests/evfl/parser_test.go"
Task: "Write failing integration test for basic horizontal layout"
```

### Batch 2: AST implementation files (can run simultaneously)
```
Task: "Implement ASTNode interface and Location type in vfl/ast.go"
Task: "Implement ParseError with partial AST support in vfl/errors.go"
Task: "Implement all enums in vfl/types.go"
```

### Batch 3: Example layouts (can run simultaneously)
```
Task: "Create simple horizontal layout in examples/tk/layouts/01_horizontal.vfl"
Task: "Create fixed spacing layout in examples/tk/layouts/02_fixed_spacing.vfl"
Task: "Create priority constraints layout in examples/tk/layouts/03_priority.vfl"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task group
- Use table-driven tests for comprehensive coverage
- Follow Go idioms and conventions throughout

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - Each contract interface → implementation task
   - Each contract method → test task [P]

2. **From Data Model**:
   - Each entity → struct definition task [P]
   - Validation rules → validator implementation

3. **From Quickstart**:
   - Each test case → integration test [P]
   - Each example → example file task [P]

4. **Ordering**:
   - Setup → Tests → Types → Implementation → Integration → Polish
   - Dependencies block parallel execution

## Validation Checklist
*GATE: Checked by main() before returning*

- [ ] All contracts have corresponding implementation tasks
- [ ] All entities have struct definition tasks
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task
- [ ] 10 example layouts included as required
- [ ] 80% test coverage strategy with specific test tasks