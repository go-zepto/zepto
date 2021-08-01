package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var PRODUCTS_VERSION uint = 20210801150315
var USERS_VERSION uint = 20210805150315

func clean() {
	os.RemoveAll("./db/development.sqlite3")
	files, err := filepath.Glob("./db/migrate/*.sql")
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
}

func TestCreate(t *testing.T) {
	clean()
	m, err := NewMigrate(Options{})
	require.NoError(t, err)
	require.NotNil(t, m)
	m.CreateMigrationFiles(CreateMigrationFilesOptions{
		Name: "products",
	})
	clean()
}

func TestUp(t *testing.T) {
	clean()
	m, err := NewMigrate(Options{
		dir: "testdata/db/migrate",
	})
	require.NoError(t, err)
	err = m.Up()
	require.NoError(t, err)
}

func TestUpSteps(t *testing.T) {
	clean()
	m, err := NewMigrate(Options{
		dir: "testdata/db/migrate",
	})
	require.NoError(t, err)
	err = m.UpSteps(1)
	require.NoError(t, err)
}

func TestDown(t *testing.T) {
	clean()
	m, err := NewMigrate(Options{
		dir: "testdata/db/migrate",
	})
	require.NoError(t, err)
	err = m.Up()
	require.NoError(t, err)
	err = m.Down()
}

func TestDownSteps(t *testing.T) {
	clean()
	m, err := NewMigrate(Options{
		dir: "testdata/db/migrate",
	})
	require.NoError(t, err)
	err = m.UpSteps(1)
	require.NoError(t, err)
	err = m.DownSteps(1)
}

func TestGetLatestVersion(t *testing.T) {
	clean()
	m, err := NewMigrate(Options{
		dir: "testdata/db/migrate",
	})
	require.NoError(t, err)
	latestVersion, err := m.getLatestVersion()
	require.NoError(t, err)
	require.Equal(t, latestVersion, uint(USERS_VERSION))
}

func TestStatus(t *testing.T) {
	clean()
	m, err := NewMigrate(Options{
		dir: "testdata/db/migrate",
	})
	require.NoError(t, err)
	err = m.UpSteps(1)
	require.NoError(t, err)
	status, err := m.Status()
	require.NoError(t, err)
	require.EqualValues(t, status.CurrentVersion, uint(PRODUCTS_VERSION))
	require.EqualValues(t, status.CurrentVersionFile, "20210801150315_products.up.sql")
	require.EqualValues(t, status.LatestVersion, uint(USERS_VERSION))
	require.True(t, status.Pending)
	err = m.UpSteps(1)
	require.NoError(t, err)
	status, err = m.Status()
	require.NoError(t, err)
	require.EqualValues(t, status.CurrentVersion, uint(USERS_VERSION))
	require.False(t, status.Pending)
}
