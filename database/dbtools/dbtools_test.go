package dbtools

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAndDropDB(t *testing.T) {
	if os.Getenv("POSTGRES_TEST") != "true" {
		t.Skip("Skipping tests that requires postgres")
	}
	os.Setenv("ZEPTO_DB_ADAPTER", "postgres")
	os.Setenv("ZEPTO_DB_HOST", "127.0.0.1")
	os.Setenv("ZEPTO_DB_DATABASE", "mycustomdb")
	os.Setenv("ZEPTO_DB_PORT", "15432")
	os.Setenv("ZEPTO_DB_SSLMODE", "disable")
	os.Setenv("ZEPTO_DB_USERNAME", "postgres")
	os.Setenv("ZEPTO_DB_PASSWORD", "postgres")
	dt, err := NewDBTools()
	require.NoError(t, err)
	require.NotNil(t, dt)
	err = dt.CreateDB()
	require.NoError(t, err)
	err = dt.DropDB()
	require.NoError(t, err)
}
