// Package api
package api

import (
	"errors"
	"net/http"

	"github.com/joaoguilherme2909/goBid/internal/jsonutils"
	"github.com/joaoguilherme2909/goBid/internal/services"
	"github.com/joaoguilherme2909/goBid/internal/usecase/user"
)

func (api *Api) handleSignupUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.CreateUserReq](r)
	if err != nil {
		_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
			_ = jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "email or username already exists",
			})
			return
		}
	}

	_ = jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
}

func (api *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[user.LoginUserReq](r)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
	}

	id, err := api.UserService.AuthenticateUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]any{
				"error": "invalid email or password",
			})
			return
		}
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpedcted internal server error",
		})
		return
	}

	err = api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpedcted server error",
		})
		return
	}

	api.Sessions.Put(r.Context(), "AuthenticatedUserId", id)

	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{
		"message": "logged in, sucessfully",
	})
}

func (api *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := api.Sessions.RenewToken(r.Context())
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]any{
			"error": "unexpedcted server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), "AuthenticatedUserId")
	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]any{
		"message": "logged out, sucessfully",
	})
}
