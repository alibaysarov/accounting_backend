package middlewares

import (
	"acc_backend/internal/app/helpers"
	"acc_backend/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// userRepo :=

// *model.User

type AuthMiddleware struct {
	jwtService *service.JwtService
	userRepo   service.UserRepository
}

func NewAuthMiddleware(jwtService *service.JwtService, userRepo service.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{jwtService: jwtService, userRepo: userRepo}
}

func (mdl *AuthMiddleware) GetUser() gin.HandlerFunc {

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

		userId, err := mdl.jwtService.Verify(tokenStr)

		if err != nil {
			helpers.Fail(c, 401, "UnAuthorized")
		}
		if userId == "" {
			helpers.Fail(c, 401, "UnAuthorized")
		}
		user, err := mdl.userRepo.GetById(c.Request.Context(), userId)

		if err != nil {
			helpers.Fail(c, 401, "UnAuthorized")
		}
		c.Set("user", user)
		// Pre-handler phase
		c.Next()
	}
}

func extracTokenFromHeaderString(authHeaderStr string) string {
	return authHeaderStr[8:]
}
