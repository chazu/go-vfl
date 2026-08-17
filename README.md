# go-vfl

A Go implementation of Visual Format Language (VFL) parser and composer with Extended VFL (EVFL) support.

## Features

- ✅ **Complete VFL Parser**: Full support for Apple's Visual Format Language syntax
- ✅ **Extended VFL (EVFL)**: Percentage-based layouts, expressions, view stacks, and attribute references
- ✅ **Constraint Composition**: Combine multiple VFL programs with named view support
- ✅ **Multiple Output Formats**: Auto Layout, Cassowary, Kiwi constraint formats
- ✅ **Error Recovery**: Partial AST generation for better error handling
- ✅ **Circular Reference Detection**: Automatic detection of circular dependencies
- ✅ **Visitor Pattern**: Transform and analyze AST structures
- ✅ **Comprehensive Testing**: TDD approach with extensive test coverage
Move parser out of internal you dimwit
Move parser into module
Wrap participle in parser type with functional options pattern
 - Lookahead length
Move test cases to actual test file

WONTFIX
IN PROGRESS - figure out how to use either Capture interface or Parseable to, e.g. parse connections into something not totally idiotic. Probably Parseable tbh
