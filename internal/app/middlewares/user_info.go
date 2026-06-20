package middlewares

import (
	"acc_backend/internal/app/helpers"
	"acc_backend/internal/container"
	"strings"

	"github.com/gin-gonic/gin"
)

// userRepo :=

// *model.User
func GetUser(cnt *container.Container) gin.HandlerFunc {

	return func(c *gin.Context) {
		prefix := "Bearer: "
		authHeaderStr := c.GetHeader("Authorization")
		if authHeaderStr == "" {
			helpers.Fail(c, 401, "UnAuthorized")
		}

		if !strings.HasPrefix(authHeaderStr, prefix) {
			helpers.Fail(c, 401, "UnAuthorized")
		}

		tokenStr, found := strings.CutPrefix(authHeaderStr, prefix)

		if !found {
			helpers.Fail(c, 401, "UnAuthorized")
		}

		userId, err := cnt.JwtService.Verify(tokenStr)

		if err != nil {
			helpers.Fail(c, 401, "UnAuthorized")
		}
		if userId == "" {
			helpers.Fail(c, 401, "UnAuthorized")
		}
		c.Set("userId", userId)
		// Pre-handler phase
		c.Next()
	}
}

func extracTokenFromHeaderString(authHeaderStr string) string {
	return authHeaderStr[8:]
}
