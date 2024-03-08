package service

import (
	"context"
	"errors"
	"sawitpro/entity"
	"sawitpro/error_list"
	"sawitpro/helper"
	"sawitpro/mocks"
	"sawitpro/repository"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestNewProfileService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepository := mocks.NewMockUserProfileRepositoryInterface(ctrl)
	mockHelper := mocks.NewMockAuthHelperInterface(ctrl)

	type args struct {
		deps ProfileServiceDeps
	}
	tests := []struct {
		name string
		args args
		want profileService
	}{
		{
			name: "return profile service instance",
			args: args{
				deps: ProfileServiceDeps{
					ProfileRepository: mockProfileRepository,
					Authhelper:        mockHelper,
				},
			},
			want: profileService{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewProfileService(tt.args.deps)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_profileService_Register(t *testing.T) {
	mockTx := &sqlx.Tx{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepository := mocks.NewMockUserProfileRepositoryInterface(ctrl)
	mockHelper := mocks.NewMockAuthHelperInterface(ctrl)

	type fields struct {
		profileRepository repository.UserProfileRepositoryInterface
		authhelper        helper.AuthHelperInterface
	}
	type args struct {
		ctx     context.Context
		request entity.ProfileRegisterRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.ProfileRegisterResponse
		wantErr error
		mock    func()
	}{
		{
			name: "success register",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345",
				},
			},
			want: entity.ProfileRegisterResponse{
				Id: "profil-id-1",
			},
			wantErr: nil,
			mock: func() {
				mockHelper.EXPECT().HashPassword(gomock.Any(), "12345").Return("hashedPassword", nil)
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), mockTx, "+62345").Return(
					entity.UserProfile{}, nil,
				)
				mockProfileRepository.EXPECT().InsertProfile(gomock.Any(), mockTx, entity.UserProfile{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "hashedPassword",
				}).Return("profil-id-1", nil)
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
		{
			name: "error when insert profile",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345",
				},
			},
			want:    entity.ProfileRegisterResponse{},
			wantErr: errors.New("error when register a new profile"),
			mock: func() {
				mockHelper.EXPECT().HashPassword(gomock.Any(), "12345").Return("hashedPassword", nil)
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), mockTx, "+62345").Return(
					entity.UserProfile{}, nil,
				)
				mockProfileRepository.EXPECT().InsertProfile(gomock.Any(), mockTx, entity.UserProfile{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "hashedPassword",
				}).Return("", errors.New("error insert"))
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
		{
			name: "error duplicate phone number",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345",
				},
			},
			want:    entity.ProfileRegisterResponse{},
			wantErr: errors.New("error there existing data conficted with new data"),
			mock: func() {
				mockHelper.EXPECT().HashPassword(gomock.Any(), "12345").Return("hashedPassword", nil)
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), mockTx, "+62345").Return(
					entity.UserProfile{
						FullName:    "jonathan",
						PhoneNumber: "+62345",
						Password:    "12345",
					}, nil,
				)
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
		{
			name: "error when get profile",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345",
				},
			},
			want:    entity.ProfileRegisterResponse{},
			wantErr: errors.New("error when register a new profile"),
			mock: func() {
				mockHelper.EXPECT().HashPassword(gomock.Any(), "12345").Return("hashedPassword", nil)
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), mockTx, "+62345").Return(
					entity.UserProfile{}, errors.New("error get"),
				)
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
		{
			name: "error when hashing password",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345",
				},
			},
			want:    entity.ProfileRegisterResponse{},
			wantErr: errors.New("error when register a new profile"),
			mock: func() {
				mockHelper.EXPECT().HashPassword(gomock.Any(), "12345").Return("", errors.New("error hash password"))

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			p := profileService{
				profileRepository: tt.fields.profileRepository,
				authhelper:        tt.fields.authhelper,
			}
			got, err := p.Register(tt.args.ctx, tt.args.request)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_profileService_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepository := mocks.NewMockUserProfileRepositoryInterface(ctrl)
	mockHelper := mocks.NewMockAuthHelperInterface(ctrl)

	type fields struct {
		profileRepository repository.UserProfileRepositoryInterface
		authhelper        helper.AuthHelperInterface
	}
	type args struct {
		ctx     context.Context
		request entity.GetProfileRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.GetProfileResponse
		wantErr error
		mock    func()
	}{
		{
			name: "succes get profile",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.GetProfileRequest{
					ProfileId: "profile-id-1",
				},
			},
			want: entity.GetProfileResponse{
				FullName:    "jonathan",
				PhoneNumber: "+62345",
			},
			wantErr: nil,
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileById(gomock.Any(), nil, "profile-id-1").Return(
					entity.UserProfile{
						Id:          "profile-id-1",
						FullName:    "jonathan",
						PhoneNumber: "+62345",
						Password:    "12345",
					}, nil,
				)
			},
		},
		{
			name: "profile not found",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.GetProfileRequest{
					ProfileId: "profile-id-1",
				},
			},
			want:    entity.GetProfileResponse{},
			wantErr: errors.New("error profile not found"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileById(gomock.Any(), nil, "profile-id-1").Return(
					entity.UserProfile{}, nil,
				)
			},
		},
		{
			name: "error get profile from repository",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.GetProfileRequest{
					ProfileId: "profile-id-1",
				},
			},
			want:    entity.GetProfileResponse{},
			wantErr: errors.New("error when get user profile"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileById(gomock.Any(), nil, "profile-id-1").Return(
					entity.UserProfile{}, errors.New("error get"),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			p := profileService{
				profileRepository: tt.fields.profileRepository,
				authhelper:        tt.fields.authhelper,
			}
			got, err := p.GetProfile(tt.args.ctx, tt.args.request)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_profileService_Login(t *testing.T) {
	mockTx := &sqlx.Tx{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepository := mocks.NewMockUserProfileRepositoryInterface(ctrl)
	mockHelper := mocks.NewMockAuthHelperInterface(ctrl)

	type fields struct {
		profileRepository repository.UserProfileRepositoryInterface
		authhelper        helper.AuthHelperInterface
	}
	type args struct {
		ctx     context.Context
		request entity.LoginRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entity.LoginResponse
		wantErr error
		mock    func()
	}{
		{
			name: "success login",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345",
				},
			},
			want: entity.LoginResponse{
				Token: "token-1",
			},
			wantErr: nil,
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), nil, "+62345").Return(
					entity.UserProfile{
						Id:          "profile-id-1",
						FullName:    "jonathan",
						PhoneNumber: "+62345",
						Password:    "12345",
					}, nil,
				)
				mockHelper.EXPECT().VerifyPassword(gomock.Any(), "12345", "12345").Return(nil)
				mockHelper.EXPECT().GenerateToken(gomock.Any(), "profile-id-1").Return("token-1", nil)
				mockProfileRepository.EXPECT().IncreaseSuccessLoginCount(gomock.Any(), mockTx, "profile-id-1").Return(nil)
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
		{
			name: "error when increasing counter",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345",
				},
			},
			want:    entity.LoginResponse{},
			wantErr: errors.New("error when try to login"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), nil, "+62345").Return(
					entity.UserProfile{
						Id:          "profile-id-1",
						FullName:    "jonathan",
						PhoneNumber: "+62345",
						Password:    "12345",
					}, nil,
				)
				mockHelper.EXPECT().VerifyPassword(gomock.Any(), "12345", "12345").Return(nil)
				mockHelper.EXPECT().GenerateToken(gomock.Any(), "profile-id-1").Return("token-1", nil)
				mockProfileRepository.EXPECT().IncreaseSuccessLoginCount(gomock.Any(), mockTx, "profile-id-1").Return(errors.New("error update"))
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
		{
			name: "error when generate token",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345",
				},
			},
			want:    entity.LoginResponse{},
			wantErr: errors.New("error when try to login"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), nil, "+62345").Return(
					entity.UserProfile{
						Id:          "profile-id-1",
						FullName:    "jonathan",
						PhoneNumber: "+62345",
						Password:    "12345",
					}, nil,
				)
				mockHelper.EXPECT().VerifyPassword(gomock.Any(), "12345", "12345").Return(nil)
				mockHelper.EXPECT().GenerateToken(gomock.Any(), "profile-id-1").Return("", errors.New("error token"))
			},
		},
		{
			name: "error when password not match",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "123456",
				},
			},
			want:    entity.LoginResponse{},
			wantErr: errors.New("error credentials combination not match"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), nil, "+62345").Return(
					entity.UserProfile{
						Id:          "profile-id-1",
						FullName:    "jonathan",
						PhoneNumber: "+62345",
						Password:    "12345",
					}, nil,
				)
				mockHelper.EXPECT().VerifyPassword(gomock.Any(), "123456", "12345").Return(error_list.ErrPasswordNotMatch)
			},
		},
		{
			name: "error when verifying password",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "123456",
				},
			},
			want:    entity.LoginResponse{},
			wantErr: errors.New("error when try to login"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), nil, "+62345").Return(
					entity.UserProfile{
						Id:          "profile-id-1",
						FullName:    "jonathan",
						PhoneNumber: "+62345",
						Password:    "12345",
					}, nil,
				)
				mockHelper.EXPECT().VerifyPassword(gomock.Any(), "123456", "12345").Return(errors.New("error verify"))
			},
		},
		{
			name: "error when profile not found",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "123456",
				},
			},
			want:    entity.LoginResponse{},
			wantErr: errors.New("error credentials combination not match"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), nil, "+62345").Return(
					entity.UserProfile{}, nil,
				)
			},
		},
		{
			name: "error when get profile",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "123456",
				},
			},
			want:    entity.LoginResponse{},
			wantErr: errors.New("error when try to login"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), nil, "+62345").Return(
					entity.UserProfile{}, errors.New("error get"),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			p := profileService{
				profileRepository: tt.fields.profileRepository,
				authhelper:        tt.fields.authhelper,
			}
			got, err := p.Login(tt.args.ctx, tt.args.request)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func Test_profileService_UpdateProfile(t *testing.T) {

	mockTx := &sqlx.Tx{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileRepository := mocks.NewMockUserProfileRepositoryInterface(ctrl)
	mockHelper := mocks.NewMockAuthHelperInterface(ctrl)

	type fields struct {
		profileRepository repository.UserProfileRepositoryInterface
		authhelper        helper.AuthHelperInterface
	}
	type args struct {
		ctx     context.Context
		request entity.UpdateProfileRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		mock    func()
	}{
		{
			name: "success update profile",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.UpdateProfileRequest{
					Id:          "profile-id-1",
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				},
			},
			wantErr: nil,
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), mockTx, "+62345").Return(
					entity.UserProfile{}, nil,
				)
				mockProfileRepository.EXPECT().UpdateProfileById(gomock.Any(), mockTx, "profile-id-1", entity.UserProfile{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				}).Return(nil)
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
		{
			name: "error when update the profile",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.UpdateProfileRequest{
					Id:          "profile-id-1",
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				},
			},
			wantErr: errors.New("error when updating profile"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), mockTx, "+62345").Return(
					entity.UserProfile{}, nil,
				)
				mockProfileRepository.EXPECT().UpdateProfileById(gomock.Any(), mockTx, "profile-id-1", entity.UserProfile{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				}).Return(errors.New("error update"))
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
		{
			name: "error duplicate phone number",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.UpdateProfileRequest{
					Id:          "profile-id-1",
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				},
			},
			wantErr: errors.New("error there existing data conficted with new data"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), mockTx, "+62345").Return(
					entity.UserProfile{
						Id:          "profile-id-1",
						FullName:    "jonathan",
						PhoneNumber: "+62345",
						Password:    "12345",
					}, nil,
				)
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
		{
			name: "error when get existing profile",
			fields: fields{
				profileRepository: mockProfileRepository,
				authhelper:        mockHelper,
			},
			args: args{
				ctx: context.TODO(),
				request: entity.UpdateProfileRequest{
					Id:          "profile-id-1",
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				},
			},
			wantErr: errors.New("error when updating profile"),
			mock: func() {
				mockProfileRepository.EXPECT().GetProfileByPhoneNumber(gomock.Any(), mockTx, "+62345").Return(
					entity.UserProfile{}, errors.New("error select"),
				)
				mockProfileRepository.EXPECT().RunWithTransaction(gomock.Any(), gomock.Any()).DoAndReturn(
					func(ctx context.Context, handleFunc func(tx *sqlx.Tx) error) interface{} {
						return handleFunc(mockTx)
					},
				)
			},
		},
	}
	for _, tt := range tests {
		tt.mock()

		t.Run(tt.name, func(t *testing.T) {
			p := profileService{
				profileRepository: tt.fields.profileRepository,
				authhelper:        tt.fields.authhelper,
			}
			err := p.UpdateProfile(tt.args.ctx, tt.args.request)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
