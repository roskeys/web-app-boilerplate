package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/roskeys/app/utils"
)

func AccessTokenCheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")
		if err != nil {
			utils.SendErrorResponse(c, utils.NO_COOKIE_FOUND)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, utils.UnexpectedSigningMethodErr
			}
			return []byte(utils.JWT_SECRET), nil
		})
		if err != nil {
			utils.SendErrorResponse(c, utils.COOKIE_EXPIRED)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// check token is expired or not
			if int64(claims["exp"].(float64)) < time.Now().Unix() {
				utils.SendErrorResponse(c, utils.COOKIE_EXPIRED)
				return
			}
			c.Set("uid", claims["uid"])
		} else {
			utils.SendErrorResponse(c, utils.COOKIE_EXPIRED)
			return
		}
		c.Next()
	}
}
