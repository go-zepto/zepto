package auth

import (
	"encoding/json"
	"strings"

	"github.com/go-zepto/zepto/web"
)

func getCredentialsFromCtx(ctx web.Context) (*AuthCredentials, error) {
	credentials := AuthCredentials{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&credentials)
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

func getTokenFromCtx(ctx web.Context) string {
	authStr := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authStr, "Bearer ") {
		return ""
	}
	return authStr[7:]
}
