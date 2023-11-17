package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rosaekapratama/go-starter/constant/sym"
	"github.com/rosaekapratama/go-starter/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"strconv"
	"strings"
	"time"

	"github.com/rosaekapratama/go-starter/constant/integer"

	"github.com/rosaekapratama/go-starter/config"
	"github.com/rosaekapratama/go-starter/constant/str"
	"github.com/rosaekapratama/go-starter/constant/timezone"
	"github.com/rosaekapratama/go-starter/log"
	"github.com/rosaekapratama/go-starter/response"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

const (
	errDatabaseManagerIsDisabled = "Database manager is disabled"

	PostgreSQL = "postgresql"
	MySQL      = "mysql"
	MSSQL      = "mssql"

	PostgreSQLPort = 5432
	MySQLPort      = 3306
	MSSQLPort      = 1433
)

var (
	errInvalidConnectionId = errors.New("connection Is Invalid")

	logrusLevelMap = map[logrus.Level]logger.LogLevel{
		logrus.PanicLevel: logger.Error,
		logrus.FatalLevel: logger.Error,
		logrus.ErrorLevel: logger.Error,
		logrus.WarnLevel:  logger.Warn,
		logrus.InfoLevel:  logger.Warn,
		logrus.DebugLevel: logger.Info,
		logrus.TraceLevel: logger.Info,
	}

	Manager IManager
)

func Init(ctx context.Context, config config.Config) {
	logLevel := log.GetLogger().GetLevel()
	cfg := config.GetObject().Database

	if cfg == nil || len(cfg) == integer.Zero {
		log.Warn(ctx, errDatabaseManagerIsDisabled)
		return
	}

	gormDBMap := make(map[string]*gorm.DB)
	sqlDBMap := make(map[string]*sql.DB)
	for id, databaseConfig := range cfg {
		if databaseConfig.Driver == str.Empty {
			log.Fatal(ctx, response.ConfigNotFound, "Missing database driver")
			return
		}

		if databaseConfig.Address == str.Empty {
			log.Fatal(ctx, response.ConfigNotFound, "Missing database address")
			return
		}
		if databaseConfig.Database == str.Empty {
			log.Fatal(ctx, response.ConfigNotFound, "Missing database name")
			return
		}
		if databaseConfig.Username == str.Empty {
			log.Warn(ctx, response.ConfigNotFound, "Missing database username")
		}

		if databaseConfig.Password == str.Empty {
			log.Warn(ctx, response.ConfigNotFound, "Missing database password")
		}

		address := strings.Split(databaseConfig.Address, sym.Colon)
		host := address[integer.Zero]
		var port int
		var err error
		var dialector gorm.Dialector
		switch databaseConfig.Driver {
		case PostgreSQL:
			if len(address) > integer.One {
				port, err = strconv.Atoi(address[integer.One])
				if err != nil {
					log.Fatalf(ctx, response.InvalidConfig, "Invalid database port, host=%s, port=%s", host, address[integer.One])
					return
				}
			} else {
				port = PostgreSQLPort
			}
			dialector = postgres.Open(
				fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
					host,
					databaseConfig.Username,
					databaseConfig.Password,
					databaseConfig.Database,
					port,
					timezone.AsiaJakarta,
				),
			)
		case MySQL:
			if len(address) > integer.One {
				port, err = strconv.Atoi(address[integer.One])
				if err != nil {
					log.Fatalf(ctx, response.InvalidConfig, "invalid database port, host=%s, port=%s", host, address[integer.One])
					return
				}
			} else {
				port = MySQLPort
			}
			dialector = mysql.Open(
				fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
					databaseConfig.Username,
					databaseConfig.Password,
					host,
					port,
					databaseConfig.Database,
				),
			)
		case MSSQL:
			if len(address) > integer.One {
				port, err = strconv.Atoi(address[integer.One])
				if err != nil {
					log.Fatalf(ctx, response.InvalidConfig, "invalid database port, host=%s, port=%s", host, address[integer.One])
					return
				}
			} else {
				port = MSSQLPort
			}
			dialector = sqlserver.Open(
				fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
					databaseConfig.Username,
					databaseConfig.Password,
					host,
					port,
					databaseConfig.Database,
				),
			)
		default:
			log.Fatalf(ctx, response.ConfigNotFound, "Unsupported database driver %s", databaseConfig.Driver)
			return
		}

		gormLoggerConfig := logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logrusLevelMap[logLevel],
			IgnoreRecordNotFoundError: databaseConfig.IgnoreRecordNotFoundError,
			Colorful:                  utils.IsRunLocally(ctx),
		}
		if databaseConfig.SlowThreshold > integer.Zero {
			gormLoggerConfig.SlowThreshold = time.Duration(databaseConfig.SlowThreshold) * time.Millisecond
		}
		gormConfig := gorm.Config{
			SkipDefaultTransaction: databaseConfig.SkipDefaultTransaction,
			Logger:                 logger.New(log.GetLogger().GetLogrusLogger(), gormLoggerConfig),
		}

		gormDB, err := gorm.Open(dialector, &gormConfig)
		if err != nil {
			log.Fatalf(ctx, err, "Failed to connect to database, id=%s", id)
			return
		}

		sqlDB, err := gormDB.DB()
		if err != nil {
			log.Fatalf(ctx, err, "Failed to get generic database object *sql.DB, id=%s", id)
			return
		}

		if databaseConfig.Conn != nil {
			if databaseConfig.Conn.MaxIdle > integer.Zero {
				// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
				sqlDB.SetMaxIdleConns(databaseConfig.Conn.MaxIdle)
			}

			if databaseConfig.Conn.MaxOpen > integer.Zero {
				// SetMaxOpenConns sets the maximum number of open connections to the database.
				sqlDB.SetMaxOpenConns(databaseConfig.Conn.MaxOpen)
			}

			if databaseConfig.Conn.MaxLifeTime > integer.Zero {
				// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
				sqlDB.SetConnMaxLifetime(time.Duration(databaseConfig.Conn.MaxLifeTime) * time.Millisecond)
			}
		}

		gormDBMap[id] = gormDB
		sqlDBMap[id] = sqlDB
	}

	Manager = &ManagerImpl{
		gormDBMap: gormDBMap,
		sqlDBMap:  sqlDBMap,
	}
}

func (manager *ManagerImpl) GetConnectionIds() []string {
	var keys []string
	for k := range manager.gormDBMap {
		keys = append(keys, k)
	}

	return keys
}

func (manager *ManagerImpl) DB(ctx context.Context, connectionID string) (*gorm.DB, *sql.DB, error) {

	if gormDB, ok := manager.gormDBMap[connectionID]; ok {
		if sqlDB, ok := manager.sqlDBMap[connectionID]; ok {
			return gormDB, sqlDB, nil
		}
	}

	log.Errorf(ctx, response.GeneralError, "*gorm.DB and *sql.DB are not found, id=%s", connectionID)
	return nil, nil, errInvalidConnectionId
}
