package database

import (
	"errors"
	"fmt"

	"github.com/go-zepto/zepto/config"
	"github.com/xo/dburl"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Connection struct {
	Adapter  string `json:"adapter"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
	SSLMode  string `json:"sslmode"`
}

func NewConnectionDataFromConfig(c *config.Config) *Connection {
	adapter := c.DB.Adapter
	if adapter == "sqlite" {
		adapter = "sqlite3"
	}
	return &Connection{
		Adapter:  adapter,
		Host:     c.DB.Host,
		Port:     c.DB.Port,
		Username: c.DB.Username,
		Password: c.DB.Password,
		Database: c.DB.Database,
		SSLMode:  c.DB.SSLMode,
	}
}

func (c *Connection) GenerateAdapterDBURL() (*dburl.URL, error) {
	userAndPass := ""
	if c.Username != "" {
		userAndPass += c.Username
		if c.Password != "" {
			userAndPass += ":" + c.Password
		}
		userAndPass += "@"
	}
	port := ""
	if c.Port != 0 {
		port = fmt.Sprintf(":%d", c.Port)
	}
	query := "?"
	if c.SSLMode != "" {
		query = query + "sslmode=" + c.SSLMode
	}
	url := fmt.Sprintf(
		"%s://%s%s%s/%s%s",
		c.Adapter,
		userAndPass,
		c.Host,
		port,
		c.Database,
		query,
	)
	if c.Adapter == "sqlite" || c.Adapter == "sqlite3" {
		url = fmt.Sprintf("sqlite3://%s", c.Database)
	}
	dsnURL, err := dburl.Parse(url)
	if err != nil {
		return nil, err
	}
	return dsnURL, nil
}

func (c *Connection) GenerateAdapterDSN() (string, error) {
	u, err := c.GenerateAdapterDBURL()
	if err != nil {
		return "", err
	}
	return u.DSN, nil
}

func (c *Connection) Open(logger logger.Interface) (*gorm.DB, error) {
	var dialector gorm.Dialector
	dsn, err := c.GenerateAdapterDSN()
	if err != nil {
		return nil, err
	}
	switch c.Adapter {
	case "postgres":
		dialector = postgres.Open(dsn)
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(dsn)
	}
	if dialector == nil {
		return nil, errors.New("database adapter not supported")
	}
	return gorm.Open(dialector, &gorm.Config{
		Logger: logger,
	})
}
