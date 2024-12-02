package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/database"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/email"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/handlers"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/middlewares"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/service"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/tokens"

	docs "github.com/AV3RAGE-ENJOYER/medods_test_task/docs"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           MEDODS Golang test task
// @version         1.0
// @description     This is a test task for Juniour Go Developer in MEDODS.

// @contact.name   Andrei Dombrovskii
// @contact.email  andrushathegames@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey AccessToken
// @in header
// @name Authorization

// @host      127.0.0.1:8080
// @BasePath  /api/v1
func main() {
	// Setup logger

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, handlerOpts))

	slog.SetDefault(logger)

	slog.Info("Logger is set up.")

	// Load .env file

	godotenv.Load("config.env")

	GIN_MODE := os.Getenv("GIN_MODE")
	GIN_ADDR := os.Getenv("GIN_ADDR")

	// Setup Database

	POSTGRES_URL := os.Getenv("POSTGRES_URL")

	db, err := database.NewDB(context.Background(), POSTGRES_URL)

	if err != nil {
		slog.Error("Failed to establish a connection to PostgresDB.", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Database is Set Up.")

	defer db.Pool.Close()

	userService := service.NewUserService(&db)

	// Setup migrations

	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("Failed to set dialect for postgres.", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Migrations are Set Up.")

	driver := stdlib.OpenDBFromPool(db.Pool)

	slog.Info("Migrating database...")

	if err := goose.Up(driver, "migrations"); err != nil {
		slog.Error("Failed to ran migrations.", slog.Any("error", err))
		os.Exit(1)
	}

	if err := driver.Close(); err != nil {
		slog.Error("Failed to close driver.", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Migrations ran successfully!")

	// Setup Token Controller

	JWT_SECRET_KEY := os.Getenv("JWT_SECRET_KEY")

	tokenController, err := tokens.NewTokenController([]byte(JWT_SECRET_KEY), 15*time.Minute, 72*time.Hour)

	if err != nil {
		slog.Error("Failed to initialize jwt token controller.", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Token Controller is Set Up.")

	// Setup EmailService

	mockEmailRepository := email.MockEmailRepository{}
	emailService := service.NewEmailService(&mockEmailRepository)

	slog.Info("Email Service is Set Up.")

	// Setup Gin

	gin.SetMode(GIN_MODE)
	r := setupRouter(userService, tokenController, emailService)

	slog.Info("Gin Router is Set Up")

	r.Run(GIN_ADDR)
}

func setupRouter(userService *service.UserService, tc tokens.TokenController, es *service.EmailService) *gin.Engine {
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/user")
		user.Use(middlewares.AuthMiddleware(userService.Repo, tc))
		user.GET("/ping", handlers.PingHandler())
		user.POST("/add", handlers.AddUserHandler(userService.Repo))

		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.LoginHandler(userService.Repo, tc))
			auth.POST("/refresh", handlers.RefreshTokensHandler(userService.Repo, tc, es.Repo))
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return r
}
