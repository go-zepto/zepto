package auth

import (
	"encoding/json"
	"strings"

	"github.com/go-zepto/zepto/plugins/auth/authcore"
	"github.com/go-zepto/zepto/web"
)

func getCredentialsFromCtx(ctx web.Context) (*authcore.AuthCredentials, error) {
	credentials := authcore.AuthCredentials{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&credentials)
	if err != nil {
		return nil, err
	}
	return &credentials, nil
}

func getResetPasswordRequestFromCtx(ctx web.Context) (*authcore.AuthResetPasswordRequest, error) {
	req := authcore.AuthResetPasswordRequest{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	if req.Password == "" || req.Token == "" {
		return nil, authcore.ErrBadRequest
	}
	return &req, nil
}

func getRecoveryPasswordRequestFromCtx(ctx web.Context) (*authcore.AuthRecoveryPasswordRequest, error) {
	req := authcore.AuthRecoveryPasswordRequest{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	if req.Email == "" {
		return nil, authcore.ErrBadRequest
	}
	return &req, nil
}

func getTokenFromCtx(ctx web.Context) string {
	authStr := ctx.Request().Header.Get("Authorization")
	if !strings.HasPrefix(authStr, "Bearer ") {
		return ""
	}
	return authStr[7:]
}
