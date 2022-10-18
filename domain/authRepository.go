package domain

import (
	"banking-auth/errs"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type AuthRepositoryDb struct {
	client *sqlx.DB
}

type AuthRepository interface {
	FindBy(username string, password string) (*Login, *errs.AppError)
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
			// TODO: Add error handling
			fmt.Println("invalid credentials")
			fmt.Println(err.Error())
		} else {
			// TODO: Add error handling
			fmt.Println("Unexpected database error")
			fmt.Println(err.Error())
		}
	}

	return &login, nil
}
