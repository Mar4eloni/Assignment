package parser

import (
	"encoding/json"
	"regexp"
	"strings"
)

// EmailParts represents the components of a parsed email adress.
// It includes the adress specification, display name and any parsing error.
type EmailParts struct {
	AddrSpec    string  `json:"addr_spec"`    // The email adress
	DisplayName string  `json:"display_name"` // The display name
	Error       *string `json:"error"`        // Parsing error
}

// emailRegex is the regular expression for validating email addresses according the RFC 5322
var emailRegex = regexp.MustCompile(
	`(?i)^("(?:[!#-\[\]-~]|\\[\t -~])*"|[!#-'*+\-/-9=?A-Z\^-~]+(?:\.[!#-'*+\-/-9=?A-Z\^-~]+)*)@([!#-'*+\-/-9=?A-Z\^-~]+(?:\.[!#-'*+\-/-9=?A-Z\^-~]+)*|\[[!-Z\^-~]*\])$`,
)

// FormatResultsToJson converts a slice of EmailParts to a JSON formatted byte array
// Each EmailParts is seperated by blank line in the oputput.
// Returns the JSON data or an error if marshaling fails.
func FormatResultsToJson(results []EmailParts) ([]byte, error) {
	var output []byte
	for _, result := range results {
		jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			return nil, err
		}
		output = append(output, jsonData...)
		output = append(output, '\n', '\n')
	}
	return output, nil
}

// Parsline parses a single line of thext containing an email adresses and display name(optional)
// It handlesvarious formats including :
//   - "Display Name" <email@domain.com>
//   - email@domain.com
//   - (comments) email@domain.com
//
// Returns EmailParts containing the parsed components or an error if parsing fails
func ParseLine(line string) EmailParts {
	line = strings.TrimSpace(line)
	if line == "" {
		return EmailParts{Error: strPtr("empty line")}
	}

	addrSpec := extractEmailAddress(line)

	displayName := extractDisplayName(line, addrSpec)

	if addrSpec == "" {
		return EmailParts{
			Error: strPtr("no email address found"),
		}
	}

	if !emailRegex.MatchString(addrSpec) {
		return EmailParts{
			AddrSpec:    addrSpec,
			DisplayName: displayName,
			Error:       strPtr("invalid email format"),
		}
	}

	return EmailParts{
		AddrSpec:    addrSpec,
		DisplayName: displayName,
	}
}

// strPtr returns a pointer to the given string.
// This is a helper function for creating error message pointers.
func strPtr(s string) *string { return &s }

// cleanInput prepares the input string for parsing by:
// - Removing byte order marks
// - Normalizing whitespace
func cleanInput(input string) string {
	// Remove byte order mark and normalize spaces
	input = strings.TrimPrefix(input, "\ufeff")
	return strings.Join(strings.Fields(input), " ")
}

// extractDisplayName extracts the display name from an email string.
// It handles both quoted ("Name") and unquoted names.
func extractDisplayName(input, addrSpec string) string {

	//input = strings.TrimSpace(input)
	remaining := strings.Replace(input, addrSpec, "", 1)
	remaining = removeComments(remaining)

	if strings.Contains(remaining, "\"") {
		start := strings.Index(remaining, "\"")
		end := strings.LastIndex(remaining, "\"")
		if end > start {
			name := remaining[start+1 : end]
			return unescapeDisplayName(name)
		}
	}

	if before := strings.Split(remaining, "<")[0]; before != remaining {
		return strings.TrimSpace(before)
	}

	return ""
}

// extractEmailAddress extracts the mail address from a string, like:
// - Angle bracket notation (<email@domain.com>)
// - Bare email addresses
// - Comments in parentheses
func extractEmailAddress(input string) string {
	input = strings.TrimSpace(input)
	input = removeComments(input)

	if start := strings.Index(input, "<"); start != -1 {
		if end := strings.Index(input, ">"); end > start {
			candidate := strings.TrimSpace(input[start+1 : end])
			if emailRegex.MatchString(candidate) {
				return candidate
			}
		}
	}

	matches := emailRegex.FindAllString(input, -1)
	if len(matches) > 0 {
		return matches[len(matches)-1]
	}

	return ""
}

// unescapeDisplayName processes excape equences in display names.
// Currently handles \" and \\ sequences."
func unescapeDisplayName(name string) string {
	var result strings.Builder
	escape := false

	for _, r := range name {
		if escape {
			result.WriteRune(r)
			escape = false
		} else if r == '\\' {
			escape = true
		} else {
			result.WriteRune(r)
		}
	}
	return strings.TrimSpace(result.String())
}

// removeComents strips comments enclosed in paretheses from the input string,
// except when they appear witthin quoted strings.
func removeComments(input string) string {
	var result strings.Builder
	commentDepth := 0
	inQuote := false

	for _, r := range input {
		switch {
		case r == '"' && commentDepth == 0:
			inQuote = !inQuote
			result.WriteRune(r)
		case r == '(' && !inQuote:
			commentDepth++
		case r == ')' && !inQuote && commentDepth > 0:
			commentDepth--
		case commentDepth == 0:
			result.WriteRune(r)
		}
	}
	return result.String()
}
