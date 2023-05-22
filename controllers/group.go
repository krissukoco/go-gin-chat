package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/krissukoco/go-gin-chat/models"
	"github.com/krissukoco/go-gin-chat/schema"
	"go.mongodb.org/mongo-driver/mongo"
)

type Group struct {
	Mongo *mongo.Database
}

type NewGroupRequest struct {
	Name      string   `json:"name"`
	MemberIds []string `json:"member_ids"`
}

func (req *NewGroupRequest) Validate() (int, string) {
	if req.Name == "" {
		return schema.ErrFieldRequired, "Name is required"
	}
	return 0, ""
}

func (g *Group) CreateNew(c *gin.Context) {
	userId := c.GetString("userId")
	if userId == "" {
		c.JSON(401, &schema.ErrorResponse{
			Code:    schema.ErrAuthenticationRequired,
			Message: "Unauthorized",
		})
		return
	}
	var req NewGroupRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, &schema.ErrorResponse{
			Code:    schema.ErrUnparsableJSON,
			Message: "Unparsable request",
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
	// Create and save group
	memberIds := req.MemberIds
	if memberIds == nil {
		memberIds = []string{}
	}
	memberIds = append(memberIds, userId)
	now := time.Now().UnixMilli()
	group := &models.Group{
		Name:      req.Name,
		MemberIds: memberIds,
		AdminIds:  []string{userId},
		CreatedAt: now,
		CreatedBy: userId,
		UpdatedAt: now,
	}
	err := group.Save(g.Mongo)
	if err != nil {
		c.JSON(500, &schema.ErrorResponse{
			Code:    schema.ErrInternalServer,
			Message: "Internal server error",
		})
		return
	}
	c.JSON(200, group)
}
