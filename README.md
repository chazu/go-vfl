TODO

reify shitty parsing result into something less awful for the public API

Add golang cassowary implementation (https://github.com/lithdew/casso)
Design API for accessing parsed program attributes

DONE
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
