package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func CreateUser(c *gin.Context) {
	abortWithErrorMsg := func(msg string) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
	}

	var u *user = &user{}
	var err error = c.ShouldBind(u)

	if err != nil {
		log.Error("Fail to create user: user schema mismatch", err)
		abortWithErrorMsg("Fail to create user")
		return
	}

	if isNameExisted(u) {
		log.Message("Fail to create user: username already exists")
		abortWithErrorMsg("Username already exists")
		return
	}

	//Hash password
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		log.Error("Fail to create user: fail to hash password", err)
		abortWithErrorMsg("Fail to create user")
		return
	}
	u.Password = hashedPassword

	isSaveSuccessful := save(u)

	if isSaveSuccessful {
		c.JSON(http.StatusAccepted, gin.H{"message": "Successful created"})
	} else {
		abortWithErrorMsg("Fail to create user")
	}
}

func AuthUser(c *gin.Context) {
	// Get username and password off request body (can modify in the future to avoid DRY)
	var u *user = &user{}
	var err error = c.ShouldBind(u)
	if err != nil {
		c.String(http.StatusBadRequest, "Fail to log in: ser schema mismatch")
		return
	}

	// Look up requested user from db
	isMatch := isUsernameAndPwdMatch(u.Username, u.Password)
	if !isMatch {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid username or password",
		})
		return
	}

	// Generate jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": u.Username,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), // Jwt expires after 30 days
		"iat": time.Now().Unix(),                          // issue time
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(middleware.SECRET)) // SECRET here should be a env variable e.g. []byte(os.Getenv("SECRET"))

	if err != nil {
		log.Error("", err)
		c.String(http.StatusBadRequest, "Failed to create JWT")
		return
	}

	// Sent jwt back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully generate JWT",
	})
}

func ValidateUser(c *gin.Context) {
	username := c.GetString("Username")
	c.Cookie("Authorization")
	c.JSON(http.StatusOK, gin.H{
		"message":  "Successfully log in",
		"username": username,
	})
}

func LogoutUser(c *gin.Context) {
	// Get username and password off request body (can modify in the future to avoid DRY)
	var u *user = &user{}

	var err error = c.ShouldBind(u)
	if err != nil {
		log.Error("Fail to log out: User schema mismatch", err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "Fail to log out"})
		return
	}

	log.Message(u.toString())

	jwtUsername := c.GetString("Username")
	if jwtUsername != u.Username {
		log.Message("Jwt username and request body username unmatched")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Jwt username and request body username unmatched"})
		return
	}

	// Blacklist user's current jwt
	tokenString, err := c.Cookie("Authorization")
	expiration := c.GetInt64("Exp")

	if err != nil {
		log.Error("Fail to get JWT from cookie", err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "Fail to log out"})
		return
	}

	// save user's jwt to cache
	isSaveSuccessful := saveToCache(u, tokenString, time.Unix(expiration, 0))

	if isSaveSuccessful {
		c.Status(http.StatusAccepted)
	} else {
		log.Message("Fail to log out: fail to save user jwt in cache")
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{"error": "Fail to log out"})
	}

	// Delete corresponding cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete JWT from cookie",
	})
}

func DeleteUser(c *gin.Context) {
	abortWithErrorMsg := func() {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Fail to delete user"})
	}

	// Get username and password off request body (can modify in the future to avoid DRY)
	var u *user = &user{}
	var err error = c.ShouldBind(u)
	if err != nil {
		log.Error("Fail to delete user: user schema mismatch", err)
		abortWithErrorMsg()
		return
	}

	// Logout user
	jwtUsername := c.GetString("Username")
	if jwtUsername != u.Username {
		log.Message("Jwt username and request body username unmatched")
		abortWithErrorMsg()
		return
	}

	// Get user's current jwt; extract expiration time
	tokenString, err := c.Cookie("Authorization")
	expiration := c.GetInt64("Exp")

	if err != nil {
		log.Error("Fail to get JWT from cookie", err)
	}

	// save user's jwt to cache (in order to blacklist the jwt)
	isSaveSuccessful := saveToCache(u, tokenString, time.Unix(expiration, 0))

	if !isSaveSuccessful {
		log.Message("Fail to delete user: fail to save user jwt in cache")
		abortWithErrorMsg()
		return
	}

	// Delete corresponding cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	log.Message("Deleted JWT for user " + u.Username)

	// Delete user info from db
	isDeleteSuccessful := delete(u)
	if isDeleteSuccessful {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Successfully delete user " + u.Username,
		})
	} else {
		log.Message("Fail to delete user: fail to delete user from database")
		abortWithErrorMsg()
	}
}

func ChangePassword(c *gin.Context) {
	abortWithErrorMsg := func(msg string) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
	}

	var p *passwords = &passwords{}

	err := c.ShouldBind(p)
	if err != nil {
		log.Error("Fail to change password: passwords schema unmatch", err)
		abortWithErrorMsg("Fail to change password")
		return
	}

	log.Message("Password binded: " + p.toString())

	username := c.GetString("Username")
	isMatch := isUsernameAndPwdMatch(username, p.OldPassword)
	if !isMatch {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong old password",
		})
		return
	}

	// If old pwd is the same as the new one, abort
	if p.isSame() {
		log.Error("Fail to change password: old password is the same as new password", err)
		abortWithErrorMsg("Old password cannot be the same as new password")
		return
	}

	// hash new password
	hashedPassword, err := hashPassword(p.NewPassword)
	if err != nil {
		log.Error("Fail to create user: fail to hash password", err)
		abortWithErrorMsg("Fail to change password")
		return
	}

	// update username in db
	res := updatePassword(username, hashedPassword)
	if !res {
		log.Message("Fail to store new password to db")
		abortWithErrorMsg("Fail to change password")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed password"})
}
