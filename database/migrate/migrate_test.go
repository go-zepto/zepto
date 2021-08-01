package migrate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

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
