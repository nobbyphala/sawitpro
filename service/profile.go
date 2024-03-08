package service

import (
	"context"
	"sawitpro/entity"
	"sawitpro/error_list"
	"sawitpro/helper"
	"sawitpro/repository"

	"github.com/jmoiron/sqlx"
)

type profileService struct {
	profileRepository repository.UserProfileRepositoryInterface
	authhelper        helper.AuthHelperInterface
}

type ProfileServiceDeps struct {
	ProfileRepository repository.UserProfileRepositoryInterface
	Authhelper        helper.AuthHelperInterface
}

func NewProfileService(deps ProfileServiceDeps) profileService {
	return profileService{
		profileRepository: deps.ProfileRepository,
		authhelper:        deps.Authhelper,
	}
}

func (p profileService) Register(ctx context.Context, request entity.ProfileRegisterRequest) (entity.ProfileRegisterResponse, error) {
	var res = entity.ProfileRegisterResponse{}

	hashedPassword, err := p.authhelper.HashPassword(ctx, request.Password)
	if err != nil {
		return res, error_list.ErrProfileRegister
	}

	var profileId string

	err = p.profileRepository.RunWithTransaction(ctx, func(tx *sqlx.Tx) error {
		existingProfile, err := p.profileRepository.GetProfileByPhoneNumber(ctx, tx, request.PhoneNumber)
		if err != nil {
			return error_list.ErrProfileRegister
		}

		if existingProfile.PhoneNumber == request.PhoneNumber {
			return error_list.ErrDataConflict
		}

		profileId, err = p.profileRepository.InsertProfile(ctx, tx, entity.UserProfile{
			FullName:    request.FullName,
			PhoneNumber: request.PhoneNumber,
			Password:    hashedPassword,
		})
		if err != nil {
			return error_list.ErrProfileRegister
		}

		return nil
	})
	if err != nil {
		return res, err
	}

	res = entity.ProfileRegisterResponse{
		Id: profileId,
	}

	return res, nil
}

func (p profileService) GetProfile(ctx context.Context, request entity.GetProfileRequest) (entity.GetProfileResponse, error) {
	var res = entity.GetProfileResponse{}
	profileId := request.ProfileId

	profile, err := p.profileRepository.GetProfileById(ctx, nil, profileId)
	if err != nil {
		return res, error_list.ErrGetProfile
	}

	if profile.Id == "" {
		return res, error_list.ErrProfileNotFound
	}

	res = entity.GetProfileResponse{
		FullName:    profile.FullName,
		PhoneNumber: profile.PhoneNumber,
	}

	return res, nil
}

func (p profileService) Login(ctx context.Context, request entity.LoginRequest) (entity.LoginResponse, error) {
	var res = entity.LoginResponse{}

	profile, err := p.profileRepository.GetProfileByPhoneNumber(ctx, nil, request.PhoneNumber)
	if err != nil {
		return res, error_list.ErrLogin
	}

	if profile.Id == "" {
		return res, error_list.ErrLoginCredential
	}

	err = p.authhelper.VerifyPassword(ctx, request.Password, profile.Password)
	if err != nil {
		if err == error_list.ErrPasswordNotMatch {
			return res, error_list.ErrLoginCredential
		}
		return res, error_list.ErrLogin
	}

	token, err := p.authhelper.GenerateToken(ctx, profile.Id)
	if err != nil {
		return res, error_list.ErrLogin
	}

	err = p.profileRepository.RunWithTransaction(ctx, func(tx *sqlx.Tx) error {
		// assumption the counter must be sync
		err := p.profileRepository.IncreaseSuccessLoginCount(ctx, tx, profile.Id)
		if err != nil {
			return error_list.ErrLogin
		}

		return nil
	})
	if err != nil {
		return res, err
	}

	res = entity.LoginResponse{
		Token: token,
	}

	return res, nil
}

func (p profileService) UpdateProfile(ctx context.Context, request entity.UpdateProfileRequest) error {
	err := p.profileRepository.RunWithTransaction(ctx, func(tx *sqlx.Tx) error {
		existingProfile, err := p.profileRepository.GetProfileByPhoneNumber(ctx, tx, request.PhoneNumber)
		if err != nil {
			return error_list.ErrUpdateProfile
		}

		if existingProfile.PhoneNumber == request.PhoneNumber {
			return error_list.ErrDataConflict
		}

		err = p.profileRepository.UpdateProfileById(ctx, tx, request.Id, entity.UserProfile{
			FullName:    request.FullName,
			PhoneNumber: request.PhoneNumber,
		})
		if err != nil {
			return error_list.ErrUpdateProfile
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
