package auth

import (
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

type AuthTokenResponse struct {
	Token *Token `json:"token"`
}

type AuthTokenErrorResponse struct {
	Error string `json:"error"`
}

func NewAuthToken(opts AuthTokenOptions) *AuthToken {
	at := AuthToken{
		ds:           opts.Datasource,
		tokenEncoder: opts.TokenEncoder,
		store:        opts.Store,
	}
	return &at
}

func (at *AuthToken) middleware() web.MiddlewareFunc {
	return func(next web.RouteHandler) web.RouteHandler {
		return func(ctx web.Context) error {
			tokenValue := getTokenFromCtx(ctx)
			if tokenValue == "" {
				return next(ctx)
			}
			pid, err := at.store.GetAuthTokenPID(tokenValue)
			if err != nil {
				return next(ctx)
			}
			ctx.Set("auth_user_pid", pid)
			return next(ctx)
		}
	}
}

func (at *AuthToken) Create(z *zepto.Zepto) {
	z.Use(at.middleware())
}

func (at *AuthToken) SetupAuthEndpoint(z *zepto.Zepto, router *web.Router) {
	router.Post("/", func(ctx web.Context) error {
		credentials, err := getCredentialsFromCtx(ctx)
		if err != nil {
			ctx.SetStatus(400)
			return ctx.RenderJson(AuthTokenErrorResponse{
				Error: ErrMissingUsernameOrPassword.Error(),
			})
		}
		pid, err := at.ds.Auth(credentials.Username, credentials.Password)
		if err != nil {
			ctx.SetStatus(401)
			return ctx.RenderJson(AuthTokenErrorResponse{
				Error: ErrUnauthorized.Error(),
			})
		}
		token, err := at.tokenEncoder.GenerateTokenFromPID(pid)
		if err != nil {
			ctx.SetStatus(500)
			return ctx.RenderJson(AuthTokenErrorResponse{
				Error: ErrInternalServerError.Error(),
			})
		}
		err = at.store.StoreAuthToken(token, pid)
		if err != nil {
			ctx.SetStatus(500)
			return ctx.RenderJson(AuthTokenErrorResponse{
				Error: ErrInternalServerError.Error(),
			})
		}
		return ctx.RenderJson(AuthTokenResponse{
			Token: token,
		})
	})
}

func (at *AuthToken) Init(z *zepto.Zepto) {
	if at.ds == nil {
		panic("[auth] you must define a datasource")
	}
	if at.store == nil {
		panic("[auth] you must define a store")
	}
	router := z.Router("/auth")
	at.SetupAuthEndpoint(z, router)
}
