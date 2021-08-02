package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GeneratePostgresDSN(t *testing.T) {
	c := Connection{
		Adapter:  "postgres",
		Host:     "177.50.20.100",
		Port:     1234,
		Database: "metropolis",
		Username: "clark",
		Password: "super-password",
	}
	dsn, err := c.GenerateAdapterDSN()
	require.NoError(t, err)
	require.Equal(t, "dbname=metropolis host=177.50.20.100 password=super-password port=1234 user=clark", dsn)
}

func Test_GeneratePostgresDSN_No_Password(t *testing.T) {
	c := Connection{
		Adapter:  "postgres",
		Host:     "177.50.20.100",
		Port:     1234,
		Database: "metropolis",
		Username: "clark",
	}
	dsn, err := c.GenerateAdapterDSN()
	require.NoError(t, err)
	require.Equal(t, "dbname=metropolis host=177.50.20.100 port=1234 user=clark", dsn)
}

func Test_GeneratePostgresDSN_No_Username(t *testing.T) {
	c := Connection{
		Adapter:  "postgres",
		Host:     "177.50.20.100",
		Port:     1234,
		Database: "metropolis",
		Password: "super-password",
	}
	dsn, err := c.GenerateAdapterDSN()
	require.NoError(t, err)
	require.Equal(t, "dbname=metropolis host=177.50.20.100 port=1234", dsn)
}

func Test_GenerateSqliteDSN(t *testing.T) {
	c := Connection{
		Adapter:  "sqlite",
		Database: "./db/development.sqlite3",
	}
	dsn, err := c.GenerateAdapterDSN()
	require.NoError(t, err)
	require.Equal(t, "./db/development.sqlite3", dsn)
}

func Test_GenerateSqlite3DSN(t *testing.T) {
	c := Connection{
		Adapter:  "sqlite3",
		Database: "./db/development.sqlite3",
	}
	dsn, err := c.GenerateAdapterDSN()
	require.NoError(t, err)
	require.Equal(t, "./db/development.sqlite3", dsn)
}
