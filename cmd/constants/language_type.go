package constants

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LanguageType string

const (
	Typescript LanguageType = "ts"
	Javascript LanguageType = "js"
)

func (lt LanguageType) String() string {
	return string(lt)
}

func (lt *LanguageType) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("language should be a string, got %s: %w", data, err)
	}

	value := LanguageType(strings.ToLower(s))

	switch value {
	case Typescript, Javascript:
		*lt = value
		return nil
	default:
		return fmt.Errorf("invalid language value '%s', expected '%s' or '%s'", s, Typescript, Javascript)
	}
}
