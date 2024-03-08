package helper

import "context"

type AuthHelperInterface interface {
	HashPassword(ctx context.Context, password string) (string, error)
	VerifyPassword(ctx context.Context, plainPassword string, hashedPassword string) error
	GenerateToken(ctx context.Context, profileId string) (string, error)
	VerifyToken(ctx context.Context, token string) (string, error)
}

type ValidatorHelperInterface interface {
	ValidateStruct(s interface{}) error
}
