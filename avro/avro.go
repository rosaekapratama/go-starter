package avro

import (
	"context"
	"github.com/hamba/avro/v2"
	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/otel"
	"github.com/rosaekapratama/go-starter/response"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const spanRegisterSchema = "avro.RegisterSchema"
const spanGetSchema = "avro.GetSchema"
const defaultAvscPath = "avro"
const avscExt = ".avsc"

const (
	errAvroSchemaManagerIsDisabled = "Avro schema manager is disabled"
	errFailedToReadDirectory       = "Failed to read directory, path=%s"
)

var (
	SchemaManager ISchemaManager
)

func init() {
	SchemaManager = &SchemaManagerImpl{schemas: &avro.SchemaCache{}}
}

func Init(ctx context.Context, _ config.Config) {
	dir, err := os.ReadDir(defaultAvscPath)
	if err != nil {
		if _, ok := err.(*fs.PathError); ok {
			log.Warn(ctx, errAvroSchemaManagerIsDisabled)
			return
		}
		log.Fatalf(ctx, err, errFailedToReadDirectory, defaultAvscPath)
		return
	}
	abs, err := filepath.Abs(defaultAvscPath)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to get absolute path, path=%s", defaultAvscPath)
		return
	}

	// Load all avsc files under default dir
	for _, entry := range dir {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), avscExt) {
			path := filepath.Join(abs, entry.Name())
			err := SchemaManager.RegisterSchema(ctx, strings.Split(entry.Name(), sym.Dot)[0], path)
			if err != nil {
				log.Fatalf(ctx, err, "Failed to register avro codec, path=%s", path)
				return
			}
			log.Infof(ctx, "Avro codec registered successfully, path=%s", path)
		}
	}
}

func (manager *SchemaManagerImpl) RegisterSchema(ctx context.Context, name string, avscPath string) error {
	ctx, span := otel.Trace(ctx, spanRegisterSchema)
	defer span.End()

	schema, err := avro.ParseFiles(avscPath)
	if err != nil {
		log.Fatalf(ctx, err, "Failed to parse schema, path=%s", avscPath)
		return err
	}
	manager.schemas.Add(name, schema)
	return nil
}

func (manager *SchemaManagerImpl) GetSchema(ctx context.Context, name string) (avro.Schema, error) {
	ctx, span := otel.Trace(ctx, spanGetSchema)
	defer span.End()

	schema := manager.schemas.Get(name)
	if schema != nil {
		return schema, nil
	} else {
		log.Errorf(ctx, response.AvroSchemaNotFound, "Avro schema not found, name=%s", name)
		return nil, response.AvroSchemaNotFound
	}
}
