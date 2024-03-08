package handler

import (
	"context"
	"net/http"
	"sawitpro/constant"
	"sawitpro/error_list"
	"sawitpro/generated"
	"sawitpro/helper"
	"sawitpro/service"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
)

type Server struct {
	profileService  service.ProfileServiceInterface
	authHelper      helper.AuthHelperInterface
	validatorHelper helper.ValidatorHelperInterface
}

type NewServerOptions struct {
	ProfileService  service.ProfileServiceInterface
	AuthHelper      helper.AuthHelperInterface
	ValidatorHelper helper.ValidatorHelperInterface
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		profileService:  opts.ProfileService,
		authHelper:      opts.AuthHelper,
		validatorHelper: opts.ValidatorHelper,
	}
}

func (srv *Server) sendErrorResponse(ctx echo.Context, err error) error {
	var statusCode int
	var errorMessage = err.Error()

	code, exists := statusResponseMap[errorMessage]
	if exists {
		statusCode = code
	} else {
		statusCode = http.StatusInternalServerError
	}

	resp := generated.ErrorResponse{
		Message: errorMessage,
	}

	return ctx.JSON(statusCode, resp)
}

func (srv *Server) sendValidationErrorResponse(ctx echo.Context, err error) error {
	var errorMessage = err.Error()

	resp := generated.ErrorResponse{
		Message: errorMessage,
	}

	return ctx.JSON(http.StatusBadRequest, resp)
}

func (srv *Server) getJWSFromRequest(req *http.Request) (string, error) {
	authHdr := req.Header.Get("Authorization")

	if authHdr == "" {
		return "", error_list.ErrNotAuthenticated
	}

	prefix := "Bearer "
	if !strings.HasPrefix(authHdr, prefix) {
		return "", error_list.ErrNotAuthenticated
	}
	return strings.TrimPrefix(authHdr, prefix), nil
}

func (srv *Server) CreateMiddleware() ([]echo.MiddlewareFunc, error) {
	spec, err := generated.GetSwagger()
	if err != nil {
		return nil, err
	}

	authenticator := middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
				token, err := srv.getJWSFromRequest(input.RequestValidationInput.Request)
				if err != nil {
					return err
				}

				profileId, err := srv.authHelper.VerifyToken(ctx, token)
				if err != nil {
					return err
				}

				eCtx := middleware.GetEchoContext(ctx)
				eCtx.Set(constant.ProfileIdJwtField, profileId)

				return nil
			},
		},
	})

	return []echo.MiddlewareFunc{authenticator}, nil
}

func (srv *Server) validate(obj interface{}) error {
	return srv.validatorHelper.ValidateStruct(obj)
}
