# GoType

GoType parses Golang's source code files and generates a Golang's type representation statically.


## Motivation

GoType is used to manipulate Golang's type without using Golang's reflection.
The problem with reflection is we can't get a Golang's type representation which identified by the package path and type's name.
