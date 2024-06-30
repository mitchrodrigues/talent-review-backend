package helpers

import (
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
