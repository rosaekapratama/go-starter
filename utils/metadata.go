package utils

import (
	"github.com/rosaekapratama/go-starter/constant/sym"
	"google.golang.org/grpc/metadata"
)

// TruncateMetadata truncates metadata values to a specified max length.
func TruncateMetadata(md metadata.MD, maxLength int) metadata.MD {
	truncatedMD := metadata.MD{}
	for key, values := range md {
		truncatedValues := make([]string, len(values))
		for i, value := range values {
			if len(value) > maxLength {
				truncatedValues[i] = value[:maxLength] + sym.Ellipsis
			} else {
				truncatedValues[i] = value
			}
		}
		truncatedMD[key] = truncatedValues
	}
	return truncatedMD
}
