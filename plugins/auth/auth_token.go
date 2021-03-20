package auth

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/plugins/auth/authcore"
	"github.com/go-zepto/zepto/web"
)

type AuthTokenOptions struct {
	Datasource     authcore.AuthDatasource
	TokenEncoder   authcore.AuthEncoder
	Store          authcore.AuthStore
	Notifier       authcore.AuthNotifier
	AuthContextKey string
}

type AuthToken struct {
	core     *authcore.AuthCore
	opts     *AuthTokenOptions
	instance AuthTokenInstance
}

func NewAuthTokenPlugin(opts AuthTokenOptions) *AuthToken {
	if opts.AuthContextKey == "" {
		opts.AuthContextKey = "auth_user_pid"
	}
	at := AuthToken{
		opts: &opts,
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
			pid, err := at.core.Store.GetAuthTokenPID(tokenValue)
			if err != nil {
				return next(ctx)
			}
			ctx.Set(at.opts.AuthContextKey, pid)
			return next(ctx)
		}
	}
}

func jsonError(ctx web.Context, err error) error {
	status := 401
	switch err {
	case authcore.ErrBadRequest, authcore.ErrMissingUsernameOrPassword, authcore.ErrInvalidToken:
		status = 400
	case authcore.ErrInternalServerError:
		status = 500
	}
	ctx.SetStatus(status)
	return ctx.RenderJson(authcore.AuthTokenErrorResponse{
		Error: err.Error(),
	})
}

func (at *AuthToken) setupAuthEndpoint(z *zepto.Zepto, router *web.Router) {
	router.Post("/", func(ctx web.Context) error {
		credentials, err := getCredentialsFromCtx(ctx)
		if err != nil {
			return jsonError(ctx, err)
		}
		token, err := at.instance.Auth(credentials.Username, credentials.Password)
		if err != nil {
			return jsonError(ctx, err)
		}
		return ctx.RenderJson(authcore.AuthTokenResponse{
			Token: token,
		})
	})
}

func (at *AuthToken) setupLogoutEndpoint(z *zepto.Zepto, router *web.Router) {
	router.Post("/logout", func(ctx web.Context) error {
		authToken := getTokenFromCtx(ctx)
		err := at.instance.Logout(authToken)
		if err != nil {
			jsonError(ctx, err)
		}
		return ctx.RenderJson(map[string]bool{
			"status": true,
		})
	})
}

func (at *AuthToken) setupRecoveryPasswordEndpoint(z *zepto.Zepto, router *web.Router) {
	router.Post("/recovery-password", func(ctx web.Context) error {
		req, err := getRecoveryPasswordRequestFromCtx(ctx)
		if err != nil {
			return jsonError(ctx, err)
		}
		err = at.instance.PasswordRecovery(req.Email)
		if err != nil {
			return jsonError(ctx, err)
		}
		return ctx.RenderJson(authcore.AuthRecoveryPasswordResponse{
			Status: true,
			Error:  nil,
		})
	})
}

func (at *AuthToken) setupResetPasswordEndpoint(z *zepto.Zepto, router *web.Router) {
	router.Post("/reset-password", func(ctx web.Context) error {
		req, err := getResetPasswordRequestFromCtx(ctx)
		if err != nil {
			ctx.SetStatus(400)
			return ctx.RenderJson(authcore.AuthTokenErrorResponse{
				Error: authcore.ErrBadRequest.Error(),
			})
		}
		err = at.instance.ResetPassword(req.Token, req.Password)
		if err != nil {
			return jsonError(ctx, err)
		}
		return ctx.RenderJson(map[string]bool{
			"status": true,
		})
	})
}

func (at *AuthToken) Name() string {
	return "auth"
}

func (at *AuthToken) Instance() interface{} {
	return at.instance
}

func (at *AuthToken) PrependMiddlewares() []web.MiddlewareFunc {
	return []web.MiddlewareFunc{at.middleware()}
}

func (at *AuthToken) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (at *AuthToken) OnCreated(z *zepto.Zepto) {
	at.core = CreateAuthCore(z, authcore.AuthCore{
		DS:           at.opts.Datasource,
		TokenEncoder: at.opts.TokenEncoder,
		Store:        at.opts.Store,
		Notifier:     at.opts.Notifier,
	})
	at.core.AssertConfigured()
	at.instance = &DefaultAuthTokenInstance{
		AuthCore:       at.core,
		AuthContextKey: "auth_user_pid",
	}
	return
}

func (at *AuthToken) OnZeptoInit(z *zepto.Zepto) {
	router := z.Router("/auth")
	at.setupAuthEndpoint(z, router)
	at.setupLogoutEndpoint(z, router)
	at.setupRecoveryPasswordEndpoint(z, router)
	at.setupResetPasswordEndpoint(z, router)
}

func (at *AuthToken) OnZeptoStart(z *zepto.Zepto) {
	return
}

func (at *AuthToken) OnZeptoStop(z *zepto.Zepto) {
	return
}
