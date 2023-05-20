package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-gin-chat/models"
	"github.com/krissukoco/go-gin-chat/schema"
	"github.com/krissukoco/go-gin-chat/security"
	"gorm.io/gorm"
)

type Auth struct {
	Pg        *gorm.DB
	JwtSecret string
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate: validate the request body
// return common error code, 0 if no error
func (lreq *LoginRequest) Validate() int {
	if lreq.Username == "" {
		return schema.ErrFieldRequired
	}
	if lreq.Password == "" {
		return schema.ErrFieldRequired
	}
	return 0
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Name            string `json:"name"`
	Location        string `json:"location"`
}

func (reg *RegisterRequest) Validate() (int, string) {
	if reg.Username == "" {
		return schema.ErrFieldRequired, "Username is required"
	}
	if len(reg.Username) < 3 {
		return schema.ErrFieldMinChar, "Username must be at least 3 characters"
	}
	if reg.Password == "" {
		return schema.ErrFieldRequired, "Password is required"
	}
	if len(reg.Password) < 8 {
		return schema.ErrPasswordMinChar, "Password must be at least 8 characters"
	}
	if reg.ConfirmPassword != reg.Password {
		return schema.ErrPasswordUnmatch, "Password and Confirm Password must match"
	}
	if reg.Name == "" {
		return schema.ErrFieldRequired, "Name is required"
	}
	return 0, ""
}

func (a *Auth) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, &schema.ErrorResponse{
			Code:    schema.ErrUnparsableJSON,
			Message: "Unparsable JSON",
		})
		return
	}
	if errCode := req.Validate(); errCode != 0 {
		c.JSON(400, &schema.ErrorResponse{
			Code:    errCode,
			Message: "Validation Error",
		})
		return
	}
	var u models.User
	err := u.FindByUsername(a.Pg, req.Username)
	if err != nil {
		c.JSON(400, &schema.ErrorResponse{
			Code:    schema.ErrEmailOrPasswordInvalid,
			Message: "Email or Password Invalid",
		})
		return
	}
	err = u.ComparePassword(a.Pg, req.Password)
	if err != nil {
		c.JSON(400, &schema.ErrorResponse{
			Code:    schema.ErrEmailOrPasswordInvalid,
			Message: "Email or Password Invalid",
		})
		return
	}
	token, err := security.JwtFromUserId(u.Id, a.JwtSecret)
	if err != nil {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal Server Error",
		})
		return
	}
	c.JSON(200, gin.H{
		"token": token,
		"user":  &u,
	})
}

func (a *Auth) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(422, &schema.ErrorResponse{
			Code:    schema.ErrUnparsableJSON,
			Message: "Unparsable JSON",
		})
		return
	}
	code, msg := req.Validate()
	if code != 0 {
		c.JSON(400, &schema.ErrorResponse{
			Code:    code,
			Message: msg,
		})
		return
	}
	// Ensure username is not taken
	var u models.User
	err := u.FindByUsername(a.Pg, req.Username)
	if err == nil {
		c.JSON(400, &schema.ErrorResponse{
			Code:    schema.ErrUsernameAlreadyTaken,
			Message: "Username is already taken",
		})
		return
	}
	u.Username = req.Username
	u.Password = req.Password
	u.Name = req.Name
	u.Location = req.Location
	err = u.HashPassword()
	if err != nil {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal Server Error",
		})
		return
	}
	err = u.Save(a.Pg)
	if err != nil {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal Server Error",
		})
		return
	}
	c.JSON(200, &u)
}

func (a *Auth) GetAccount(c *gin.Context) {
	userId := c.GetString("username")
	if userId == "" {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal Server Error",
		})
		return
	}

	var u models.User
	err := u.FindById(a.Pg, userId)
	if err != nil {
		c.JSON(401, &schema.ErrorResponse{
			Code:    schema.ErrTokenInvalid,
			Message: "Invalid token",
		})
		return
	}

	c.JSON(200, &u)
}
