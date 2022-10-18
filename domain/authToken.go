package domain

import (
	"banking-auth/errs"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
)

type AuthToken struct {
	token *jwt.Token
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}

func (t AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedString, err := t.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		// TODO: Add Error Handling
		fmt.Println("Failed while signing access token: " + err.Error())
	}
	return signedString, nil
}
