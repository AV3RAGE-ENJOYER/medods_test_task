package handlers

import (
	"log/slog"
	"time"

	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/service"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/tokens"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// @BasePath /api/v1

// @Summary Refresh tokens
// @Description Returns new access and refresh tokens if refresh token is correct and not expired.
// @Tags Authentication
// @Produce json
// @Param refresh_token query string true "Refresh Token"
// @Param email query string true "User email"
// @Success 200 {object} tokens.TokenPair
// @Failure 400 {string} Bad request
// @Failure 404 {string} No refresh token found
// @Failure 500 {string} Internal server error
// @Router /auth/refresh [post]
func RefreshTokensHandler(db service.UserRepository, tc tokens.TokenController, es service.EmailRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Request.URL.Query()

		refreshTokenString := query["refresh_token"]
		email := query["email"]

		if refreshTokenString == nil || email == nil {
			c.String(400, "Bad request")
			return
		}

		refreshToken, err := db.GetRefreshTokenProps(
			c, email[0])

		if err == pgx.ErrNoRows {
			c.String(404, "No refresh token found")
			return
		} else if err != nil {
			slog.Error("Failed to get refresh token info.", slog.Any("error", err))
			c.String(500, "Internal server error")
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(refreshToken.RefreshTokenHash),
			[]byte(refreshTokenString[0]))

		if err != nil {
			slog.Error("Invalid Refresh Token.", slog.Any("error", err))
			c.String(401, "Unauthorized")
			return
		}

		if time.Now().Unix()-refreshToken.ExpiresAt.Unix() > 0 {
			slog.Info("Refresh Token has Expired.", slog.Any("email", email[0]))
			c.String(401, "Unauthorized")
			return
		}

		if refreshToken.IpAddress != c.ClientIP() {
			es.NotifyUser(c, email[0])
		}

		newTokenPair, err := tc.NewJWT(email[0], c.ClientIP())

		if err != nil {
			slog.Error("Failed to Generate Token Pair", slog.Any("error", err))
			c.String(500, "Internal server error")
			return
		}

		newRefreshTokenHash, err := utils.HashBCrypt([]byte(newTokenPair.RefreshToken))

		if err != nil {
			slog.Error("Failed to Hash Refresh Token.", slog.Any("error", err))
			c.String(500, "Internal server error")
			return
		}

		err = db.AddRefreshToken(
			c,
			email[0],
			string(newRefreshTokenHash),
			tc.RefreshTokenTTL,
			c.ClientIP(),
		)

		if err != nil {
			slog.Error("Failed to Add Refresh Token.")
			c.String(500, "Internal server error")
			return
		}

		slog.Info("New Refresh Token Added.")

		c.IndentedJSON(200, newTokenPair)
	}
}
