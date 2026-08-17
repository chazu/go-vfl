# Feature Specification: VFL Parser Library with UI Integration

**Feature Branch**: `001-create-a-golang`
**Created**: 2025-09-27
**Status**: Draft
**Input**: User description: "create a golang library for parsing vfl and extended vfl, and an example library showing how to use it with modernc.org/tk to generate UIs - see plan.md for details"

## Execution Flow (main)
```
1. Parse user description from Input
   → If empty: ERROR "No feature description provided"
2. Extract key concepts from description
   → Identify: actors, actions, data, constraints
3. For each unclear aspect:
   → Mark with [NEEDS CLARIFICATION: specific question]
4. Fill User Scenarios & Testing section
   → If no clear user flow: ERROR "Cannot determine user scenarios"
5. Generate Functional Requirements
   → Each requirement must be testable
   → Mark ambiguous requirements
6. Identify Key Entities (if data involved)
7. Run Review Checklist
   → If any [NEEDS CLARIFICATION]: WARN "Spec has uncertainties"
   → If implementation details found: ERROR "Remove tech details"
8. Return: SUCCESS (spec ready for planning)
```

---

## ⚡ Quick Guidelines
- ✅ Focus on WHAT users need and WHY
- ❌ Avoid HOW to implement (no tech stack, APIs, code structure)
- 👥 Written for business stakeholders, not developers

### Section Requirements
- **Mandatory sections**: Must be completed for every feature
- **Optional sections**: Include only when relevant to the feature
- When a section doesn't apply, remove it entirely (don't leave as "N/A")

### For AI Generation
When creating this spec from a user prompt:
1. **Mark all ambiguities**: Use [NEEDS CLARIFICATION: specific question] for any assumption you'd need to make
2. **Don't guess**: If the prompt doesn't specify something (e.g., "login system" without auth method), mark it
3. **Think like a tester**: Every vague requirement should fail the "testable and unambiguous" checklist item
4. **Common underspecified areas**:
   - User types and permissions
   - Data retention/deletion policies
   - Performance targets and scale
   - Error handling behaviors
   - Integration requirements
   - Security/compliance needs

---

## Clarifications

### Session 2025-09-27
- Q: How many example layouts should the integration library demonstrate? → A: 10 examples (comprehensive showcase)
- Q: What performance target should the parser achieve for constraint parsing? → A: No specific target (best effort optimization)
- Q: What is the maximum VFL program complexity the parser should handle? → A: Unlimited (memory-bound only)
- Q: How should the parser handle invalid VFL syntax errors? → A: Return partial AST with error details
- Q: How should circular references in named views be handled? → A: Detect and return error immediately

## User Scenarios & Testing *(mandatory)*

### Primary User Story
As a developer building user interfaces, I need to parse Visual Format Language (VFL) strings and Extended VFL strings to generate layout constraints that can be applied to UI components. I want to define layouts using a readable string format rather than manual constraint programming, and see these layouts rendered in a UI framework.

### Acceptance Scenarios
1. **Given** a valid VFL string like "H:|[button(100)]-[textField]-|", **When** the parser processes this string, **Then** it produces an AST representing horizontal layout constraints with a 100-point button and flexible text field
2. **Given** an Extended VFL string with nested views and connections, **When** the parser processes this string, **Then** it produces an AST with proper view hierarchy and constraint relationships
3. **Given** multiple VFL programs with named views, **When** these are composed together, **Then** the parser links named views correctly to create a unified layout program
4. **Given** a parsed VFL AST, **When** integrated with a UI library, **Then** the UI renders according to the specified constraints

### Edge Cases
- When invalid VFL syntax is provided, parser returns partial AST with detailed error information
- Circular references in named views are detected during composition and return an error immediately
- When conflicting constraints are specified, the parser returns a CompositionError with details about the conflicting constraints, allowing the developer to resolve them manually
- Malformed Extended VFL syntax results in partial AST with specific error locations
- When referenced named views don't exist, the composer returns a MissingNamedView error immediately during composition phase (consistent with FR-003)

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: Library MUST parse standard Apple Visual Format Language syntax into an Abstract Syntax Tree (AST)
- **FR-002**: Library MUST parse Extended Visual Format Language (EVFL) syntax including:
  - Standard VFL operators (|, -, [], (), ==, >=, <=, @priority)
  - Percentage-based constraints (e.g., [view(==50%)])
  - Mathematical expressions (e.g., [view(==otherView.width*2+10)])
  - View stacks notation (e.g., [column{header,content,footer}])
  - Attribute references (e.g., view.width, superview.centerX)
  - Equal spacers (~) and disconnections (→)
- **FR-003**: Library MUST support composition of multiple VFL programs through named view references and detect circular references with immediate error reporting
- **FR-004**: Library MUST validate VFL syntax and return partial AST with error messages containing:
  - Line and column position of error
  - Error type (syntax, semantic, reference)
  - Expected tokens at error location
  - Actual token found
  - Contextual snippet with error position marked
- **FR-005**: Library MUST expose parsed AST in a format consumable by UI libraries
- **FR-006**: Example integration MUST demonstrate rendering parsed VFL constraints in a UI framework
- **FR-007**: Library MUST handle view metrics/constants with standard values (standard=8, medium=20, default=8) matching Apple VFL conventions
- **FR-008**: Parser MUST be reusable across multiple parsing operations without state conflicts
- **FR-009**: Library MUST provide comprehensive test coverage for all VFL syntax variations
- **FR-010**: Example integration MUST show at least 10 example layouts demonstrating comprehensive VFL capabilities
- **FR-011**: Library performance MUST be optimized for typical UI layout needs without specific throughput targets
- **FR-012**: Library MUST handle VFL programs of unlimited complexity (bounded only by available memory)

### Key Entities *(include if feature involves data)*
- **VFL Program**: A string representation of layout constraints using Visual Format Language syntax
- **AST Node**: A node in the parsed abstract syntax tree representing layout elements, constraints, or relationships
- **Named View**: A reusable view definition that can be referenced by name in other VFL programs
- **Constraint**: A layout rule specifying size, position, or relationship between UI elements
- **View Metrics**: Predefined constant values used in VFL programs for spacing and sizing
- **Parser Instance**: The stateful or stateless parser that processes VFL strings
- **Layout Result**: The final composed AST ready for UI rendering

---

## Review & Acceptance Checklist
*GATE: Automated checks run during main() execution*

### Content Quality
- [x] No implementation details (languages, frameworks, APIs)
- [x] Focused on user value and business needs
- [x] Written for non-technical stakeholders
- [x] All mandatory sections completed

### Requirement Completeness
- [x] No [NEEDS CLARIFICATION] markers remain
- [x] Requirements are testable and unambiguous
- [x] Success criteria are measurable
- [x] Scope is clearly bounded
- [x] Dependencies and assumptions identified

---

## Execution Status
*Updated by main() during processing*

- [x] User description parsed
- [x] Key concepts extracted
- [x] Ambiguities marked
- [x] User scenarios defined
- [x] Requirements generated
- [x] Entities identified
- [x] Review checklist passed

---