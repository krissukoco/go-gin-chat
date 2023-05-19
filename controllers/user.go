package controllers

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-gin-chat/models"
	"github.com/krissukoco/go-gin-chat/schema"
	"gorm.io/gorm"
)

type User struct {
	Pg *gorm.DB
}

var (
	ErrInvalidPage = errors.New("invalid page")
	ErrInvalidSize = errors.New("invalid size")
)

func (u *User) getPage(c *gin.Context) (int, error) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return 0, err
	}
	if page < 1 {
		return 0, ErrInvalidPage
	}
	return page, nil
}

func (u *User) getSize(c *gin.Context) (int, error) {
	sizeStr := c.DefaultQuery("size", "10")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return 0, err
	}
	if size < 1 {
		return 0, ErrInvalidSize
	}
	return size, nil
}

func (u *User) GetUserById(id string) (*models.User, error) {
	var user models.User
	err := user.FindById(u.Pg, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(400, &schema.ErrorResponse{
			Code:    schema.ErrFieldRequired,
			Message: "user id is required",
		})
		return
	}
	// Find user by id
	user, err := u.GetUserById(id)
	if err != nil {
		if err == models.ErrUserNotFound {
			c.JSON(404, &schema.ErrorResponse{
				Code:    schema.ErrResourceNotFound,
				Message: "user not found",
			})
			return
		}
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal server error",
		})
		return
	}
	c.JSON(200, &user)
}

func (u *User) GetAll(c *gin.Context) {
	page, err := u.getPage(c)
	if err != nil {
		c.JSON(400, &schema.ErrorResponse{
			Code:    schema.ErrFieldInvalid,
			Message: "Invalid page query",
		})
		return
	}
	size, err := u.getSize(c)
	if err != nil {
		c.JSON(400, &schema.ErrorResponse{
			Code:    schema.ErrFieldInvalid,
			Message: "Invalid size query",
		})
		return
	}
	offset := (page - 1) * size
	users, err := models.GetUsers(u.Pg, offset, size)
	if err != nil {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(200, &users)
}
