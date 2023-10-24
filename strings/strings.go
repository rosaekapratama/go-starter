package strings

import (
	"github.com/rosaekapratama/go-starter/constant/str"
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
