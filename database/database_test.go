package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_PostgresDSN(t *testing.T) {
	c := Connection{
		Adapter:    "postgres",
		Host:       "177.50.20.100",
		Port:       1234,
		Datababase: "metropolis",
		Username:   "clark",
		Password:   "super-password",
	}
	require.Equal(
		t,
		c.postgresDsn(),
		"host=177.50.20.100 user=clark dbname=metropolis port=1234 password=super-password",
	)
}
