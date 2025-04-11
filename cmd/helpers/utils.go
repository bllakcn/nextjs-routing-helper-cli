package helpers

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func titleCase(s string) string {
	caser := cases.Title(language.Und) // This creates a caser that works for any language
	return caser.String(s)
}

func ToPascalCase(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	s = titleCase(s)
	return strings.ReplaceAll(s, " ", "")
}
