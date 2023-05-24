package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	ChatTypeText    = "text"
	ChatTypePoll    = "poll"
	ChatTypeInfo    = "info"
	ChatCollection  = "chats"
	GroupCollection = "groups"
)

type ChatRoom struct {
	ChatId string      `json:"chat_id" bson:"chat_id"`
	User   interface{} `json:"user" bson:"user"`
	Group  *Group      `json:"group" bson:"group"`
	Chats  []*Chat     `json:"chats" bson:"chats"`
}

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

func (c *Chat) Save(db *mongo.Database) error {
	if c.ObjectId.IsZero() {
		r, err := db.Collection(ChatCollection).InsertOne(context.Background(), &c)
		if err != nil {
			return err
		}
		oid, ok := r.InsertedID.(primitive.ObjectID)
		if ok {
			c.ObjectId = oid
		}
		return nil
	}
	_, err := db.Collection(ChatCollection).UpdateOne(context.Background(), bson.M{"_id": c.ObjectId}, &c)
	return err
}

func GetUserChatRooms(db *mongo.Database, userId string) ([]*ChatRoom, error) {
	ctx := context.Background()
	rooms := map[string]*ChatRoom{}
	chatRooms := []*ChatRoom{}
	cursor, err := db.Collection(ChatCollection).Find(
		ctx,
		bson.M{"$or": []bson.M{
			{"sender_id": userId},
			{"chat_id": userId},
		}},
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var chat Chat
		if err := cursor.Decode(&chat); err != nil {
			return nil, err
		}
		// Find chat room
		room, ok := rooms[chat.ChatId]
		if !ok {
			// Create one
			room = &ChatRoom{
				ChatId: chat.ChatId,
				Chats:  make([]*Chat, 0),
			}
			rooms[chat.ChatId] = room
			if chat.IsGroup {
				// Find group
				group := &Group{}
				if err = group.FindById(db, chat.ChatId); err != nil {
					continue
				}
				room.Group = group
			}
			chatRooms = append(chatRooms, room)
		}
		room.Chats = append(room.Chats, &chat)
	}

	return chatRooms, nil
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
