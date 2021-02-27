package testutils

import (
	"encoding/json"
	"net/http/httptest"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/plugins/auth"
	gormds "github.com/go-zepto/zepto/plugins/auth/datasources/gorm"
	"github.com/go-zepto/zepto/plugins/auth/datasources/gorm/testutils"
	"github.com/go-zepto/zepto/plugins/auth/datasources/gorm/testutils/models"
	"github.com/go-zepto/zepto/plugins/auth/encoders/uuid"
	"github.com/go-zepto/zepto/plugins/auth/stores/inmemory"
	"github.com/go-zepto/zepto/web"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthTokenTestKit struct {
	Zepto     *zepto.Zepto
	DB        *gorm.DB
	AuthToken *auth.AuthToken
}

func generatePasswordHash(pwd string) string {
	res, _ := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	return string(res)
}

func createUsers(db *gorm.DB) {
	db.Create([]models.User{
		{
			Username:     "carlos.strand",
			PasswordHash: generatePasswordHash("carlos-pwd-test-123"),
		},
		{
			Username:     "bill.gates",
			PasswordHash: generatePasswordHash("bill-pwd-test-123"),
		},
	})
}

func NewAuthTokenTestKit() *AuthTokenTestKit {
	z := zepto.NewZepto(
		zepto.Name("auth-token-test"),
	)
	db := testutils.SetupDB()
	authToken := auth.NewAuthToken(auth.AuthTokenOptions{
		Datasource: gormds.NewGormAuthDatasoruce(gormds.GormAuthDatasourceOptions{
			DB:        db,
			UserModel: &models.User{},
		}),
		TokenEncoder: uuid.NewUUIDTokenEncoder(),
		Store:        inmemory.NewInMemoryStore(),
	})
	authToken.Create(z)
	authToken.Init(z)
	z.Get("/me", func(ctx web.Context) error {
		return ctx.RenderJson(map[string]interface{}{
			"pid": ctx.Value("auth_user_pid"),
		})
	})
	z.InitApp()
	createUsers(db)
	return &AuthTokenTestKit{
		Zepto:     z,
		DB:        db,
		AuthToken: authToken,
	}
}

func (k *AuthTokenTestKit) GetMe(tokenValue string) int {
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/me", nil)
	req.Header.Add("Authorization", "Bearer "+tokenValue)
	k.Zepto.ServeHTTP(w, req)
	var res struct {
		PID int `json:"pid"`
	}
	json.Unmarshal(w.Body.Bytes(), &res)
	return res.PID
}
