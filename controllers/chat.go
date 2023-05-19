package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-gin-chat/schema"
	"go.mongodb.org/mongo-driver/mongo"
)

type Chat struct {
	Mongo   *mongo.Database
	UserCtl *User // bridge to user controller to get user data
}

func (chat *Chat) GetAll(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal server error",
		})
		return
	}
	c.JSON(200, gin.H{"username": username})
}
