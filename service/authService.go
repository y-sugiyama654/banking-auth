package service

import (
	"banking-auth/domain"
	"banking-auth/dto"
	"github.com/golang-jwt/jwt/v4"
	"github.com/y-sugiyama654/banking-lib/errs"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Refresh(dto.RefreshTokenRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(map[string]string) *errs.AppError
}

type DefaultAuthService struct {
	repo            domain.AuthRepository
	rolePermissions domain.RolePermissions
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

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	// convert the string token to JWT struct
	if jtwToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		// TODO: Add Error Handling
		return nil
	} else {
		// Checking the validity of the token, this verifies the expiry time and the signature of the token
		if jtwToken.Valid {
			// type cast the token claims to jwt.MapClaims
			claims := jtwToken.Claims.(*domain.AccessTokenClaims)
			if claims.IsUserRole() {
				if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
					// TODO: Add Error Handling
					return nil
				}
			}
			isAuthorized := s.rolePermissions.IsAuthorizedFor(claims.Role, urlParams["routeName"])
			if !isAuthorized {
				// TODO: Add Error Handling
				return nil
			}
			return nil
		} else {
			// TODO: Add Error Handling
			return nil
		}
	}
}

func NewAuthService(repository domain.AuthRepositoryDb, permissions domain.RolePermissions) DefaultAuthService {
	return DefaultAuthService{repository, permissions}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		// TODO: Add Error log
		return nil, err
	}
	return token, nil
}
