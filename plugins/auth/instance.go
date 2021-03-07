package auth

import (
	"github.com/go-zepto/zepto/plugins/auth/authcore"
	"github.com/go-zepto/zepto/web"
)

type AuthTokenInstance interface {
	Core() *authcore.AuthCore
	LoggedPIDFromCtx(ctx web.Context) authcore.PID
}

type DefaultAuthTokenInstance struct {
	AuthCore       *authcore.AuthCore
	AuthContextKey string
}

func (d *DefaultAuthTokenInstance) Core() *authcore.AuthCore {
	return d.AuthCore
}

func (d *DefaultAuthTokenInstance) LoggedPIDFromCtx(ctx web.Context) authcore.PID {
	auth_user_pid := ctx.Value(d.AuthContextKey)
	if auth_user_pid == nil {
		return nil
	}
	return auth_user_pid.(authcore.PID)
}

func InstanceFromCtx(ctx web.Context) AuthTokenInstance {
	i := ctx.PluginInstance("auth")
	AuthTokenInstance, ok := i.(AuthTokenInstance)
	if !ok {
		return nil
	}
	return AuthTokenInstance
}
