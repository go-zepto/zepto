package auth_token_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-zepto/zepto/plugins/auth"
	"github.com/go-zepto/zepto/plugins/auth/tests/testutils"
	"github.com/stretchr/testify/assert"
)

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
		w := httptest.NewRecorder()
		data := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, ta.Username, ta.Password)
		bodyReader := strings.NewReader(data)
		kit.Zepto.ServeHTTP(w, httptest.NewRequest("POST", "/auth", bodyReader))
		assert.Equal(t, 200, w.Result().StatusCode)
		var res auth.AuthTokenResponse
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.NotNil(t, res.Token)
		assert.NotZero(t, res.Token.Expiration)
		assert.NotZero(t, res.Token.Value)
		assert.Equal(t, ta.PID, kit.GetMe(res.Token.Value))
	}
}

func TestAuthToken_WrongCredentials(t *testing.T) {
	kit := testutils.NewAuthTokenTestKit()
	w := httptest.NewRecorder()
	bodyReader := strings.NewReader(`{"username": "unknown.user", "password": "wrong-password-123"}`)
	kit.Zepto.ServeHTTP(w, httptest.NewRequest("POST", "/auth", bodyReader))
	assert.Equal(t, 401, w.Result().StatusCode)
	var res auth.AuthTokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &res)
	assert.NoError(t, err)
	assert.Nil(t, res.Token)
}
