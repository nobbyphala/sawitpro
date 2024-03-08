package helper

import (
	"context"
	"sawitpro/constant"
	"sawitpro/error_list"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authHelper struct {
}

func NewAuthHelper() authHelper {
	return authHelper{}
}

func (hlp authHelper) HashPassword(ctx context.Context, password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashed), err
}

func (hlp authHelper) VerifyPassword(ctx context.Context, plainPassword string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return error_list.ErrPasswordNotMatch
		}

		return err
	}

	return nil
}

func (hlp authHelper) GenerateToken(ctx context.Context, profileId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		constant.ProfileIdJwtField: profileId,
	})

	return token.SignedString([]byte(constant.EnvJWTSecretKey))
}

func (hlp authHelper) VerifyToken(ctx context.Context, token string) (string, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(constant.EnvJWTSecretKey), nil
	})
	if err != nil {
		return "", error_list.ErrInvalidToken
	}

	claims, claimsExist := jwtToken.Claims.(jwt.MapClaims)
	if !claimsExist {
		return "", error_list.ErrInvalidToken
	}

	profileId, profileIdExists := claims[constant.ProfileIdJwtField]
	if !profileIdExists {
		return "", error_list.ErrInvalidToken
	}

	profileIdStr, ok := profileId.(string)
	if !ok {
		return "", error_list.ErrInvalidToken
	}

	return profileIdStr, nil
}
