package service

import (
	"banking-auth/domain"
	"banking-auth/dto"
	"banking-auth/errs"
	"github.com/golang-jwt/jwt/v4"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Refresh(dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError)
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

	// TODO: Implementing Refresh Token

	return &dto.LoginResponse{AccessToken: accessToken}, nil
}

func (s DefaultAuthService) Refresh(req dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError) {
	// TODO: エラーハンドリングは後で実装する
	var errTmp *errs.AppError

	if vErr := req.IsAccessTokenValid(); vErr != nil {
		if vErr.Errors == jwt.ValidationErrorExpired {
			// continue with the refresh token functionality
			var appErr *errs.AppError
			if appErr = s.repo.RefreshTokenExists(req.RefreshToken); appErr != nil {
				return nil, appErr
			}
			// generate a access token from refresh token
			var accessToken string
			if accessToken, appErr = domain.NewAccessTokenFromRefreshToken(req.RefreshToken); appErr != nil {
				return nil, appErr
			}
			return &dto.LoginResponse{AccessToken: accessToken}, nil
		}
		return nil, errTmp
	}
	return nil, errTmp
}

func NewAuthService(repository domain.AuthRepositoryDb) DefaultAuthService {
	return DefaultAuthService{repository}
}
