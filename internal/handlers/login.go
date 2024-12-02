package handlers

import (
	"log/slog"

	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/service"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/tokens"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// @BasePath /api/v1

// @Summary Login user
// @Description Login user via email and password
// @Tags Authentication
// @Produce json
// @Accept json
// @Param body body handlers.UserRequest true "User email and password"
// @Success 200 {object} tokens.TokenPair
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/login [post]
func LoginHandler(db service.UserRepository, tc tokens.TokenController) gin.HandlerFunc {
	return func(c *gin.Context) {
		var userRequest UserRequest

		if err := c.BindJSON(&userRequest); err != nil {
			slog.Error("Failed to Unmarshall JSON", slog.Any("error", err))
			c.String(500, "Internal server error")
			return
		}

		user, err := db.GetUser(c, userRequest.Email)

		if err == pgx.ErrNoRows {
			slog.Error("Failed to Get User.", slog.Any("error", err))
			c.String(404, "User not found")
			return
		}

		err = user.CheckPassword(userRequest.Password)

		if err != nil {
			slog.Error("Invalid Password.", slog.Any("error", err))
			c.String(401, "Unauthorized")
			return
		}

		tokenPair, err := tc.NewJWT(user.Email, c.ClientIP())

		if err != nil {
			slog.Error("Failed to Generate Token Pair.", slog.Any("error", err))
			c.String(500, "Internal server error")
			return
		}

		refreshTokenHash, err := utils.HashBCrypt([]byte(tokenPair.RefreshToken))

		if err != nil {
			slog.Any("Failed to Hash Refresh Token.", slog.Any("error", err))
			c.String(500, "Internal server error")
			return
		}

		err = db.AddRefreshToken(
			c,
			user.Email,
			string(refreshTokenHash),
			tc.RefreshTokenTTL,
			c.ClientIP(),
		)

		if err != nil {
			slog.Error("Failed to Add Refresh Token.", slog.Any("error", err))
			c.String(500, "Internal server error")
			return
		}

		c.IndentedJSON(200, tokenPair)
	}
}
