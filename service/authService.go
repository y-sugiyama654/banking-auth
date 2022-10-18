package service

import (
	"banking-auth/domain"
	"banking-auth/dto"
	"banking-auth/errs"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
}

type DefaultAuthService struct {
	repo domain.AuthRepository
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, *errs.AppError) {
	var appErr *errs.AppError
	var login *domain.Login

	if login, appErr = s.repo.FindBy(req.UserName, req.UserPassword); appErr != nil {
		return nil, appErr
	}

	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)

	var accessToken string
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}
	// Debug
	//fmt.Println(accessToken)

	// TODO: Implementing Refresh Token

	return &dto.LoginResponse{AccessToken: accessToken}, nil
}

func NewAuthService(repository domain.AuthRepositoryDb) DefaultAuthService {
	return DefaultAuthService{repository}
}
