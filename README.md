Go-VFL
------

A hopelessly naive attempt to implement Visual Format Language parsing (and rendering) in golang.

All of the layout code has been lifted wholesale from lithdew/blanc - in the interests of figuring out how to use cassowary at a lower level by studying his higher-level usage.


TODO
- Figure out how to plug a reified program AST into lithdew's layout API, hoping to understand how it works
- Figure out how to reify a layout comprising multiple programs using lithdew's Layout API - or at least what changes are needed to do so



DONE
Add golang cassowary implementation (https://github.com/lithdew/casso)
reify shitty parsing result into something less awful for the public API
Create struct(s) to hold program ast reified from horrible private parser structs
Make current parser structs private
Move parser back into internal you dimwit
Move parser out of internal you dimwit
Move parser into module
Wrap participle in parser type with functional options pattern
 - Lookahead length
Move test cases to actual test file

WONTFIX
IN PROGRESS - figure out how to use either Capture interface or Parseable to, e.g. parse connections into something not totally idiotic. Probably Parseable tbh
