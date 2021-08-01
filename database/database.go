package database

import (
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DEFAULT_SQLITE_PATH = "db/development.sqlite3"

type Connection struct {
	Adapter    string `json:"adapter"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Datababase string `json:"database"`
}

func (c *Connection) postgresDsn() string {
	port := 5432
	if c.Port != 0 {
		port = c.Port
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%d",
		c.Host,
		c.Username,
		c.Datababase,
		port,
	)
	if c.Password != "" {
		dsn = dsn + fmt.Sprintf(" password=%s", c.Password)
	}
	return dsn
}

func (c *Connection) sqliteDsn() string {
	if c.Datababase == DEFAULT_SQLITE_PATH {
		// Create DB folder if not exists
		_ = os.Mkdir("db", os.ModePerm)
	}
	return c.Datababase
}

func (c *Connection) Open(logger logger.Interface) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch c.Adapter {
	case "postgres":
		dialector = postgres.Open(c.postgresDsn())
	case "sqlite":
		dialector = sqlite.Open(c.sqliteDsn())
	}
	if dialector == nil {
		return nil, errors.New("database adapter not supported")
	}
	return gorm.Open(dialector, &gorm.Config{
		Logger: logger,
	})
}
