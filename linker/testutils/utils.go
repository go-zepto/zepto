package testutils

import (
	"testing"

	"github.com/go-zepto/zepto/linker"
	"github.com/go-zepto/zepto/linker/datasource/gorm/testutils"
	"github.com/go-zepto/zepto/web"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestKit struct {
	app    *web.App
	router *web.Router
	db     *gorm.DB
	linker *linker.Linker
}

func NewTestKit(t *testing.T) TestKit {
	r := require.New(t)
	app := web.NewApp()
	apiRouter := app.Router("/api")
	r.NotNil(apiRouter)
	db := testutils.SetupGorm()
	linker := linker.NewLinker(apiRouter)
	return TestKit{
		app:    app,
		router: apiRouter,
		db:     db,
		linker: linker,
	}
}
