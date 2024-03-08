package handler

import (
	"net/http"

	"sawitpro/constant"
	"sawitpro/entity"
	"sawitpro/error_list"
	"sawitpro/generated"

	"github.com/labstack/echo/v4"
)

func (s *Server) RegisterProfile(ctx echo.Context) error {
	var req generated.RegisterProfileRequest

	err := ctx.Bind(&req)
	if err != nil {
		return s.sendErrorResponse(ctx, err)
	}

	registerReq := entity.ProfileRegisterRequest{
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}

	err = s.validate(registerReq)
	if err != nil {
		return s.sendValidationErrorResponse(ctx, err)
	}

	result, err := s.profileService.Register(ctx.Request().Context(), registerReq)
	if err != nil {
		return s.sendErrorResponse(ctx, err)
	}

	resp := generated.RegisterProfileResponse{
		ProfileId: result.Id,
	}
	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) Login(ctx echo.Context) error {
	var req generated.LoginRequest

	err := ctx.Bind(&req)
	if err != nil {
		return s.sendErrorResponse(ctx, error_list.ErrInvalidRequest)
	}

	loginReq := entity.LoginRequest{
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
	}
	err = s.validate(loginReq)
	if err != nil {
		return s.sendValidationErrorResponse(ctx, err)
	}

	result, err := s.profileService.Login(ctx.Request().Context(), loginReq)
	if err != nil {
		return s.sendErrorResponse(ctx, err)
	}

	resp := generated.LoginResponse{
		Token: result.Token,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) GetProfile(ctx echo.Context, params generated.GetProfileParams) error {
	profileId, ok := ctx.Get(constant.ProfileIdJwtField).(string)
	if !ok {
		return s.sendErrorResponse(ctx, error_list.ErrInvalidRequest)
	}

	getProfileReq := entity.GetProfileRequest{
		ProfileId: profileId,
	}
	err := s.validate(getProfileReq)
	if err != nil {
		return s.sendValidationErrorResponse(ctx, err)
	}

	result, err := s.profileService.GetProfile(ctx.Request().Context(), getProfileReq)
	if err != nil {
		return s.sendErrorResponse(ctx, err)
	}

	resp := generated.GetProfileResponse{
		FullName:    result.FullName,
		PhoneNumber: result.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) UpdateProfile(ctx echo.Context, params generated.UpdateProfileParams) error {
	profileId, ok := ctx.Get(constant.ProfileIdJwtField).(string)
	if !ok {
		return s.sendErrorResponse(ctx, error_list.ErrInvalidRequest)
	}

	var req generated.UpdateProfileRequest
	err := ctx.Bind(&req)
	if err != nil {
		return s.sendErrorResponse(ctx, error_list.ErrInvalidRequest)
	}

	updateProfileReq := entity.UpdateProfileRequest{
		Id:          profileId,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
	}
	err = s.validate(updateProfileReq)
	if err != nil {
		return s.sendValidationErrorResponse(ctx, err)
	}

	err = s.profileService.UpdateProfile(ctx.Request().Context(), updateProfileReq)
	if err != nil {
		return s.sendErrorResponse(ctx, err)
	}

	resp := generated.UpdateProfileResponse{
		Message: "Success update profile",
	}

	return ctx.JSON(http.StatusOK, resp)
}
