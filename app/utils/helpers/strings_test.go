package helpers

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCamelCaseToDelimited(t *testing.T) {
	testCases := []struct {
		name      string
		input     string
		delimiter rune
		expected  string
	}{
		{"Simple", "ClientConnected", ':', "client:connected"},
		{"Empty", "", ':', ""},
		{"NoCamelCase", "client", ':', "client"},
		{"MultipleWords", "MySampleText", ':', "my:sample:text"},
		{"StartingUpper", "StartingUpper", ':', "starting:upper"},
		{"AllUpper", "ALLUPPER", ':', "a:l:l:u:p:p:e:r"},
		{"MixedCase", "JSONDataFormat", ':', "j:s:o:n:data:format"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CamelCaseToDelimited(tc.input, tc.delimiter)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGenerateCode(t *testing.T) {
	code, err := GenerateCode()
	if err != nil {
		t.Fatalf("GenerateCode returned an error: %v", err)
	}

	// Check the length of the code
	if len(code) != 11 {
		t.Fatalf("Expected code length of 11, got %d", len(code))
	}

	// Check the format of the code using a regular expression
	matched, err := regexp.MatchString(`^[A-F0-9]{5}-[A-F0-9]{5}$`, code)
	if err != nil {
		t.Fatalf("Error matching regex: %v", err)
	}
	if !matched {
		t.Fatalf("Code format is incorrect: %s", code)
	}
}
