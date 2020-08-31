package gotype

var defaultAstTypeGenerator = &astTypeGenerator{
	sourceFinder: &defaultSourceFinder{},
}

func GenerateTypesFromSpecs(typeSpecs ...TypeSpec) ([]Type, error) {
	return defaultAstTypeGenerator.GenerateTypesFromSpecs(typeSpecs...)
}