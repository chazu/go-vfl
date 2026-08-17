# Research: VFL Parser Library Implementation

**Date**: 2025-09-27
**Feature**: VFL Parser Library with UI Integration

## Technology Decisions

### 1. Parsing Framework
**Decision**: participle/v2 with extended lexer
**Rationale**: Already integrated, provides recursive descent parsing suitable for nested VFL structures
**Alternatives considered**:
- Hand-written parser: More control but significant development effort
- ANTLR: Powerful but adds Java dependency and complexity
- goyacc: Lower-level, less maintainable

### 2. AST Representation
**Decision**: Separate AST packages for VFL and Extended VFL with shared interfaces
**Rationale**: Clear separation of concerns while allowing polymorphic handling
**Alternatives considered**:
- Single unified AST: Would complicate standard VFL parsing
- Fully separate implementations: Code duplication

### 3. Error Handling Strategy
**Decision**: Custom error types with location tracking, partial AST generation
**Rationale**: Per requirements, must return partial AST with detailed error information
**Alternatives considered**:
- Fail-fast approach: Would not meet partial parsing requirement
- Warning-based system: Too permissive for syntax errors

### 4. Circular Reference Detection
**Decision**: Graph-based cycle detection during composition phase
**Rationale**: Efficient O(V+E) detection using DFS traversal
**Alternatives considered**:
- Runtime detection: Would allow invalid ASTs to be created
- Reference counting: More complex, doesn't handle all cases

### 5. Testing Approach
**Decision**: Table-driven tests with testdata directories
**Rationale**: Go idiom, easy to add test cases, clear test organization
**Alternatives considered**:
- Individual test functions: Less maintainable at scale
- Property-based testing: Overkill for deterministic parsing

## VFL Syntax Analysis

### Standard VFL Components
- **Orientation prefixes**: H: (horizontal), V: (vertical)
- **View notation**: [viewName(predicates)]
- **Superview edges**: | (pipe character)
- **Connections**: - (standard spacing), -N- (fixed spacing)
- **Relations**: ==, >=, <=
- **Priority**: @N (numeric priority)

### Extended VFL Additions
1. **Percentage sizing**: [view(==50%)]
2. **Mathematical expressions**: [view1(==view2/2-10)]
3. **Attribute references**: [view2(view1.width)]
4. **Z-ordering**: Z:|[child1][child2]
5. **Equal spacers**: |~[view]~|
6. **View stacks**: V:|[column:[views]]|
7. **View ranges**: H:[view1..5(size)]
8. **Multiple views**: [view1,view2,view3]
9. **Disconnections**: [left]→[right]
10. **Negative spacing**: [view1]-(-10)-[view2]

## Parsing Challenges & Solutions

### Challenge 1: Ambiguous Grammar
**Issue**: Minus sign as both connection and arithmetic operator
**Solution**: Context-aware lexing with different modes

### Challenge 2: Complex Expressions
**Issue**: Nested mathematical expressions with operator precedence
**Solution**: Recursive descent with explicit precedence levels

### Challenge 3: Lookahead Requirements
**Issue**: Some constructs require significant lookahead
**Solution**: Configure participle with appropriate lookahead value (minimum 5)

### Challenge 4: Error Recovery
**Issue**: Must continue parsing after errors to generate partial AST
**Solution**: Error nodes in AST, skip-to-recovery-point strategy

## modernc.org/tk Integration Research

### TK Constraint System
- Uses Cassowary constraint solver (same as iOS AutoLayout)
- Requires constraint equations in specific format
- Supports priority-based constraint resolution

### Integration Approach
1. Convert VFL AST to constraint equations
2. Map VFL priorities to TK constraint strengths
3. Handle view hierarchy through TK widget tree
4. Render using TK's layout engine

## Performance Considerations

### Parsing Performance
- No specific targets per requirements
- Focus on: Minimize allocations, reuse parser instances, efficient string operations

### Memory Management
- AST nodes should be garbage-collected normally
- Parser state should be minimal between operations
- Consider object pooling for frequently created types

## Best Practices Implementation

### Go Idioms
1. **Error wrapping**: Use fmt.Errorf with %w verb
2. **Interface design**: Small, focused interfaces (e.g., Node, Visitor)
3. **Package organization**: Internal details in internal/
4. **Documentation**: Every exported type/function with godoc

### Testing Strategy
1. **Unit tests**: Each package >80% coverage
2. **Table-driven tests**: All parser rules
3. **Integration tests**: Multi-program composition
4. **Benchmark tests**: Parser performance
5. **Fuzz testing**: Invalid input handling

## Dependencies Analysis

### Current Dependencies (keep)
- github.com/alecthomas/participle/v2: Parsing framework
- github.com/stretchr/testify: Testing assertions
- github.com/lithdew/casso: Constraint solver (for layout)

### New Dependencies (add)
- modernc.org/tk: UI toolkit for examples
- No additional runtime dependencies for core library

### Dependency Constraints
- Keep core library dependency-free except participle
- Example code can have additional UI dependencies
- Test dependencies are acceptable

## Implementation Phases

### Phase 1: Core VFL Parser
- Standard VFL syntax support
- Basic AST generation
- Error handling with location info

### Phase 2: Extended VFL Support
- Mathematical expressions
- Percentage sizing
- Attribute references

### Phase 3: Advanced Features
- View stacks and ranges
- Z-ordering
- Equal spacers

### Phase 4: Integration & Polish
- TK example application
- Performance optimization
- Comprehensive documentation

## Open Questions Resolved

All clarifications from spec have been addressed:
- ✅ Example count: 10 comprehensive examples
- ✅ Performance: Best-effort optimization
- ✅ Complexity limits: Memory-bound only
- ✅ Error handling: Partial AST with details
- ✅ Circular references: Immediate error detection