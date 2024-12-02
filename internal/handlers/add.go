package handlers

import (
	"log/slog"

	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/models"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/service"
	"github.com/AV3RAGE-ENJOYER/medods_test_task/internal/utils"

	"github.com/gin-gonic/gin"
)

// @BasePath /api/v1

// @Summary Add user
// @Description Adds user's email and hashed password to the database
// @Tags User
// @Produce json
// @Param body body handlers.UserRequest true "Request body"
// @Success 200 {object} models.User
// @Failure 400 {string} Bad request
// @Failure 500 {string} Internal server error
// @Security AccessToken
// @Router /user/add [post]
func AddUserHandler(db service.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user UserRequest

		if err := c.BindJSON(&user); err != nil {
			c.String(400, "Bad request")
			return
		}

		slog.Debug("Add user request.", slog.Any("request", user))

		if user.Email != "" && user.Password != "" {
			hashedPassword, err := utils.HashBCrypt([]byte(user.Password))

			if err != nil {
				slog.Error("Failed to hash a password.", slog.Any("error", err))
				c.String(500, "Failed to hash a password")
				return
			}

			user := models.User{
				Email:          user.Email,
				HashedPassword: string(hashedPassword),
			}

			err = db.AddUser(c, user)

			if err != nil {
				slog.Error("Failed to add user", slog.Any("error", err))
				c.String(500, "Failed to add user")
				return
			}

			slog.Debug("Successfully added user.")

			c.IndentedJSON(200, user)
			return
		}

		c.String(400, "Bad request")
	}
}
