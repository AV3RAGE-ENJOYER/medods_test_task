package tests

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/database"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/email"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/handlers"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/middlewares"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/models"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/service"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/tokens"

	"github.com/gin-gonic/gin"
)

func setupRouter(userService *service.UserService, tc tokens.TokenController, es *service.EmailService) *gin.Engine {
	// Disable logs in test mode

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	slog.SetDefault(logger)

	r := gin.New()

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
	return r
}

func sendRequest(user handlers.UserRequest, path string) *http.Request {
	userJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest(
		"POST",
		path,
		strings.NewReader(string(userJSON)),
	)
	req.Header.Set("Content-Type", "application/json")
	return req
}

var MockDB database.MockDBUserRepository = database.MockDBUserRepository{
	Users: map[string]models.User{
		"admin@gmail.com": {
			Email:          "admin@gmail.com",
			HashedPassword: "$2a$12$7TAiwtMzHZg49781OZwni.CTTeBIWKYhjkNIh/1uL8MdRK9RFMwmK",
		},
	},
	RefreshTokens: make(map[string]models.RefreshToken),
}

var MockUserService *service.UserService = service.NewUserService(&MockDB)

var MockTokenController tokens.TokenController = tokens.TokenController{
	SigningKey:      []byte("test-secret-key"),
	AccessTokenTTL:  15 * time.Minute,
	RefreshTokenTTL: 72 * time.Hour,
}

var MockEmailRepository email.MockEmailRepository = email.MockEmailRepository{}
var MockEmailService *service.EmailService = service.NewEmailService(&MockEmailRepository)

var router *gin.Engine = setupRouter(MockUserService, MockTokenController, MockEmailService)
