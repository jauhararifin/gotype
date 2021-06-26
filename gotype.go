// gotype package provides functions to parse Golang's source code files and generates a Golang's type representation
// statically.
package gotype

var defaultAstTypeGenerator = &astTypeGenerator{
	sourceFinder: &defaultSourceFinder{},
}

// GenerateTypesFromSpecs find and parses Golang's source code to generate the `Type`s specified by the `typeSpecs`.
func GenerateTypesFromSpecs(typeSpecs ...TypeSpec) ([]Type, error) {
	return defaultAstTypeGenerator.GenerateTypesFromSpecs(typeSpecs...)
}
