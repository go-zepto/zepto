package auth

import (
	"fmt"

	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/web"
)

type AuthToken struct {
	ds           AuthDatasource
	tokenEncoder AuthEncoder
	store        AuthStore
}

type AuthTokenOptions struct {
	Datasource   AuthDatasource
	TokenEncoder AuthEncoder
	Store        AuthStore
}

func NewAuthToken(opts AuthTokenOptions) *AuthToken {
	at := AuthToken{
		ds:           opts.Datasource,
		tokenEncoder: opts.TokenEncoder,
		store:        opts.Store,
	}
	return &at
}

func (at *AuthToken) Middleware() web.MiddlewareFunc {
	return func(next web.RouteHandler) web.RouteHandler {
		return func(ctx web.Context) error {
			tokenValue := getTokenFromCtx(ctx)
			if tokenValue == "" {
				return next(ctx)
			}
			pid, err := at.store.GetAuthTokenPID(tokenValue)
			fmt.Println(pid)
			if err != nil {
				return next(ctx)
			}
			ctx.Set("auth_user_pid", pid)
			return next(ctx)
		}
	}
}

func (at *AuthToken) Create(z *zepto.Zepto) {
	z.Use(at.Middleware())
}

func (at *AuthToken) Init(z *zepto.Zepto) {
	if at.ds == nil {
		panic("[auth] you must define a datasource")
	}
	if at.store == nil {
		panic("[auth] you must define a store")
	}
	authRouter := z.Router("/auth")
	authRouter.Post("/", func(ctx web.Context) error {
		credentials, err := getCredentialsFromCtx(ctx)
		if err != nil {
			ctx.SetStatus(400)
			return ctx.RenderJson(map[string]interface{}{
				"error": "missing username or password",
			})
		}
		pid, err := at.ds.Auth(credentials.Username, credentials.Password)
		if err != nil {
			ctx.SetStatus(401)
			return ctx.RenderJson(map[string]interface{}{
				"error": "unauthorized",
			})
		}
		token, err := at.tokenEncoder.GenerateTokenFromPID(pid)
		if err != nil {
			ctx.SetStatus(401)
			return ctx.RenderJson(map[string]interface{}{
				"error": "internal server error",
			})
		}
		err = at.store.StoreAuthToken(token, pid)
		if err != nil {
			ctx.SetStatus(500)
			return ctx.RenderJson(map[string]interface{}{
				"error": err.Error(),
			})
		}
		return ctx.RenderJson(map[string]interface{}{
			"token": token,
		})
	})
}
