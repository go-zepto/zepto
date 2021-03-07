package auth_token_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zepto/zepto/plugins/auth/authcore"
	"github.com/go-zepto/zepto/plugins/auth/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func assertLogin(t *testing.T, kit *testutils.AuthTokenTestKit, username string, password string) *authcore.Token {
	w := httptest.NewRecorder()
	data := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, username, password)
	bodyReader := strings.NewReader(data)
	kit.Zepto.ServeHTTP(w, httptest.NewRequest("POST", "/auth", bodyReader))
	assert.Equal(t, 200, w.Result().StatusCode)
	var res authcore.AuthTokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.NotNil(t, res.Token)
	assert.NotZero(t, res.Token.Expiration)
	assert.NotZero(t, res.Token.Value)
	return res.Token
}

func TestAuthToken(t *testing.T) {
	var testAsserts = []struct {
		Username string
		Password string
		PID      int
	}{
		{
			Username: "carlos.strand",
			Password: "carlos-pwd-test-123",
			PID:      1,
		},
		{
			Username: "bill.gates",
			Password: "bill-pwd-test-123",
			PID:      2,
		},
	}
	for _, ta := range testAsserts {
		kit := testutils.NewAuthTokenTestKit()
		token := assertLogin(t, kit, ta.Username, ta.Password)
		assert.Equal(t, ta.PID, kit.GetMe(token.Value))
	}
}

func TestAuthToken_WrongCredentials(t *testing.T) {
	kit := testutils.NewAuthTokenTestKit()
	w := httptest.NewRecorder()
	bodyReader := strings.NewReader(`{"username": "unknown.user", "password": "wrong-password-123"}`)
	kit.Zepto.ServeHTTP(w, httptest.NewRequest("POST", "/auth", bodyReader))
	assert.Equal(t, 401, w.Result().StatusCode)
	var res authcore.AuthTokenErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "unauthorized", res.Error)
}

func generateUserResetToken(t *testing.T, kit *testutils.AuthTokenTestKit, email string) *authcore.Token {
	w := httptest.NewRecorder()
	bodyReader := strings.NewReader(fmt.Sprintf(`{"email": "%s"}`, email))
	kit.Zepto.ServeHTTP(w, httptest.NewRequest("POST", "/auth/recovery-password", bodyReader))
	assert.Equal(t, 200, w.Result().StatusCode)
	var res authcore.AuthResetPasswordResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	m := kit.MailerStub.LastSentFromTemplate()
	assert.NotNil(t, m)
	assert.NotNil(t, m.Opts.Vars["token"])
	return m.Opts.Vars["token"].(*authcore.Token)
}

func TestAuthToken_RecoveryPassword(t *testing.T) {
	kit := testutils.NewAuthTokenTestKit()
	generateUserResetToken(t, kit, "carlos.strand@go-zepto.com")
}

func TestAuthToken_ResetPassword(t *testing.T) {
	kit := testutils.NewAuthTokenTestKit()
	assertLogin(t, kit, "carlos.strand", "carlos-pwd-test-123")
	token := generateUserResetToken(t, kit, "carlos.strand@go-zepto.com")
	w := httptest.NewRecorder()
	data := fmt.Sprintf(`{"token": "%s", "password": "my-new-password"}`, token.Value)
	bodyReader := strings.NewReader(data)
	kit.Zepto.ServeHTTP(w, httptest.NewRequest("POST", "/auth/reset-password", bodyReader))
	assert.Equal(t, 200, w.Result().StatusCode)
	var res authcore.AuthResetPasswordResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Nil(t, res.Error)
	assertLogin(t, kit, "carlos.strand", "my-new-password")
}

func TestAuthToken_ResetPassword_InvalidToken(t *testing.T) {
	kit := testutils.NewAuthTokenTestKit()
	w := httptest.NewRecorder()
	bodyReader := strings.NewReader(`{"token": "some-invalid-token", "password": "my-new-password"}`)
	kit.Zepto.ServeHTTP(w, httptest.NewRequest("POST", "/auth/reset-password", bodyReader))
	assert.Equal(t, 400, w.Result().StatusCode)
	var res authcore.AuthResetPasswordResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Equal(t, "invalid token", *res.Error)
}
