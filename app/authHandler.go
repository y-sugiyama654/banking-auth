package app

import (
	"banking-auth/dto"
	"banking-auth/service"
	"encoding/json"
	"github.com/y-sugiyama654/banking-lib/logger"
	"net/http"
)

type AuthHandler struct {
	service service.AuthService
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logger.Error("Error while decoding login request: " + err.Error())
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
		logger.Error("Error while decoding refresh token request: " + err.Error())
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

func (ah *AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	// converting from query to map type
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		appErr := ah.service.Verify(urlParams)
		if appErr != nil {
			writeResponse(w, http.StatusForbidden, notAuthorizedResponse(appErr.Message))
		} else {
			writeResponse(w, http.StatusOK, authorizedResponse())
		}
	} else {
		writeResponse(w, http.StatusForbidden, notAuthorizedResponse("missing token"))
	}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func notAuthorizedResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      msg,
	}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": true}
}
