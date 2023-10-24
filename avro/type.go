package avro

import (
	"context"
	"github.com/hamba/avro/v2"
)

type ISchemaManager interface {
	RegisterSchema(ctx context.Context, name string, avscPath string) error
	GetSchema(ctx context.Context, name string) (avro.Schema, error)
}

type SchemaManagerImpl struct {
	schemas *avro.SchemaCache
}
