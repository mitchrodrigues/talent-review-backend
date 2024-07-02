package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func InterfaceToCamelCaseDelimited(event any, delimiter rune) string {
	name := reflect.TypeOf(event).Name()
	return CamelCaseToDelimited(name, delimiter)
}

// CamelCaseToDelimited converts camelCase strings to a delimited format.
// For example, "ClientConnected" becomes "client:connected".
func CamelCaseToDelimited(s string, delimiter rune) string {
	var sb strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				sb.WriteRune(delimiter)
			}
			sb.WriteRune(unicode.ToLower(r))
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func Coalesce(s1, s2 string) string {
	if s1 != "" {
		return s1
	}

	return s2
}

// GenerateCode generates a 10-digit alphanumeric code in the format XXXXX-XXXXX
func GenerateCode() (string, error) {
	// Generate 5 bytes of random data (10 hex characters)
	bytes := make([]byte, 5)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Convert bytes to a hexadecimal string
	code := hex.EncodeToString(bytes)

	// Insert hyphen after 5 characters
	return strings.ToUpper(fmt.Sprintf("%s-%s", code[:5], code[5:])), nil
}
