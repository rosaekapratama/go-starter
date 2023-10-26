package utils

import (
	"context"
	"github.com/rosaekapratama/go-starter/constant/env"
	"github.com/rosaekapratama/go-starter/constant/integer"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/log"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func BoolP(b bool) *bool {
	return &b
}

func PBool(b *bool) bool {
	return *b
}

func StringP(s string) *string {
	if s == str.Empty {
		return nil
	}
	return &s
}

func PString(s *string) string {
	if s == nil {
		return str.Empty
	}
	return *s
}

func IntP(i int) *int {
	return &i
}

func PInt(i *int) int {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Int32P(i int32) *int32 {
	return &i
}

func PInt32(i *int32) int32 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Int64P(i int64) *int64 {
	return &i
}

func PInt64(i *int64) int64 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func UintP(i uint) *uint {
	return &i
}

func PUint(i *uint) uint {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Uint32P(i uint32) *uint32 {
	return &i
}

func PUint32(i *uint32) uint32 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func Uint64P(i uint64) *uint64 {
	return &i
}

func PUint64(i *uint64) uint64 {
	if i == nil {
		return integer.Zero
	}
	return *i
}

func GenerateSalt(ctx context.Context, size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := rand.Read(salt)
	if err != nil {
		log.Error(ctx, err)
		return nil, err
	}
	return salt, nil
}

func StringToSliceOfString(s string, sep string) []string {
	return strings.Split(s, sep)
}

func StringToSliceOfInt(ctx context.Context, s string, sep string) ([]int, error) {
	ints := make([]int, integer.Zero)
	for _, s := range strings.Split(s, sep) {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Errorf(ctx, err, "Failed parse '%s' to uint64", s)
			return nil, err
		}
		ints = append(ints, i)
	}
	return ints, nil
}

func StringToSliceOfInt16(s string, sep string) ([]int16, error) {
	uints := make([]int16, integer.Zero)
	for _, s := range strings.Split(s, sep) {
		if s == str.Empty {
			continue
		}
		i, err := strconv.ParseInt(s, integer.Ten, integer.I16)
		if err != nil {
			log.Errorf(context.Background(), err, "Failed parse '%s' to uint64", s)
			return nil, err
		}
		uints = append(uints, int16(i))
	}
	return uints, nil
}

func StringToSliceOfUint64(s string, sep string) ([]uint64, error) {
	uints := make([]uint64, integer.Zero)
	for _, s := range strings.Split(s, sep) {
		if s == str.Empty {
			continue
		}
		i, err := strconv.ParseUint(s, integer.Ten, integer.I64)
		if err != nil {
			log.Errorf(context.Background(), err, "Failed parse '%s' to uint64", s)
			return nil, err
		}
		uints = append(uints, i)
	}
	return uints, nil
}

func StringToSliceOfFloat64(s string, sep string) ([]float64, error) {
	uints := make([]float64, integer.Zero)
	for _, s := range strings.Split(s, sep) {
		if s == str.Empty {
			continue
		}
		i, err := strconv.ParseFloat(s, integer.I64)
		if err != nil {
			log.Errorf(context.Background(), err, "Failed parse '%s' to uint64", s)
			return nil, err
		}
		uints = append(uints, i)
	}
	return uints, nil
}

func SliceOfStringToString(strs []string, sep string) string {
	v := strings.Builder{}
	for i, s := range strs {
		v.WriteString(s)
		if i < len(strs)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func SliceOfIntToString(ints []int, sep string) string {
	v := strings.Builder{}
	for idx, in := range ints {
		v.WriteString(strconv.Itoa(in))
		if idx < len(ints)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func SliceOfInt16ToString(ints []int16, sep string) string {
	v := strings.Builder{}
	for i, u := range ints {
		v.WriteString(strconv.FormatInt(int64(u), integer.Ten))
		if i < len(ints)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func SliceOfUint64ToString(uints []uint64, sep string) string {
	v := strings.Builder{}
	for i, u := range uints {
		v.WriteString(strconv.FormatUint(u, integer.Ten))
		if i < len(uints)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func SliceOfFloat64ToString(floats []float64, sep string) string {
	v := strings.Builder{}
	for i, u := range floats {
		v.WriteString(strconv.FormatFloat(u, 'f', integer.Two, integer.I64))
		if i < len(floats)-1 {
			v.WriteString(sep)
		}
	}
	return v.String()
}

func IsRunLocally(ctx context.Context) bool {
	if localRunStr, ok := os.LookupEnv(env.EnvLocalRun); localRunStr != str.Empty && ok {
		localRun, err := strconv.ParseBool(localRunStr)
		if err != nil {
			log.Warnf(ctx, "Failed to parse %s env var '%s' to boolean, %s", env.EnvLocalRun, localRunStr, err.Error())
		} else {
			return localRun
		}
	}

	return false
}

func IsZeroValue(v interface{}) bool {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true
		}
		val = val.Elem()
	}
	zeroValue := reflect.Zero(val.Type()).Interface()
	return reflect.DeepEqual(val.Interface(), zeroValue)
}
