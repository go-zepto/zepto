package auth

import (
	"github.com/go-zepto/zepto"
	"github.com/go-zepto/zepto/plugins/auth/authcore"
	"github.com/go-zepto/zepto/web"
	"go.uber.org/thriftrw/ptr"
)

type AuthTokenOptions struct {
	Datasource     authcore.AuthDatasource
	TokenEncoder   authcore.AuthEncoder
	Store          authcore.AuthStore
	Notifier       authcore.AuthNotifier
	AuthContextKey string
}

type AuthToken struct {
	core *authcore.AuthCore
	opts *AuthTokenOptions
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

func (at *AuthToken) setupAuthEndpoint(z *zepto.Zepto, router *web.Router) {
	router.Post("/", func(ctx web.Context) error {
		credentials, err := getCredentialsFromCtx(ctx)
		if err != nil {
			ctx.SetStatus(400)
			return ctx.RenderJson(authcore.AuthTokenErrorResponse{
				Error: authcore.ErrMissingUsernameOrPassword.Error(),
			})
		}
		pid, err := at.core.DS.Auth(credentials.Username, credentials.Password)
		if err != nil {
			ctx.SetStatus(401)
			return ctx.RenderJson(authcore.AuthTokenErrorResponse{
				Error: authcore.ErrUnauthorized.Error(),
			})
		}
		token, err := at.core.TokenEncoder.GenerateTokenFromPID(pid)
		if err != nil {
			ctx.SetStatus(500)
			return ctx.RenderJson(authcore.AuthTokenErrorResponse{
				Error: authcore.ErrInternalServerError.Error(),
			})
		}
		err = at.core.Store.StoreAuthToken(token, pid)
		if err != nil {
			ctx.SetStatus(500)
			return ctx.RenderJson(authcore.AuthTokenErrorResponse{
				Error: authcore.ErrInternalServerError.Error(),
			})
		}
		return ctx.RenderJson(authcore.AuthTokenResponse{
			Token: token,
		})
	})
}

func (at *AuthToken) setupRecoveryPasswordEndpoint(z *zepto.Zepto, router *web.Router) {
	router.Post("/recovery-password", func(ctx web.Context) error {
		req, err := getRecoveryPasswordRequestFromCtx(ctx)
		if err != nil {
			ctx.SetStatus(400)
			return ctx.RenderJson(authcore.AuthTokenErrorResponse{
				Error: authcore.ErrBadRequest.Error(),
			})
		}
		pid, err := at.core.DS.FindPIDByEmail(req.Email)
		if pid != nil {
			token, _ := at.core.TokenEncoder.GenerateTokenFromPID(pid)
			at.core.Store.StoreResetPasswordToken(token, pid)
			err := at.core.Notifier.NotifyResetPasswordToken(req.Email, token, pid)
			if err != nil {
				ctx.SetStatus(500)
				return ctx.RenderJson(authcore.AuthRecoveryPasswordResponse{
					Status: false,
					Error:  ptr.String(authcore.ErrInternalServerError.Error()),
				})
			}
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
		pid, err := at.core.Store.GetResetPasswordTokenPID(req.Token)
		at.core.DS.ResetPassword(pid, req.Password)
		if err != nil {
			ctx.SetStatus(400)
			return ctx.RenderJson(authcore.AuthResetPasswordResponse{
				Status: false,
				Error:  ptr.String(authcore.ErrInvalidToken.Error()),
			})
		}
		return ctx.RenderJson(authcore.AuthResetPasswordResponse{
			Status: true,
			Error:  nil,
		})
	})
}

func (at *AuthToken) Name() string {
	return "auth"
}

func (at *AuthToken) Instance() interface{} {
	return &DefaultAuthTokenInstance{
		AuthCore:       at.core,
		AuthContextKey: "auth_user_pid",
	}
}

func (at *AuthToken) PrependMiddlewares() []web.MiddlewareFunc {
	return []web.MiddlewareFunc{at.middleware()}
}

func (at *AuthToken) AppendMiddlewares() []web.MiddlewareFunc {
	return nil
}

func (at *AuthToken) OnCreated(z *zepto.Zepto) {
	return
}

func (at *AuthToken) OnZeptoInit(z *zepto.Zepto) {
	at.core = CreateAuthCore(z, authcore.AuthCore{
		DS:           at.opts.Datasource,
		TokenEncoder: at.opts.TokenEncoder,
		Store:        at.opts.Store,
		Notifier:     at.opts.Notifier,
	})
	at.core.AssertConfigured()
	router := z.Router("/auth")
	at.setupAuthEndpoint(z, router)
	at.setupRecoveryPasswordEndpoint(z, router)
	at.setupResetPasswordEndpoint(z, router)
}

func (at *AuthToken) OnZeptoStart(z *zepto.Zepto) {
	return
}

func (at *AuthToken) OnZeptoStop(z *zepto.Zepto) {
	return
}
