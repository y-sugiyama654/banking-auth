package app

import (
	"banking-auth/domain"
	"banking-auth/service"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"time"
)

func Start() {
	// 環境変数のチェック
	sanityCheck()
	router := mux.NewRouter()

	// dependency injection
	dbClient := getDbClient()
	authRepository := domain.NewAuthRepository(dbClient)
	authService := service.NewAuthService(authRepository, domain.GetRolePermissions())
	ah := AuthHandler{authService}

	// router
	router.HandleFunc("/auth/login", ah.Login).Methods(http.MethodPost)
	router.HandleFunc("/auth/register", ah.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/refresh", ah.Refresh).Methods(http.MethodPost)
	router.HandleFunc("/auth/verify", ah.Verify).Methods(http.MethodGet)

	// starting server
	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), router))
}

func sanityCheck() {
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_ADDR",
		"DB_PORT",
		"DB_NAME",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			// TODO: Edit Error Log
			fmt.Sprintf("Environment variable %s not defined. Terminating application.", k)
		}
	}
}

func getDbClient() *sqlx.DB {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbAddr := os.Getenv("DB_ADDR")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbAddr, dbPort, dbName)
	client, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	// 接続を再利用できる期間
	client.SetConnMaxLifetime(time.Minute * 3)
	// 最大何本接続できるかの上限
	client.SetMaxOpenConns(10)
	// idle接続を最大何本保持できるかの上限
	client.SetMaxIdleConns(10)

	return client
}
