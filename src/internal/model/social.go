package model

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	Username  string               `bson:"username"`
	Password  string               `bson:"password"`
	Email     string               `bson:"email"`
	Tags      []string             `bson:"tags"`
	Following []primitive.ObjectID `bson:"following"`
	Followers []primitive.ObjectID `bson:"followers"`
	CreatedAt time.Time            `bson:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt"`
}

type UserDetail struct {
	ID       primitive.ObjectID `bson:"_id"`
	Username string             `bson:"username"`
}

type UserUpdates struct {
	Username  string    `bson:"username,omitempty"`
	Email     string    `bson:"email,omitempty"`
	Tags      []string  `bson:"tags,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt"`
}

type SocialRepo struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewSocialRepo(db *mongo.Database, logger log.Logger) *SocialRepo {
	collection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Log(log.LevelError, "msg", "Failed to connect to collection users", "error", err)
	} else {
		logger.Log(log.LevelInfo, "msg", "Successfully connected to collection users")
	}

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_email_unique"),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_username_unique"),
		},
		{
			Keys:    bson.D{{Key: "tags", Value: 1}},
			Options: options.Index().SetName("idx_tags"),
		},
		{
			Keys: bson.D{
				{Key: "following", Value: 1},
				{Key: "followers", Value: 1},
			},
			Options: options.Index().SetName("idx_social_graph"),
		},
		{
			Keys: bson.D{
				{Key: "createdAt", Value: -1},
				{Key: "updatedAt", Value: -1},
			},
			Options: options.Index().SetName("idx_timestamps"),
		},
	}

	_, err = collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		panic(fmt.Sprintf("Failed to create indexes: %v", err))
	}

	return &SocialRepo{
		collection: collection,
		logger:     logger,
	}
}
