package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rosaekapratama/go-starter/database"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupMockDB sets up the mock gorm.DB and sql.DB
func SetupMockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	// Create sqlmock database connection
	sqlDB, mockSqlDB, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	// Initialize GORM with the sqlmock DB and the postgres driver
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB, // Use the sql.DB object created by sqlmock
	}), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	// Mock the manager and its func call
	mockManager := new(MockIManager)
	mockManager.On("GetConnectionIds").Return([]string{"test"})
	mockManager.On("DB", mock.Anything, mock.Anything).Return(gormDB, mockSqlDB, nil)
	mockManager.On("Begin", mock.Anything, mock.Anything).Return(gormDB, nil)
	database.Manager = mockManager

	return gormDB, mockSqlDB, nil
}
