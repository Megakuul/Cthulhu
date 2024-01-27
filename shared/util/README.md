# util

This directory includes go and cpp utility libraries that are independent of the application itself. It includes datatypes or functions that simplify a specific task, but they are absolutly independent of any other part of the cthulhu system or library.

Every function / datatype uses its own file and respective bazel rule, only dependency that is used is the standard library.
