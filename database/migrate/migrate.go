package migrate

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/go-zepto/zepto/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/iancoleman/strcase"
)

const defaultTimeFormat = "20060102150405"

var defaultMigrationsDir = path.Join("db", "migrate")

type Migrate struct {
	mig *migrate.Migrate
	dir string
}

type Options struct {
	dir string
}

type CreateMigrationFilesOptions struct {
	Name    string
	UpSQL   string
	DownSQL string
}

type StatusResponse struct {
	CurrentVersion     uint   `json:"current_version"`
	CurrentVersionFile string `json:"current_version_file"`
	LatestVersion      uint   `json:"latest_version"`
	Dirty              bool   `json:"dity"`
	Pending            bool   `json:"pending"`
}

func databaseURLFromConfig() (string, error) {
	c, err := config.NewConfigFromFile()
	if err != nil {
		return "", err
	}
	db := c.DB
	if db.Adapter == "sqlite" {
		return fmt.Sprintf("sqlite3://%s", db.Database), nil
	}
	userAndPassword := c.DB.Username
	if db.Password != "" {
		userAndPassword = userAndPassword + ":" + db.Password
	}
	dbURL := fmt.Sprintf(
		"%s://%s@%s:%d/%s",
		db.Adapter,
		userAndPassword,
		db.Host,
		db.Port,
		db.Database,
	)
	return dbURL, nil
}

func NewMigrate(opts Options) (*Migrate, error) {
	mdir := opts.dir
	if mdir == "" {
		mdir = defaultMigrationsDir
	}
	_ = os.MkdirAll(defaultMigrationsDir, os.ModePerm)
	dbURL, err := databaseURLFromConfig()
	m, err := migrate.New("file://"+mdir, dbURL)
	if err != nil {
		return nil, err
	}
	return &Migrate{
		mig: m,
		dir: mdir,
	}, nil
}

func (m *Migrate) CreateMigrationFiles(opts CreateMigrationFilesOptions) error {
	migrationName := strcase.ToSnake(opts.Name)
	timestamp := time.Now().Format(defaultTimeFormat)
	filenameBase := fmt.Sprintf("%s_%s", timestamp, migrationName)
	filepathBase := path.Join(m.dir, filenameBase)
	filepathUp := fmt.Sprintf("%s.up.sql", filepathBase)
	filepathDown := fmt.Sprintf("%s.down.sql", filepathBase)
	fileUp, err := os.Create(filepathUp)
	if err != nil {
		return err
	}
	_, err = fileUp.WriteString(opts.UpSQL)
	if err != nil {
		return err
	}
	fileDown, err := os.Create(filepathDown)
	if err != nil {
		return err
	}
	_, err = fileDown.WriteString(opts.DownSQL)
	if err != nil {
		return err
	}
	fmt.Println("Migration created successfully")
	fmt.Printf("\t%s\t%s\n", color.GreenString("create"), filepathUp)
	fmt.Printf("\t%s\t%s\n", color.GreenString("create"), filepathDown)
	return nil
}

func (m *Migrate) Up() error {
	return m.mig.Up()
}

func (m *Migrate) UpSteps(steps int) error {
	return m.mig.Steps(steps)
}

func (m *Migrate) Down() error {
	return m.mig.Down()
}

func (m *Migrate) DownSteps(steps int) error {
	return m.mig.Steps(steps * -1)
}

func (m *Migrate) findFilenameFromVersion(version uint) (string, error) {
	var filename string
	err := filepath.Walk(m.dir, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), fmt.Sprintf("%d", version)) {
			filename = info.Name()
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return filename, nil
}

func (m *Migrate) Status() (*StatusResponse, error) {
	version, dirty, err := m.mig.Version()
	if err != nil {
		return nil, err
	}
	filename, err := m.findFilenameFromVersion(version)
	if err != nil {
		return nil, err
	}
	latestVersion, err := m.getLatestVersion()
	if err != nil {
		return nil, err
	}
	return &StatusResponse{
		CurrentVersion:     version,
		CurrentVersionFile: filename,
		LatestVersion:      latestVersion,
		Dirty:              dirty,
		Pending:            version != latestVersion,
	}, nil
}

func (m *Migrate) getLatestVersion() (uint, error) {
	latestVersion := uint(0)
	err := filepath.Walk(m.dir, func(path string, info os.FileInfo, err error) error {
		r := regexp.MustCompile(`(\d+)(\_)(.*)`)
		rs := r.FindStringSubmatch(info.Name())
		if len(rs) > 0 {
			u64, err := strconv.ParseUint(rs[1], 10, 64)
			if err != nil {
				return err
			}
			v := uint(u64)
			if v > latestVersion {
				latestVersion = v
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return latestVersion, nil
}
