package domain

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/y-sugiyama654/banking-lib/errs"
	"github.com/y-sugiyama654/banking-lib/logger"
)

type AuthRepositoryDb struct {
	client *sqlx.DB
}

type AuthRepository interface {
	FindBy(username string, password string) (*Login, *errs.AppError)
	RefreshTokenExists(string) *errs.AppError
}

func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}

func (d AuthRepositoryDb) FindBy(username string, password string) (*Login, *errs.AppError) {
	var login Login
	sqlVerify := `SELECT username, u.customer_id, role, group_concat(a.account_id) as account_numbers FROM users u
                  LEFT JOIN accounts a ON a.customer_id = u.customer_id
                WHERE username = ? and password = ?
                GROUP BY a.customer_id`
	if err := d.client.Get(&login, sqlVerify, username, password); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database: " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error.")
		}
	}
	return &login, nil
}

func (d AuthRepositoryDb) RefreshTokenExists(refreshToken string) *errs.AppError {
	sqlSelect := "select refresh_token from refresh_token_store where refresh_token = ?"
	var token string
	err := d.client.Get(&token, sqlSelect, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewAuthenticationError("refresh token not registered in the store")
		} else {
			logger.Error("Unexpected database error: " + err.Error())
			return errs.NewUnexpectedError("unexpected database error")
		}
	}
	return nil
}
