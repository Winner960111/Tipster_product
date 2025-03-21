package repository

import (
	"context"
	"time"

	"src/internal/errors"
	"src/internal/model"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SocialRepository interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	GetUser(ctx context.Context, userID primitive.ObjectID) (*model.User, error)
	GetUserDetails(ctx context.Context, userIDs []primitive.ObjectID) ([]*model.UserDetail, error)
	UpdateUser(ctx context.Context, userID primitive.ObjectID, updates *model.UserUpdates) error
	DeleteUser(ctx context.Context, userID primitive.ObjectID) error
	ListUsers(ctx context.Context, filter bson.M, limit int64) ([]*model.User, error)
	FollowTipster(ctx context.Context, userID, tipsterID primitive.ObjectID) error
	UnfollowTipster(ctx context.Context, userID, tipsterID primitive.ObjectID) error
	CreateTip(ctx context.Context, tip *model.Tip) (primitive.ObjectID, error)
	GetTip(ctx context.Context, tipID primitive.ObjectID) (*model.Tip, error)
	UpdateTip(ctx context.Context, tipID primitive.ObjectID, updates bson.M) (*model.Tip, error)
	DeleteTip(ctx context.Context, tipID primitive.ObjectID) error
	ListTips(ctx context.Context, filter bson.M, pageSize int64, nextCursor string) ([]*model.Tip, primitive.ObjectID, error)
	LikeTip(ctx context.Context, tipID, userID primitive.ObjectID) (int32, error)
	UnlikeTip(ctx context.Context, tipID, userID primitive.ObjectID) (int32, error)
	ShareTip(ctx context.Context, tipID primitive.ObjectID, shareType string, updatedAt time.Time) error
	CreateComment(ctx context.Context, comment *model.Comment) (primitive.ObjectID, error)
	UpdateComment(ctx context.Context, commentID primitive.ObjectID, content string, updatedAt time.Time) error
	DeleteComment(ctx context.Context, commentID primitive.ObjectID) error
	ListTipComments(ctx context.Context, tipID string) ([]*model.Comment, error)
	LikeComment(ctx context.Context, commentID, userID primitive.ObjectID) (int32, error)
	UnlikeComment(ctx context.Context, commentID, userID primitive.ObjectID) (int32, error)
	CreateReply(ctx context.Context, reply *model.Comment) (primitive.ObjectID, error)
	ListReplies(ctx context.Context, parentCommentID string) ([]*model.Comment, error)
	ListComments(ctx context.Context, pageSize int64, nextCursor string) ([]*model.Comment, string, error)
}

type socialRepository struct {
	collection        *mongo.Collection
	tipCollection     *mongo.Collection
	commentCollection *mongo.Collection
	logger            log.Logger
}

func NewSocialRepository(db *mongo.Database, logger log.Logger) SocialRepository {
	collection := db.Collection("users")
	tipCollection := db.Collection("tips")
	commentCollection := db.Collection("comments")

	return &socialRepository{
		collection:        collection,
		tipCollection:     tipCollection,
		commentCollection: commentCollection,
		logger:            logger,
	}
}

func (r *socialRepository) CreateUser(ctx context.Context, user *model.User) (string, error) {
	filter := bson.M{"email": user.Email}
	var result model.User
	err := r.collection.FindOne(ctx, filter).Decode(&result)
	if err == nil {
		return "exist", nil
	}
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	_, err = r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	return user.ID.Hex(), nil
}

func (r *socialRepository) GetUser(ctx context.Context, userID primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}
	return &user, nil
}

func (r *socialRepository) GetUserDetails(ctx context.Context, userIDs []primitive.ObjectID) ([]*model.UserDetail, error) {
	if len(userIDs) == 0 {
		return []*model.UserDetail{}, nil
	}

	cursor, err := r.collection.Find(ctx, bson.M{
		"_id": bson.M{"$in": userIDs},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.UserDetail
	for cursor.Next(ctx) {
		var user model.UserDetail
		if err := cursor.Decode(&user); err == nil {
			users = append(users, &user)
		}
	}
	return users, cursor.Err()
}

func (r *socialRepository) UpdateUser(ctx context.Context, userID primitive.ObjectID, updates *model.UserUpdates) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": updates},
	)
	if err != nil {
		return errors.ToRpcError(err)
	}
	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *socialRepository) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": userID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *socialRepository) ListUsers(ctx context.Context, filter bson.M, limit int64) ([]*model.User, error) {
	options := options.Find().SetLimit(limit).SetSort(bson.M{"_id": 1})
	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*model.User
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err == nil {
			users = append(users, &user)
		}
	}
	return users, cursor.Err()
}

func (r *socialRepository) FollowTipster(ctx context.Context, userID, tipsterID primitive.ObjectID) error {
	// Update user's following list
	userUpdate := bson.M{
		"$addToSet": bson.M{
			"following": tipsterID,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, userUpdate)
	if err != nil {
		return err
	}

	// Update tipster's followers list
	tipsterUpdate := bson.M{
		"$addToSet": bson.M{
			"followers": userID,
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": tipsterID}, tipsterUpdate)
	return err
}

func (r *socialRepository) UnfollowTipster(ctx context.Context, userID, tipsterID primitive.ObjectID) error {

	userUpdate := bson.M{
		"$pull": bson.M{
			"following": tipsterID,
		},
	}
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, userUpdate)
	if err != nil {
		return err
	}
	tipsterUpdate := bson.M{
		"$pull": bson.M{
			"followers": userID,
		},
	}
	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": tipsterID}, tipsterUpdate)
	return err
}
