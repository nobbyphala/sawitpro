package service

import (
	"context"
	"sawitpro/entity"
)

type ProfileServiceInterface interface {
	Register(ctx context.Context, request entity.ProfileRegisterRequest) (entity.ProfileRegisterResponse, error)
	Login(ctx context.Context, request entity.LoginRequest) (entity.LoginResponse, error)
	UpdateProfile(ctx context.Context, request entity.UpdateProfileRequest) error
	GetProfile(ctx context.Context, request entity.GetProfileRequest) (entity.GetProfileResponse, error)
}
