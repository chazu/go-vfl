# Data Model: VFL Parser Library

**Date**: 2025-09-27
**Feature**: VFL Parser Library with UI Integration

## Core Entities

### 1. VFLProgram
**Purpose**: Represents a complete VFL string to be parsed
**Attributes**:
- `input`: string - The raw VFL string
- `orientation`: enum(Horizontal, Vertical, Z, Combined) - Layout direction
- `statements`: []Statement - Parsed constraint statements

**Relationships**:
- Contains multiple Statements
- Can reference NamedViews

**State Transitions**:
- Created → Parsing → Parsed/Failed
- Parsed → Composing (if has named views)
- Composing → Composed/Failed

**Validation Rules**:
- Must have valid orientation prefix or default to Horizontal
- Must contain at least one view or constraint

### 2. ASTNode
**Purpose**: Base interface for all AST elements
**Attributes**:
- `location`: Location - Position in source string
- `nodeType`: string - Type identifier for visitor pattern

**Relationships**:
- Parent of all specific node types
- Can contain child nodes

**Validation Rules**:
- Must have valid location information
- Must implement visitor accept method

### 3. Statement
**Purpose**: Single constraint statement within a program
**Attributes**:
- `views`: []View - Views involved in constraint
- `connections`: []Connection - Spacing between views
- `superviewStart`: bool - Starts at superview edge
- `superviewEnd`: bool - Ends at superview edge

**Relationships**:
- Contains Views and Connections
- Part of VFLProgram

**Validation Rules**:
- Views and connections must alternate
- Cannot have adjacent connections

### 4. View
**Purpose**: Represents a UI view with constraints
**Attributes**:
- `name`: string - View identifier
- `predicates`: []Predicate - Size/position constraints
- `priority`: int - Constraint priority (0-1000)
- `isNamedView`: bool - References external view

**Relationships**:
- Contains Predicates
- Can reference other Views

**State Transitions**:
- Unresolved → Resolved (when named view found)
- Resolved → Bound (when constraints applied)

**Validation Rules**:
- Name must be valid identifier
- Priority must be 0-1000
- Named views must exist in registry

### 5. Predicate
**Purpose**: Constraint expression for a view
**Attributes**:
- `relation`: enum(Equal, GreaterEqual, LessEqual) - Constraint relation
- `value`: Value - Constraint value
- `attribute`: string - For EVFL attribute references

**Relationships**:
- Belongs to View
- Can reference other Views or Metrics

**Validation Rules**:
- Value must be non-negative unless explicit negative spacing
- Attribute must be valid view property

### 6. Connection
**Purpose**: Spacing between views
**Attributes**:
- `spacing`: float - Space amount
- `isDefault`: bool - Uses default spacing
- `isEqualSpace`: bool - For EVFL equal spacers (~)
- `isDisconnection`: bool - For EVFL disconnections (→)

**Relationships**:
- Connects adjacent Views
- Part of Statement

**Validation Rules**:
- Spacing can be negative for overlapping
- Cannot have both isDefault and specific spacing

### 7. Value
**Purpose**: Represents constraint values
**Attributes**:
- `type`: enum(Constant, ViewReference, Metric, Expression, Percentage)
- `constant`: float - For numeric values
- `viewRef`: string - For view references
- `metricName`: string - For named metrics
- `expression`: Expression - For EVFL math expressions
- `percentage`: float - For EVFL percentages

**Relationships**:
- Used by Predicates
- Can reference Views or Metrics

**Validation Rules**:
- Exactly one value type must be set
- Percentages must be 0-100
- View references must exist

### 8. Expression (Extended VFL)
**Purpose**: Mathematical expression in constraints
**Attributes**:
- `operator`: enum(Add, Subtract, Multiply, Divide)
- `left`: Value - Left operand
- `right`: Value - Right operand

**Relationships**:
- Contains nested Values
- Used in EVFL predicates

**Validation Rules**:
- Division by zero check
- Type compatibility for operations

### 9. ParseError
**Purpose**: Detailed error information
**Attributes**:
- `message`: string - Error description
- `location`: Location - Position in source
- `errorType`: enum(Syntax, Semantic, Reference)
- `partialAST`: ASTNode - Partial parse result

**Relationships**:
- Associated with failed parse operations
- Contains partial AST if available

**Validation Rules**:
- Must have descriptive message
- Must include location if known

### 10. Location
**Purpose**: Source position tracking
**Attributes**:
- `line`: int - Line number (1-based)
- `column`: int - Column number (1-based)
- `offset`: int - Byte offset in source

**Relationships**:
- Used by all AST nodes and errors

**Validation Rules**:
- All values must be non-negative
- Line and column start at 1

### 11. NamedView
**Purpose**: Reusable view definition
**Attributes**:
- `name`: string - Global identifier
- `program`: VFLProgram - View definition
- `isComposite`: bool - Contains other named views

**Relationships**:
- Can contain nested NamedViews
- Referenced by Views

**State Transitions**:
- Registered → Resolved → Composed

**Validation Rules**:
- Name must be globally unique
- No circular references allowed

### 12. Parser
**Purpose**: Stateless parser instance
**Attributes**:
- `lookahead`: int - Parser lookahead setting
- `options`: ParserOptions - Configuration

**Relationships**:
- Creates VFLPrograms from strings
- Uses Lexer and Grammar

**Validation Rules**:
- Must be reusable across operations
- Must be concurrent-safe

### 13. Composer
**Purpose**: Combines multiple programs with named views
**Attributes**:
- `registry`: map[string]NamedView - Named view registry
- `visitedNodes`: set[string] - For cycle detection

**Relationships**:
- Manages NamedViews
- Produces composed VFLPrograms

**Validation Rules**:
- Must detect circular references
- Must validate all references exist

### 14. Metrics
**Purpose**: Named spacing constants
**Attributes**:
- `values`: map[string]float - Metric definitions
- `defaults`: MetricDefaults - Standard spacings

**Relationships**:
- Referenced by Values
- Used during constraint resolution

**Validation Rules**:
- Metric names must be valid identifiers
- Values must be non-negative

### 15. LayoutResult
**Purpose**: Final composed AST ready for rendering
**Attributes**:
- `constraints`: []Constraint - All resolved constraints
- `viewHierarchy`: ViewTree - View relationships
- `errors`: []ParseError - Any partial errors

**Relationships**:
- Output of Composer
- Input to UI renderer

**State Transitions**:
- Composing → Complete/Failed
- Complete → Rendering

**Validation Rules**:
- All named views must be resolved
- No conflicting constraints

## Entity State Machine

```
Input String
    ↓
[Parsing] → ParseError (with partial AST)
    ↓
VFLProgram (may have named view refs)
    ↓
[Composition] → CircularRefError
    ↓
LayoutResult
    ↓
[Rendering]
    ↓
UI Display
```

## Validation Rules Summary

1. **Syntax Validation**:
   - Valid VFL/EVFL syntax
   - Proper view/connection alternation
   - Valid operators and relations

2. **Semantic Validation**:
   - View references exist
   - No circular dependencies
   - Valid attribute names
   - Priority ranges (0-1000)

3. **Composition Validation**:
   - All named views registered
   - No reference cycles
   - Compatible constraint types

4. **Runtime Validation**:
   - No conflicting constraints
   - Valid metric references
   - Memory bounds checking