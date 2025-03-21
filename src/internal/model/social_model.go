package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model
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

// Tip model
type Tip struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	TipsterID string               `bson:"tipsterId"`
	Title     string               `bson:"title"`
	Content   string               `bson:"content"`
	Tags      []string             `bson:"tags"`
	ShareType string               `bson:"shareType"`
	Likes     []primitive.ObjectID `bson:"likes"`
	Unlikes   []primitive.ObjectID `bson:"unlikes"`
	CreatedAt time.Time            `bson:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt"`
}

// Comment model
type Comment struct {
	ID        primitive.ObjectID   `bson:"_id"`
	TipID     string               `bson:"tipId"`
	UserID    string               `bson:"userId"`
	ParentID  string               `bson:"parentId"`
	Content   string               `bson:"content"`
	Likes     []primitive.ObjectID `bson:"likes"`
	Unlikes   []primitive.ObjectID `bson:"unlikes"`
	CreatedAt time.Time            `bson:"createdAt"`
	UpdatedAt time.Time            `bson:"updatedAt"`
}
