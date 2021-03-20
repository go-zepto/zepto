package auth

import (
	"github.com/go-zepto/zepto/plugins/auth/authcore"
	"github.com/go-zepto/zepto/web"
)

type AuthTokenInstance interface {
	Core() *authcore.AuthCore
	LoggedPIDFromCtx(ctx web.Context) authcore.PID
	Auth(username string, password string) (*authcore.Token, error)
	Logout(authToken string) error
	PasswordRecovery(email string) error
	ResetPassword(resetPasswordToken string, password string) error
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

func (d *DefaultAuthTokenInstance) Auth(username string, password string) (*authcore.Token, error) {
	pid, err := d.Core().DS.FindPIDByValidCredentials(username, password)
	if err != nil {
		return nil, authcore.ErrUnauthorized
	}
	token, err := d.Core().TokenEncoder.GenerateTokenFromPID(pid)
	if err != nil {
		return nil, authcore.ErrInternalServerError
	}
	err = d.Core().Store.StoreAuthToken(token, pid)
	if err != nil {
		return nil, authcore.ErrInternalServerError
	}
	return token, nil
}

func (d *DefaultAuthTokenInstance) Logout(authToken string) error {
	return d.Core().Store.DeleteAuthToken(authToken)
}

func (d *DefaultAuthTokenInstance) PasswordRecovery(email string) error {
	pid, _ := d.Core().DS.FindPIDByEmail(email)
	if pid != nil {
		token, _ := d.Core().TokenEncoder.GenerateTokenFromPID(pid)
		d.Core().Store.StoreResetPasswordToken(token, pid)
		err := d.Core().Notifier.NotifyResetPasswordToken(email, token, pid)
		return err
	}
	return nil
}

func (d *DefaultAuthTokenInstance) ResetPassword(resetPasswordToken string, password string) error {
	pid, err := d.Core().Store.GetResetPasswordTokenPID(resetPasswordToken)
	if err != nil {
		return authcore.ErrInvalidToken
	}
	d.Core().DS.ResetPassword(pid, password)
	return d.Core().Store.DeleteResetPasswordToken(resetPasswordToken)
}

func InstanceFromCtx(ctx web.Context) AuthTokenInstance {
	i := ctx.PluginInstance("auth")
	AuthTokenInstance, ok := i.(AuthTokenInstance)
	if !ok {
		return nil
	}
	return AuthTokenInstance
}
