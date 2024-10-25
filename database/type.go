package database

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
)

type IManager interface {
	GetConnectionIds() []string
	DB(ctx context.Context, connectionId string) (*gorm.DB, *sql.DB, error)
	Begin(ctx context.Context, connectionId string) (*gorm.DB, error)
}

type managerImpl struct {
	gormDBMap map[string]*gorm.DB
	sqlDBMap  map[string]*sql.DB
}
