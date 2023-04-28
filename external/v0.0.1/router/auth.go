package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/maximotejeda/auth-service/external/v0.1/database"
	"github.com/maximotejeda/auth-service/external/v0.1/helper"
	"github.com/maximotejeda/auth-service/external/v0.1/mail"
)

// AuthAddRoutes: add the different  routes to router differentiating from validated or not
func AuthAddRoutes(r *gin.Engine) {
	v01 := r.Group("/v0.1")
	auth := v01.Group("/auth")
	{
		auth.POST("/login", login)
		auth.POST("/register", register)
		auth.POST("/refresh", refresh)
		auth.GET("/validate", validate)
		auth.POST("/recover", newRecover)
		auth.PUT("/endrecover", finishRecover)
		auth.GET("/activate", activateAccount)
	}
	validated := auth.Group("/").Use(helper.Validated())
	{
		validated.GET("/logout", logout)
		validated.GET("/pubkey", publicKey)
	}

}

// login: Handler that will manage login on the app
func login(c *gin.Context) {
	input := database.UserLogin{}
	j := J
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO try to lookup DB if user exist and PASSWORd is the same proced
	// otherwise return error
	// TODO: token validation atime must be passed
	user := &database.User{UserLogin: input}
	userLoged, err := database.LoginAccount(user)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}
	if !userLoged.Active {
		fmt.Printf("User %v", user)
		c.JSON(http.StatusForbidden, gin.H{"error": "account not active check your email for activation, or ask for an activation emails"})
		return
	}
	claims := map[string]interface{}{
		"username": userLoged.Username,
		"email":    userLoged.Email,
		//	"id":        userLoged.ID,
		//	"firstname": userLoged.UserInfo.FirstName,
		//	"lastname":  userLoged.UserInfo.LastName,
		"rol":   userLoged.Rol,
		"loged": userLoged.UserInfo.Loged,
	}

	s, err := j.Create(claims)

	if err != nil {
		log.Print("Create token: ", err.Error())
	}

	_, err = j.Validate(s)
	if err != nil {
		fmt.Printf("validating err: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Authorization", "Bearer "+s)

	c.JSON(http.StatusOK, gin.H{"token": s})
}

// validate: handler that will validate if a token is valid
func validate(c *gin.Context) {
	token := c.Query("token")
	s := token[7:]
	_, err := J.Validate(s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Authorization", "Bearer "+token)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// publickKey: handler tha will give the public key for external app validation
func publicKey(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"publicKey": J.PublicPemStr})
}

// logout: will logout the session
func logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// refresh: will issue a token if the past token is expired for less than a number
func refresh(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
		return
	}
	if token == "" || len(strings.Split(token, ".")) < 3 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is needed for a refresh"})
		return
	}
	s := token[7:]
	newToken, err := J.RefreshToken(s)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Authorization", "Bearer "+newToken)
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// register: handler let user register on the service
// ALL fields are mandatory, username, password, email, firstname, lastname
// on create active actual false
func register(c *gin.Context) {
	defaultActive := os.Getenv("DEFAULTACTIVEUSER")
	input := database.Register{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// User active or inactive for default
	user := database.User{Active: defaultActive != "", UserLogin: input.UserLogin, UserInfo: input.UserInfo}
	err := database.CreateAccount(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//TODO create a email verification token sender
	// to finish activating account
	token, err := J.Create(user)
	fmt.Println(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Writer.Header().Set("Authorization", "Bearer "+token)
	host := c.Request.Host
	mail.SendEmail(user.Email, token, "activate", "0", host)
	// the above code is innecessary can be deleted just for testing purpose
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

// newRecover: launch account recovery proccess
// if usser remmber user name or email
// if that info is on DB
// an email willl be sent with a code
// update record on database to expect a code and measure time to do it
// the param must be named data
// need to be a function working out of context to disable recover after the token time
func newRecover(c *gin.Context) {
	type data struct {
		Data string `json:"data" binding:"required"`
	}
	input := data{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := database.AccountExist(input.Data)
	fmt.Println("Result of look up user is: ", res)
	user := database.NewRecoverAccount(input.Data)
	//TODO logic to send email will be available for token live
	// the email will redirect to a page with a token issued in the url querys
	// if the token and the pin are ok
	// be able to change password in the token live
	if user != nil {
		token, err := J.Create(user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		host := c.Request.Host
		mail.SendEmail(user.Email, token, "recover", user.RecoverAccount.RecoverCode, host)

	}
	c.JSON(http.StatusCreated, gin.H{"status": "if data exist an email will be sent in the next 10 minutes"})
}

func finishRecover(c *gin.Context) {
	type pin struct {
		Pin int // number to compare
	}
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
		return
	}
	sToken := token[7:]
	params, err := J.Validate(sToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	input := pin{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, ok := params["username"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not in claims"})
		return
	}
	us, ok := username.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not an string claims"})
		return
	}
	/// create user to query Db
	user := database.User{UserLogin: database.UserLogin{Username: us}}
	id, err := database.GetID(&user)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "getting id"})
		return
	}
	recovery, err := database.GetRecoverInfo(id)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "getting recovery"})
		return
	}
	if recovery.RecoverCode != fmt.Sprintf("%d", input.Pin) {
		fmt.Printf("error on comparing error Want: %s, Got: %d", recovery.RecoverCode, input.Pin)
		return
	}
	database.EndRecoverAccount(id)
	c.Writer.Header().Set("Authorization", "Bearer "+sToken)
	c.JSON(http.StatusAccepted, gin.H{"status": "account proccesed"})
}

// Working expect to recive email from service with token to activate
func activateAccount(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization required"})
		return
	}
	sToken := token[7:]
	params, err := J.Validate(sToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	username, ok := params["username"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not in claims"})
		return
	}
	us, ok := username.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not an string claims"})
		return
	}
	/// create user to query Db
	user := database.User{UserLogin: database.UserLogin{Username: us}}
	database.ActivateAccount(&user)
}
