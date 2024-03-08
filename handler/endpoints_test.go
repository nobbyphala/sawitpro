package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sawitpro/entity"
	"sawitpro/generated"
	"sawitpro/helper"
	"sawitpro/mocks"
	"sawitpro/service"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestServer_RegisterProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileService := mocks.NewMockProfileServiceInterface(ctrl)
	mockAuthHelper := mocks.NewMockAuthHelperInterface(ctrl)
	mockValidatorHelper := mocks.NewMockValidatorHelperInterface(ctrl)

	type fields struct {
		profileService  service.ProfileServiceInterface
		authHelper      helper.AuthHelperInterface
		validatorHelper helper.ValidatorHelperInterface
	}
	type args struct {
		req generated.RegisterProfileRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       generated.RegisterProfileResponse
		wantErr    bool
		errResp    *generated.ErrorResponse
		statusCode int
		mock       func()
	}{
		{
			name: "success register",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.RegisterProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345A!",
				},
			},
			want: generated.RegisterProfileResponse{
				ProfileId: "profile-id-1",
			},
			wantErr:    false,
			errResp:    nil,
			statusCode: http.StatusOK,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(nil)
				mockProfileService.EXPECT().Register(gomock.Any(), entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(entity.ProfileRegisterResponse{
					Id: "profile-id-1",
				}, nil)
			},
		},
		{
			name: "error when register",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.RegisterProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345A!",
				},
			},
			want:    generated.RegisterProfileResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error when register a new profile",
			},
			statusCode: http.StatusInternalServerError,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(nil)
				mockProfileService.EXPECT().Register(gomock.Any(), entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(entity.ProfileRegisterResponse{}, errors.New("error when register a new profile"))
			},
		},
		{
			name: "error when duplicated data",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.RegisterProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345A!",
				},
			},
			want:    generated.RegisterProfileResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error there existing data conficted with new data",
			},
			statusCode: http.StatusConflict,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(nil)
				mockProfileService.EXPECT().Register(gomock.Any(), entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(entity.ProfileRegisterResponse{}, errors.New("error there existing data conficted with new data"))
			},
		},
		{
			name: "error invalid payload",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.RegisterProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "-----",
				},
			},
			want:    generated.RegisterProfileResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "invalid payload at password",
			},
			statusCode: http.StatusBadRequest,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.ProfileRegisterRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Password:    "-----",
				}).Return(errors.New("invalid payload at password"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := &Server{
				profileService:  tt.fields.profileService,
				authHelper:      tt.fields.authHelper,
				validatorHelper: tt.fields.validatorHelper,
			}

			e := echo.New()

			e.POST("/register", s.RegisterProfile)

			requestBody, _ := json.Marshal(tt.args.req)

			req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var expectBody []byte

			if tt.wantErr {
				expectBody, _ = json.Marshal(tt.errResp)
			} else {
				expectBody, _ = json.Marshal(tt.want)
			}

			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, strings.TrimSpace(string(expectBody)), strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestServer_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileService := mocks.NewMockProfileServiceInterface(ctrl)
	mockAuthHelper := mocks.NewMockAuthHelperInterface(ctrl)
	mockValidatorHelper := mocks.NewMockValidatorHelperInterface(ctrl)

	type fields struct {
		profileService  service.ProfileServiceInterface
		authHelper      helper.AuthHelperInterface
		validatorHelper helper.ValidatorHelperInterface
	}
	type args struct {
		req generated.LoginRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       generated.LoginResponse
		wantErr    bool
		errResp    *generated.ErrorResponse
		statusCode int
		mock       func()
	}{
		{
			name: "success login",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				},
			},
			want: generated.LoginResponse{
				Token: "token1",
			},
			wantErr:    false,
			errResp:    nil,
			statusCode: http.StatusOK,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(nil)
				mockProfileService.EXPECT().Login(gomock.Any(), entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(entity.LoginResponse{
					Token: "token1",
				}, nil)
			},
		},
		{
			name: "error profile not found",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				},
			},
			want:    generated.LoginResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error profile not found",
			},
			statusCode: http.StatusNotFound,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(nil)
				mockProfileService.EXPECT().Login(gomock.Any(), entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(entity.LoginResponse{}, errors.New("error profile not found"))
			},
		},
		{
			name: "error wrong credentials",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				},
			},
			want:    generated.LoginResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error credentials combination not match",
			},
			statusCode: http.StatusBadRequest,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(nil)
				mockProfileService.EXPECT().Login(gomock.Any(), entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(entity.LoginResponse{}, errors.New("error credentials combination not match"))
			},
		},
		{
			name: "error when try to login",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				},
			},
			want:    generated.LoginResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error when try to login",
			},
			statusCode: http.StatusInternalServerError,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(nil)
				mockProfileService.EXPECT().Login(gomock.Any(), entity.LoginRequest{
					PhoneNumber: "+62345",
					Password:    "12345A!",
				}).Return(entity.LoginResponse{}, errors.New("error when try to login"))
			},
		},
		{
			name: "error request not valid",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.LoginRequest{
					PhoneNumber: "62345",
					Password:    "12345",
				},
			},
			want:    generated.LoginResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error phone number not valid",
			},
			statusCode: http.StatusBadRequest,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.LoginRequest{
					PhoneNumber: "62345",
					Password:    "12345",
				}).Return(errors.New("error phone number not valid"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := &Server{
				profileService:  tt.fields.profileService,
				authHelper:      tt.fields.authHelper,
				validatorHelper: tt.fields.validatorHelper,
			}

			e := echo.New()

			e.POST("/login", s.Login)

			requestBody, _ := json.Marshal(tt.args.req)

			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var expectBody []byte

			if tt.wantErr {
				expectBody, _ = json.Marshal(tt.errResp)
			} else {
				expectBody, _ = json.Marshal(tt.want)
			}

			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, strings.TrimSpace(string(expectBody)), strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestServer_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileService := mocks.NewMockProfileServiceInterface(ctrl)
	mockAuthHelper := mocks.NewMockAuthHelperInterface(ctrl)
	mockValidatorHelper := mocks.NewMockValidatorHelperInterface(ctrl)

	type fields struct {
		profileService  service.ProfileServiceInterface
		authHelper      helper.AuthHelperInterface
		validatorHelper helper.ValidatorHelperInterface
	}
	type args struct {
		profileId string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       generated.GetProfileResponse
		wantErr    bool
		errResp    *generated.ErrorResponse
		statusCode int
		mock       func()
	}{
		{
			name: "success get profile",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				profileId: "profile-id-1",
			},
			want: generated.GetProfileResponse{
				FullName:    "jonathan",
				PhoneNumber: "+62345",
			},
			wantErr:    false,
			errResp:    nil,
			statusCode: http.StatusOK,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.GetProfileRequest{
					ProfileId: "profile-id-1",
				}).Return(nil)
				mockProfileService.EXPECT().GetProfile(gomock.Any(), entity.GetProfileRequest{
					ProfileId: "profile-id-1",
				}).Return(entity.GetProfileResponse{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				}, nil)
			},
		},
		{
			name: "error when get profile",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				profileId: "profile-id-1",
			},
			want:    generated.GetProfileResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error when get user profile",
			},
			statusCode: http.StatusInternalServerError,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.GetProfileRequest{
					ProfileId: "profile-id-1",
				}).Return(nil)
				mockProfileService.EXPECT().GetProfile(gomock.Any(), entity.GetProfileRequest{
					ProfileId: "profile-id-1",
				}).Return(entity.GetProfileResponse{}, errors.New("error when get user profile"))
			},
		},
		{
			name: "error when validate request",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				profileId: "profile-id-1",
			},
			want:    generated.GetProfileResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error invalid profile id",
			},
			statusCode: http.StatusBadRequest,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.GetProfileRequest{
					ProfileId: "profile-id-1",
				}).Return(errors.New("error invalid profile id"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := &Server{
				profileService:  tt.fields.profileService,
				authHelper:      tt.fields.authHelper,
				validatorHelper: tt.fields.validatorHelper,
			}

			wrapper := func(ctx echo.Context) error {
				ctx.Set("profile_id", tt.args.profileId)

				return s.GetProfile(ctx, generated.GetProfileParams{})
			}

			e := echo.New()

			e.GET("/profile", wrapper)

			req := httptest.NewRequest(http.MethodGet, "/profile", nil)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var expectBody []byte

			if tt.wantErr {
				expectBody, _ = json.Marshal(tt.errResp)
			} else {
				expectBody, _ = json.Marshal(tt.want)
			}

			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, strings.TrimSpace(string(expectBody)), strings.TrimSpace(rec.Body.String()))
		})
	}
}

