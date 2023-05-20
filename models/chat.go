package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	ChatTypeText = "text"
	ChatTypePoll = "poll"
	ChatTypeInfo = "info"
)

type Chat struct {
	ObjectId primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SenderId string             `bson:"sender_id" json:"sender_id"`
	// ChatId is user id if personal chat, or group id if group chat
	ChatId    string    `bson:"chat_id" json:"chat_id"`
	IsGroup   bool      `bson:"is_group" json:"is_group"`
	Type      string    `bson:"type" json:"type"`
	Text      string    `bson:"text,omitempty" json:"text,omitempty"`
	MediaUrls []string  `bson:"media_urls" json:"media_urls"`
	Poll      *Poll     `bson:"poll,omitempty" json:"poll,omitempty"`
	Info      *ChatInfo `bson:"info,omitempty" json:"info,omitempty"`
	ReadBy    []string  `bson:"read_by" json:"read_by"`
	CreatedAt int64     `bson:"created_at" json:"created_at"`
	UpdatedAt int64     `bson:"updated_at" json:"updated_at"`
}

type Group struct {
	ObjectId primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name" json:"name"`
	// MemberIds is a list of user ids who are members of this group
	MemberIds []string `bson:"member_ids" json:"member_ids"`
	// AdminIds is a list of user ids who are admins of this group
	AdminIds  []string `bson:"admin_ids" json:"admin_ids"`
	CreatedAt int64    `bson:"created_at" json:"created_at"`
	CreatedBy string   `bson:"created_by" json:"created_by"`
	UpdatedAt int64    `bson:"updated_at" json:"updated_at"`
}

type PollOption struct {
	Text string `bson:"text" json:"text"`
	// UserVotes is a list of user ids who voted this option
	UserVotes []string `bson:"user_votes" json:"user_votes"`
}

type Poll struct {
	Question string        `bson:"question" json:"question"`
	Options  []*PollOption `bson:"options" json:"options"`
}

// ChatInfo is for 'notifications' on group
// e.g. user joined, user left, group created, image changed, etc
type ChatInfo struct {
	Type string `bson:"type" json:"type"`
	// UserId is the user who did the action
	UserId    string `bson:"user_id" json:"user_id"`
	Message   string `bson:"message" json:"message"`
	Timestamp int64  `bson:"timestamp" json:"timestamp"`
}
