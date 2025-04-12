package constants

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ComponentStyleType string

const (
	Function ComponentStyleType = "function"
	Const    ComponentStyleType = "const"
)

func (cst ComponentStyleType) String() string {
	return string(cst)
}

func (cst *ComponentStyleType) UnmarshalJSON(data []byte) error {
	var s string

	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("component style should be a string, got %s: %w", data, err)
	}

	value := ComponentStyleType(strings.ToLower(s))

	switch value {
	case Function, Const:
		*cst = value
		return nil
	default:
		return fmt.Errorf("invalid component style value '%s', expected '%s' or '%s'", s, Function, Const)
	}
}
