package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func (g *Group) FindById(db *mongo.Database, id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	return db.Collection(GroupCollection).FindOne(context.Background(), bson.M{"_id": objId}).Decode(&g)
}

func (g *Group) Save(db *mongo.Database) error {
	if g.ObjectId.IsZero() {
		_, err := db.Collection(GroupCollection).InsertOne(context.Background(), &g)
		return err
	}
	_, err := db.Collection(GroupCollection).UpdateOne(context.Background(), bson.M{"_id": g.ObjectId}, &g)
	return err
}
