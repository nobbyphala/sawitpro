package main

import (
	"fmt"
	"log"
	"os"
	"sawitpro/constant"
	"sawitpro/generated"
	"sawitpro/handler"
	"sawitpro/helper"
	"sawitpro/repository"
	"sawitpro/service"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	var server = newServer()
	mw, err := server.CreateMiddleware()
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}

	e.Use(middleware.Logger())
	e.Use(mw...)

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func connectDB() (*sqlx.DB, error) {
	connString := fmt.Sprintf("user=%s dbname=%s host=%s port=%s password=%s sslmode=disable",
		constant.EnvPostgresUser,
		constant.EnvPostgresDatabase,
		constant.EnvPostgresHost,
		constant.EnvPostgresPort,
		constant.EnvPostgresPassword,
	)

	db, err := sqlx.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func newServer() *handler.Server {
	conn, err := connectDB()
	if err != nil {
		fmt.Fprintf(os.Stdout, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	//repository
	profileRepository := repository.NewUserProfileRepository(conn)

	//helper
	authHelper := helper.NewAuthHelper()
	validatorHelper := helper.NewValidatorHelper()

	//service
	profileService := service.NewProfileService(service.ProfileServiceDeps{
		ProfileRepository: profileRepository,
		Authhelper:        authHelper,
	})

	opts := handler.NewServerOptions{
		ProfileService:  profileService,
		AuthHelper:      authHelper,
		ValidatorHelper: validatorHelper,
	}

	return handler.NewServer(opts)
}
