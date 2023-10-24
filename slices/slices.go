package slices

import (
	"github.com/rosaekapratama/go-starter/constant/integer"
	"strings"
)

func ContainStringCaseInsensitive(slice []string, elem string) bool {
	for _, v := range slice {
		if strings.ToLower(elem) == strings.ToLower(v) {
			return true
		}
	}
	return false
}

func ContainStringCaseSensitive(slice []string, elem string) bool {
	for _, v := range slice {
		if elem == v {
			return true
		}
	}
	return false
}

func ContainStringsCaseInsensitive(slice []string, elems ...string) bool {
	for _, e := range elems {
		var result bool
		for _, v := range slice {
			if strings.ToLower(e) == strings.ToLower(v) {
				result = true
				break
			}
		}

		if !result {
			return false
		}
	}
	return true
}

func ContainStringsCaseSensitive(slice []string, elems ...string) bool {
	for _, e := range elems {
		var result bool
		for _, v := range slice {
			if e == v {
				result = true
				break
			}
		}

		if !result {
			return false
		}
	}
	return true
}

func ContainIn64(slice []int64, elem int64) bool {
	for _, v := range slice {
		if elem == v {
			return true
		}
	}
	return false
}

func ContainInt64s(slice []int64, elems ...int64) bool {
	for _, e := range elems {
		var result bool
		for _, v := range slice {
			if e == v {
				result = true
				break
			}
		}

		if !result {
			return false
		}
	}
	return true
}

func ContainUint64(slice []uint64, elem uint64) bool {
	for _, v := range slice {
		if elem == v {
			return true
		}
	}
	return false
}

func ContainUint64s(slice []uint64, elems ...uint64) bool {
	for _, e := range elems {
		var result bool
		for _, v := range slice {
			if e == v {
				result = true
				break
			}
		}

		if !result {
			return false
		}
	}
	return true
}

func GetNotContainedUint64s(slice []uint64, elems ...uint64) []uint64 {
	result := make([]uint64, integer.Zero)
	for _, e := range elems {
		var exists bool
		for _, v := range slice {
			if e == v {
				exists = true
				break
			}
		}

		if !exists {
			result = append(result, e)
		}
	}
	return result
}

func AppendUniqueUint64(slice []uint64, elems ...uint64) []uint64 {
	for _, e := range elems {
		var exists bool
		for _, v := range slice {
			if e == v {
				exists = true
				break
			}
		}

		if !exists {
			slice = append(slice, e)
		}
	}

	return slice
}

func ToString(slice []string, delimiter string) string {
	sb := strings.Builder{}
	for i, s := range slice {
		sb.WriteString(s)
		if i+1 < len(slice) {
			sb.WriteString(delimiter)
		}
	}
	return sb.String()
}
