package gotype

var defaultAstTypeGenerator = &astTypeGenerator{
	sourceFinder: &defaultSourceFinder{},
}

func GenerateTypesFromSpecs(typeSpecs ...TypeSpec) []Type {
	return defaultAstTypeGenerator.GenerateTypesFromSpecs(typeSpecs...)
}