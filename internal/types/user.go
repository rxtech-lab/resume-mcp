package types

import "context"

const (
	AuthenticatedUserContextKey = "authenticated_user"
)

type AuthenticatedUser struct {
	Aud      []string `json:"aud"`
	ClientId string   `json:"client_id"`
	Exp      int      `json:"exp"`
	Iat      int      `json:"iat"`
	Iss      string   `json:"iss"`
	Jti      string   `json:"jti"`
	Nbf      int      `json:"nbf"`
	Oid      string   `json:"oid"`
	Resid    string   `json:"resid"`
	Roles    []string `json:"roles"`
	Scopes   []string `json:"scopes"`
	Sid      string   `json:"sid"`
	Sub      string   `json:"sub"`
}

func WithAuthenticatedUser(ctx context.Context, user *AuthenticatedUser) context.Context {
	return context.WithValue(ctx, AuthenticatedUserContextKey, user)
}

func GetAuthenticatedUser(ctx context.Context) *AuthenticatedUser {
	return ctx.Value(AuthenticatedUserContextKey).(*AuthenticatedUser)
}
