package handler

import (
	"net/http"
	"sawitpro/error_list"
)

var statusResponseMap = map[string]int{
	error_list.ErrProfileRegister.Error():  http.StatusInternalServerError,
	error_list.ErrGetProfile.Error():       http.StatusInternalServerError,
	error_list.ErrProfileNotFound.Error():  http.StatusNotFound,
	error_list.ErrLoginCredential.Error():  http.StatusBadRequest,
	error_list.ErrLogin.Error():            http.StatusInternalServerError,
	error_list.ErrUpdateProfile.Error():    http.StatusInternalServerError,
	error_list.ErrNotAuthenticated.Error(): http.StatusForbidden,
	error_list.ErrInvalidRequest.Error():   http.StatusBadRequest,
	error_list.ErrDataConflict.Error():     http.StatusConflict,
}
