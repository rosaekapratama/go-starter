package strings

import (
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"strings"
	"unicode"
)

func ReturnEmptyIfNil(s *string) string {
	if s != nil {
		return *s
	} else {
		return str.Empty
	}
}

func ReturnZeroIfEmpty(s string) string {
	if s == str.Empty {
		return str.Zero
	} else {
		return s
	}
}

// SnakeToCamel converts a snake_case string to camelCase.
func SnakeToCamel(snake string) string {
	// Split the snake_case string into words
	words := strings.Split(snake, sym.Underscore)

	// If there are no words, return an empty string
	if len(words) == 0 {
		return str.Empty
	}

	// Process the first word (lowercase in camelCase)
	camel := words[0]

	// Capitalize the first letter of each subsequent word
	for _, word := range words[1:] {
		if len(word) > 0 {
			camel += strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return camel
}

// CamelToSnake converts a camelCase string to snake_case.
func CamelToSnake(camel string) string {
	var snake strings.Builder

	for i, ch := range camel {
		if unicode.IsUpper(ch) {
			if i > 0 {
				snake.WriteString(sym.Underscore) // Add an underscore before uppercase letters except at the start
			}
			snake.WriteRune(unicode.ToLower(ch)) // Convert uppercase to lowercase
		} else {
			snake.WriteRune(ch) // Directly write the rune
		}
	}

	return snake.String()
}

// DashToCamel converts a dash-separated string (like "Terminal-Id") to camelCase.
func DashToCamel(dash string) string {
	// Split the string by dashes
	words := strings.Split(dash, sym.Hyphen)

	// If there are no words, return an empty string
	if len(words) == 0 {
		return str.Empty
	}

	// Process the first word (lowercase in camelCase)
	camel := words[0]

	// Capitalize the first letter of each subsequent word
	for _, word := range words[1:] {
		if len(word) > 0 {
			camel += strings.ToUpper(string(word[0])) + word[1:]
		}
	}

	return camel
}

// DashToSnake converts a dash-separated string (like "Terminal-Id") to snake_case.
func DashToSnake(dash string) string {
	// Replace all dashes with underscores
	snake := strings.ReplaceAll(dash, sym.Hyphen, sym.Underscore)

	// Convert the entire string to lowercase
	return strings.ToLower(snake)
}
