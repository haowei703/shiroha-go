package auth

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApiMiddleware(keyCloakApi *KeyCloakApi) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("keyCloakApi", keyCloakApi)
		c.Next()
	}
}

func AuthenticateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, errT := c.Cookie("csrf_token")
		refreshToken, errR := c.Cookie("refresh_token")
		if refreshToken == "" || token == "" || errT != nil || errR != nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

		api, exists := c.Get("keyCloakApi")
		if !exists {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		keyCloakApi := api.(*KeyCloakApi)
		loggedIn, err := keyCloakApi.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		if !gocloak.PBool(loggedIn) {
			var jwt *gocloak.JWT
			jwt, err = keyCloakApi.RefreshToken(refreshToken)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
			c.SetCookie("csrf_token", jwt.AccessToken, jwt.ExpiresIn, "/", "localhost", false, true)
			c.SetCookie("refresh_token", jwt.RefreshToken, jwt.RefreshExpiresIn, "/", "localhost", false, true)
		}
		c.Next()
	}
}
