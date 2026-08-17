# Go-VFL

A parser for Apple's Visual Format Language and Extended Visual Format Language. The parser takes strings which are visual format language "programs" and parses them into ASTs which can then be used in combination with UI libraries to produce layouts for user interfaces.

Programs can contain named views. These named views may in turn be programs. A parser instance can accept multiple programs and compose them into a single program by linking the named views together. Of course this requires that some programs be registered with the parser with names, so that other programs which contain them can include them.

The parser should include a full suite of tests as well as a library demonstrating how to take the output from the parser and use it to construct UIs using the modernc.org/tk package.


# Sources
Apple's original documentation on Visual Format Language can be found here:
https://developer.apple.com/library/archive/documentation/UserExperience/Conceptual/AutolayoutPG/VisualFormatLanguage.html

Info about Extended Visual Format Language can be found here: https://github.com/lume/autolayout?tab=readme-ov-file#extended-visual-format-language-evfl
