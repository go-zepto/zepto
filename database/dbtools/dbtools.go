package dbtools

import (
	"database/sql"
	"errors"

	"github.com/go-zepto/zepto/config"
	"github.com/go-zepto/zepto/database"
	_ "github.com/lib/pq"
)

const DEFAULT_MANAGEMENT_PG_DB_NAME = "postgres"

var ErrAdapterNotSupported = errors.New("adapter not supported")

type DBTools struct {
	db         *sql.DB
	connConfig *database.Connection
	dbName     string
}

func NewDBTools() (*DBTools, error) {
	c, err := config.NewConfigFromFile()
	if err != nil {
		return nil, err
	}
	connConfig := database.NewConnectionDataFromConfig(c)
	dbName := connConfig.Database
	connConfig.Database = DEFAULT_MANAGEMENT_PG_DB_NAME
	dburl, err := connConfig.GenerateAdapterDBURL()
	if err != nil {
		return nil, err
	}
	db, err := sql.Open(connConfig.Adapter, dburl.DSN)
	if err != nil {
		return nil, err
	}
	return &DBTools{
		db:         db,
		connConfig: connConfig,
		dbName:     dbName,
	}, nil
}

func (dt *DBTools) DropDB() error {
	if dt.connConfig.Adapter != "postgres" {
		return ErrAdapterNotSupported
	}
	_, err := dt.db.Exec("DROP DATABASE " + dt.dbName)
	return err
}

func (dt *DBTools) CreateDB() error {
	if dt.connConfig.Adapter != "postgres" {
		return ErrAdapterNotSupported
	}
	_, err := dt.db.Exec("CREATE DATABASE " + dt.dbName)
	return err
}
