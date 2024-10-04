package utils

import (
	"github.com/inhies/go-bytesize"
	"runtime"
)

func GetMemStats() bytesize.ByteSize {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return bytesize.New(float64(m.Alloc))
}

func GetMemStatsInMB() string {
	return GetMemStats().Format("%f", "MB", false)
}

func GetMemStatsInKB() string {
	return GetMemStats().Format("%f", "KB", false)
}
