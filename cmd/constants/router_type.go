package constants

import (
	"encoding/json"
	"fmt"
	"strings"
)

type RouterType string

const (
	AppRouter   RouterType = "app"
	PagesRouter RouterType = "pages"
)

func (rt RouterType) String() string {
	return string(rt)
}

// IsValid checks if the RouterType instance is one of the defined constants.
func (rt RouterType) IsValid() bool {
	switch rt {
	case AppRouter, PagesRouter:
		return true
	default:
		return false
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// This method provides custom validation when decoding JSON into RouterType.
func (rt *RouterType) UnmarshalJSON(data []byte) error {
	var s string
	// First, unmarshal the JSON string into a temporary string variable
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("router type should be a string, got %s: %w", data, err)
	}

	// Convert the temporary string to your RouterType for comparison
	value := RouterType(strings.ToLower(s)) // Use lower case for case-insensitive comparison

	// Validate against your defined constants
	switch value {
	case AppRouter, PagesRouter:
		*rt = value // Assign the valid value to the target RouterType pointer
		return nil  // Success
	default:
		// The value from JSON is not one of the allowed types
		return fmt.Errorf("invalid router type value '%s', expected '%s' or '%s'", s, AppRouter, PagesRouter)
	}
}
