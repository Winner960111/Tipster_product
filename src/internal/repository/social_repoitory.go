package repository

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"src/internal/model"
)

type SocialRepository interface {
	CreateUser(ctx context.Context, user *model.User) (string, error)
	CheckEmail(ctx context.Context, email string) *mongo.SingleResult
	GetUser(ctx context.Context, userID primitive.ObjectID) (*model.User, error)
	GetUserDetails(ctx context.Context, userIDs []primitive.ObjectID) ([]*model.UserDetail, error)
	UpdateUser(ctx context.Context, userID primitive.ObjectID, updates *model.UserUpdates) error
	DeleteUser(ctx context.Context, userID primitive.ObjectID) error
	ListUsers(ctx context.Context, filter bson.M, limit int64) ([]*model.User, error)
	FollowTipster(ctx context.Context, userID, tipsterID primitive.ObjectID) error
	UnfollowTipster(ctx context.Context, userID, tipsterID primitive.ObjectID) error
}

type socialRepository struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewSocialRepository(db *mongo.Database, logger log.Logger) SocialRepository {
	collection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Log(log.LevelError, "msg", "Failed to connect to collection users", "error", err)
	} else {
		logger.Log(log.LevelInfo, "msg", "Successfully connected to collection users")
	}
	return &socialRepository{
		collection: collection,
		logger:     logger,
	}
}

func (r *socialRepository) CreateUser(ctx context.Context, user *model.User) (string, error) {
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}
	return user.ID.Hex(), nil
}

func (r *socialRepository) CheckEmail(ctx context.Context, email string) *mongo.SingleResult {
	filter := bson.M{"email": email}
	return r.collection.FindOne(ctx, filter)
}

func (r *socialRepository) GetUser(ctx context.Context, userID primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
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
		return err
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
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		userUpdate := bson.M{
			"$addToSet": bson.M{
				"following": tipsterID,
			},
		}
		tipsterUpdate := bson.M{
			"$addToSet": bson.M{
				"followers": userID,
			},
		}

		_, err := r.collection.UpdateOne(sessCtx, bson.M{"_id": userID}, userUpdate)
		if err != nil {
			return nil, err
		}

		_, err = r.collection.UpdateOne(sessCtx, bson.M{"_id": tipsterID}, tipsterUpdate)
		return nil, err
	})

	return err
}
func (r *socialRepository) UnfollowTipster(ctx context.Context, userID, tipsterID primitive.ObjectID) error {
	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		userUpdate := bson.M{
			"$pull": bson.M{
				"following": tipsterID,
			},
		}
		tipsterUpdate := bson.M{
			"$pull": bson.M{
				"followers": userID,
			},
		}

		_, err := r.collection.UpdateOne(sessCtx, bson.M{"_id": userID}, userUpdate)
		if err != nil {
			return nil, err
		}

		_, err = r.collection.UpdateOne(sessCtx, bson.M{"_id": tipsterID}, tipsterUpdate)
		return nil, err
	})

	return err
}
