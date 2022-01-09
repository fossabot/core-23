package transport

import (
	"regexp"

	"github.com/nasermirzaei89/core/internal/core"
)

func isValidName(name string) bool {
	return regexp.MustCompile(core.NameRegex).MatchString(name)
}

func isValidType(typ string) bool {
	return regexp.MustCompile(core.TypeRegex).MatchString(typ)
}
