package parser

import (
	"testing"
)

func TestParseLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected EmailParts
	}{
		{
			name:  "test name",
			input: "test@test.com",
			expected: EmailParts{
				AddrSpec:    "test@test.com",
				DisplayName: "",
				Error:       nil,
			},
		},
		{
			name:  "jonny joy doe",
			input: `"John Doe" <john@example.com>`,
			expected: EmailParts{
				AddrSpec:    "john@example.com",
				DisplayName: "John Doe",
				Error:       nil,
			},
		},
		{
			name:  "empty line",
			input: "",
			expected: EmailParts{
				AddrSpec:    "",
				DisplayName: "",
				Error:       strPtr("empty line"),
			},
		},
		{
			name:  "invalid email format",
			input: "not-an-email",
			expected: EmailParts{
				AddrSpec:    "not-an-email",
				DisplayName: "",
				Error:       strPtr("invalid email format"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseLine(tt.input)

			// Compare AddrSpec
			if result.AddrSpec != tt.expected.AddrSpec {
				t.Errorf("AddrSpec mismatch: got %q, want %q", result.AddrSpec, tt.expected.AddrSpec)
			}

			// Compare DisplayName
			if result.DisplayName != tt.expected.DisplayName {
				t.Errorf("DisplayName mismatch: got %q, want %q", result.DisplayName, tt.expected.DisplayName)
			}

			// Compare Error
			if (result.Error == nil) != (tt.expected.Error == nil) {
				t.Errorf("Error presence mismatch: got %v, want %v", result.Error, tt.expected.Error)
			}
			if result.Error != nil && tt.expected.Error != nil && *result.Error != *tt.expected.Error {
				t.Errorf("Error message mismatch: got %q, want %q", *result.Error, *tt.expected.Error)
			}
		})
	}
}
