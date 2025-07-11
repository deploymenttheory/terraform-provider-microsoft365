package convert

import (
	"errors"
	"strings"
)

// MockBitmaskEnum represents a mock bitmask enum for testing
type MockBitmaskEnum int

const (
	MockBitmaskNone MockBitmaskEnum = 1
	MockBitmaskOne  MockBitmaskEnum = 2
	MockBitmaskTwo  MockBitmaskEnum = 4
	MockBitmaskAll  MockBitmaskEnum = 7
)

func (e MockBitmaskEnum) String() string {
	var values []string
	options := []string{"none", "one", "two"}
	for p := 0; p < 3; p++ {
		mask := MockBitmaskEnum(1 << p)
		if e&mask == mask {
			values = append(values, options[p])
		}
	}
	return strings.Join(values, ",")
}

// MockParseBitmaskEnum simulates parsing a comma-separated string into a bitmask enum
func MockParseBitmaskEnum(input string) (any, error) {
	if input == "" {
		return nil, nil
	}

	var result MockBitmaskEnum
	parts := strings.Split(input, ",")

	for _, part := range parts {
		switch strings.TrimSpace(part) {
		case "none":
			result |= MockBitmaskNone
		case "one":
			result |= MockBitmaskOne
		case "two":
			result |= MockBitmaskTwo
		default:
			return nil, errors.New("invalid enum value: " + part)
		}
	}

	return &result, nil
}
