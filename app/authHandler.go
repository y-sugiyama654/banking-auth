package app

import (
	"banking-auth/dto"
	"banking-auth/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	service service.AuthService
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		// TODO: Add error log
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appError := ah.service.Login(loginRequest)
		if appError != nil {
			writeResponse(w, appError.Code, appError.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// TODO: Implementing Register API
	writeResponse(w, http.StatusOK, "Register API not implemented yet...")
}

func (ah *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var refreshTokenRequest dto.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&refreshTokenRequest); err != nil {
		// TODO: Add error log
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appError := ah.service.Refresh(refreshTokenRequest)
		if appError != nil {
			writeResponse(w, appError.Code, appError.AsMessage())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}