func TestServer_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProfileService := mocks.NewMockProfileServiceInterface(ctrl)
	mockAuthHelper := mocks.NewMockAuthHelperInterface(ctrl)
	mockValidatorHelper := mocks.NewMockValidatorHelperInterface(ctrl)

	type fields struct {
		profileService  service.ProfileServiceInterface
		authHelper      helper.AuthHelperInterface
		validatorHelper helper.ValidatorHelperInterface
	}
	type args struct {
		req       generated.UpdateProfileRequest
		profileId string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		want       generated.UpdateProfileResponse
		wantErr    bool
		errResp    *generated.ErrorResponse
		statusCode int
		mock       func()
	}{
		{
			name: "success update profile",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				},
				profileId: "profile-id-1",
			},
			want: generated.UpdateProfileResponse{
				Message: "Success update profile",
			},
			wantErr:    false,
			errResp:    nil,
			statusCode: http.StatusOK,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Id:          "profile-id-1",
				}).Return(nil)
				mockProfileService.EXPECT().UpdateProfile(gomock.Any(), entity.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Id:          "profile-id-1",
				}).Return(nil)
			},
		},
		{
			name: "error data conflict",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				},
				profileId: "profile-id-1",
			},
			want:    generated.UpdateProfileResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error there existing data conficted with new data",
			},
			statusCode: http.StatusConflict,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Id:          "profile-id-1",
				}).Return(nil)
				mockProfileService.EXPECT().UpdateProfile(gomock.Any(), entity.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Id:          "profile-id-1",
				}).Return(errors.New("error there existing data conficted with new data"))
			},
		},
		{
			name: "error when update profile",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
				},
				profileId: "profile-id-1",
			},
			want:    generated.UpdateProfileResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error when updating profile",
			},
			statusCode: http.StatusInternalServerError,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Id:          "profile-id-1",
				}).Return(nil)
				mockProfileService.EXPECT().UpdateProfile(gomock.Any(), entity.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+62345",
					Id:          "profile-id-1",
				}).Return(errors.New("error when updating profile"))
			},
		},
		{
			name: "error when validate",
			fields: fields{
				profileService:  mockProfileService,
				authHelper:      mockAuthHelper,
				validatorHelper: mockValidatorHelper,
			},
			args: args{
				req: generated.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+623s45",
				},
				profileId: "profile-id-1",
			},
			want:    generated.UpdateProfileResponse{},
			wantErr: true,
			errResp: &generated.ErrorResponse{
				Message: "error in phone number",
			},
			statusCode: http.StatusBadRequest,
			mock: func() {
				mockValidatorHelper.EXPECT().ValidateStruct(entity.UpdateProfileRequest{
					FullName:    "jonathan",
					PhoneNumber: "+623s45",
					Id:          "profile-id-1",
				}).Return(errors.New("error in phone number"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			s := &Server{
				profileService:  tt.fields.profileService,
				authHelper:      tt.fields.authHelper,
				validatorHelper: tt.fields.validatorHelper,
			}

			e := echo.New()

			wrapper := func(ctx echo.Context) error {
				ctx.Set("profile_id", tt.args.profileId)

				return s.UpdateProfile(ctx, generated.UpdateProfileParams{})
			}

			e.PUT("/register", wrapper)

			requestBody, _ := json.Marshal(tt.args.req)

			req := httptest.NewRequest(http.MethodPut, "/register", strings.NewReader(string(requestBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			var expectBody []byte

			if tt.wantErr {
				expectBody, _ = json.Marshal(tt.errResp)
			} else {
				expectBody, _ = json.Marshal(tt.want)
			}

			assert.Equal(t, tt.statusCode, rec.Code)
			assert.Equal(t, strings.TrimSpace(string(expectBody)), strings.TrimSpace(rec.Body.String()))
		})
	}
}
