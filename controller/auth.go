package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"github.com/golang-jwt/jwt/v4"
	"github.com/roskeys/app/db"
	"github.com/roskeys/app/middleware"
	"github.com/roskeys/app/utils"
	"golang.org/x/crypto/bcrypt"
)

func InitAuthController() {
	auth_router := utils.MainRouter.Group("/auth")
	auth_router.Use(middleware.IPRateLimiter(redis_rate.PerMinute(10)))
	{
		auth_router.POST("/signup", signupHandler)
		auth_router.POST("/login", loginHandler)
		auth_router.POST("/logout", logoutHandler)
		auth_router.GET("/refresh", middleware.AccessTokenCheckJWT(), tokenRefreshHandler)
		auth_router.GET("/test", middleware.AccessTokenCheckJWT(), test)
	}
}

func test(c *gin.Context) {
	fmt.Println(c.GetString("uid"))
	c.String(200, "Hello world")
}

func signupHandler(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		utils.SendErrorResponse(c, utils.INVALID_REGISTER_FORM)
		return
	}
	uid, e := db.CreateNewUser(input.Email, input.Username, input.Password)
	if len(uid) == 0 {
		utils.SendErrorResponse(c, e)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success", "uid": uid})
}

func loginHandler(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBind(&input); err != nil {
		utils.SendErrorResponse(c, utils.INVALID_CREDENTIALS)
		return
	}
	user := db.GetUserByEmail(input.Email)
	if len(user.UID) == 0 {
		utils.SendErrorResponse(c, utils.INVALID_CREDENTIALS)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		utils.SendErrorResponse(c, utils.INVALID_CREDENTIALS)
		return
	}
	token, err := generateJWTToken(user.UID)
	if err != nil {
		utils.SendErrorResponse(c, utils.INTERNAL_SERVER_ERROR)
		return
	}
	c.SetCookie("access_token", token, int(time.Hour*48), "/", utils.DOMAIN_NAME, false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"user": gin.H{
			"username": user.Username,
			"email":    user.Email,
			"uid":      user.UID,
		},
	})
}

func logoutHandler(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", utils.DOMAIN_NAME, false, true)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func tokenRefreshHandler(c *gin.Context) {
	uid := c.GetString("uid")
	token, err := generateJWTToken(uid)
	if err != nil {
		utils.SendErrorResponse(c, utils.INTERNAL_SERVER_ERROR)
		return
	}
	c.SetCookie("access_token", token, int(time.Hour*48), "/", utils.DOMAIN_NAME, false, true)
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func generateJWTToken(uid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  uid,
		"exp": time.Now().Add(time.Hour * 48).Unix(),
	})
	tokenString, err := token.SignedString([]byte(utils.JWT_SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
