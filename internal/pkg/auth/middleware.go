package auth

import "github.com/gin-gonic/gin"

func ApiMiddleware(keyCloakApi *KeyCloakApi) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("keyCloakApi", keyCloakApi)
		c.Next()
	}
}
