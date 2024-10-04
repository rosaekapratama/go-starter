package database

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type IManager interface {
	GetConnectionIds() []string
	DB(ctx context.Context, connectionID string) (*gorm.DB, *sql.DB, error)
}

type managerImpl struct {
	gormDBMap map[string]*gorm.DB
	sqlDBMap  map[string]*sql.DB
}
