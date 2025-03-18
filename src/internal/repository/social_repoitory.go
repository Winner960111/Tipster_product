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
	ListTip(ctx context.Context, tipID string) ([]*model.Comment, error)
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
func (r *socialRepository) CreateTip(ctx context.Context, tip *model.Tip) (primitive.ObjectID, error) {

	result, err := r.tipCollection.InsertOne(ctx, tip)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}
func (r *socialRepository) GetTip(ctx context.Context, tipID primitive.ObjectID) (*model.Tip, error) {
	var tip model.Tip
	err := r.tipCollection.FindOne(ctx, bson.M{"_id": tipID}).Decode(&tip)
	if err != nil {
		return nil, err
	}
	return &tip, nil
}
func (r *socialRepository) UpdateTip(ctx context.Context, tipID primitive.ObjectID, updates bson.M) (*model.Tip, error) {
	var updatedTip model.Tip
	err := r.tipCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": tipID},
		updates,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedTip)

	if err != nil {
		return nil, err
	}
	return &updatedTip, nil
}
func (r *socialRepository) DeleteTip(ctx context.Context, tipID primitive.ObjectID) error {
	result, err := r.tipCollection.DeleteOne(ctx, bson.M{"_id": tipID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *socialRepository) ListTips(ctx context.Context, filter bson.M, pageSize int64, nextCursor string) ([]*model.Tip, primitive.ObjectID, error) {
	if nextCursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(nextCursor)
		if err != nil {
			return nil, primitive.NilObjectID, err
		}
		filter["_id"] = bson.M{"$gt": cursorID}
	}

	findOptions := options.Find().SetLimit(pageSize).SetSort(bson.M{"_id": 1})
	cursor, err := r.tipCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, primitive.NilObjectID, err
	}
	defer cursor.Close(ctx)

	var tips []*model.Tip
	var lastTipID primitive.ObjectID
	for cursor.Next(ctx) {
		var tip model.Tip
		if err := cursor.Decode(&tip); err == nil {
			tips = append(tips, &tip)
			lastTipID = tip.ID
		}
	}
	return tips, lastTipID, cursor.Err()
}

func (r *socialRepository) LikeTip(ctx context.Context, tipID, userID primitive.ObjectID) (int32, error) {
	// Check if user already liked the tip
	var existingTip struct {
		Likes []primitive.ObjectID `bson:"likes"`
	}
	err := r.tipCollection.FindOne(ctx, bson.M{"_id": tipID, "likes": userID}).Decode(&existingTip)
	if err == nil {
		// User already liked this tip
		return int32(len(existingTip.Likes)), nil
	}

	// Add userId to `likes` array and remove from `unlikes`
	update := bson.M{
		"$pull": bson.M{
			"unlikes": userID, // Remove user from `unlikes` if they previously disliked
		},
		"$addToSet": bson.M{
			"likes": userID, // Ensure userId is only added once
		},
	}

	// Get updated document
	var updatedTip struct {
		Likes []primitive.ObjectID `bson:"likes"`
	}
	err = r.tipCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": tipID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedTip)

	if err != nil {
		return 0, err
	}

	// Compute total likes safely
	return int32(len(updatedTip.Likes)), nil
}
func (r *socialRepository) UnlikeTip(ctx context.Context, tipID, userID primitive.ObjectID) (int32, error) {
	// Check if user already unliked the tip
	var existingTip struct {
		Unlikes []primitive.ObjectID `bson:"unlikes"`
	}
	err := r.tipCollection.FindOne(ctx, bson.M{"_id": tipID, "unlikes": userID}).Decode(&existingTip)
	if err == nil {
		// User already unliked this tip
		return int32(len(existingTip.Unlikes)), nil
	}

	// Add unlike and remove like if exists
	update := bson.M{
		"$pull": bson.M{
			"likes": userID, // Remove user from `likes` if they previously liked
		},
		"$addToSet": bson.M{
			"unlikes": userID, // Ensure userId is only added once
		},
	}

	// Get updated document
	var updatedTip struct {
		Unlikes []primitive.ObjectID `bson:"unlikes"`
	}
	err = r.tipCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": tipID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedTip)

	if err != nil {
		return 0, err
	}

	// Compute total unlikes safely
	return int32(len(updatedTip.Unlikes)), nil
}
func (r *socialRepository) ShareTip(ctx context.Context, tipID primitive.ObjectID, shareType string, updatedAt time.Time) error {
	_, err := r.tipCollection.UpdateOne(
		ctx,
		bson.M{"_id": tipID},
		bson.M{
			"$set": bson.M{
				"shareType": shareType,
				"updatedAt": updatedAt,
			},
		},
	)
	return err
}

func (r *socialRepository) CreateComment(ctx context.Context, comment *model.Comment) (primitive.ObjectID, error) {
	result, err := r.commentCollection.InsertOne(ctx, comment)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *socialRepository) UpdateComment(ctx context.Context, commentID primitive.ObjectID, content string, updatedAt time.Time) error {
	update := bson.M{
		"$set": bson.M{
			"content":   content,
			"updatedAt": updatedAt,
		},
	}

	result := r.commentCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": commentID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (r *socialRepository) DeleteComment(ctx context.Context, commentID primitive.ObjectID) error {
	result, err := r.commentCollection.DeleteOne(ctx, bson.M{"_id": commentID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *socialRepository) ListTip(ctx context.Context, tipID string) ([]*model.Comment, error) {
	cursor, err := r.commentCollection.Find(ctx, bson.M{"tipId": tipID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var comments []*model.Comment
	for cursor.Next(ctx) {
		var comment model.Comment
		if err := cursor.Decode(&comment); err == nil {
			comments = append(comments, &comment)
		}
	}
	return comments, cursor.Err()
}

func (r *socialRepository) LikeComment(ctx context.Context, commentID, userID primitive.ObjectID) (int32, error) {

	// Check if user already liked the comment
	var existingComment struct {
		Likes []primitive.ObjectID `bson:"likes"`
	}
	err := r.commentCollection.FindOne(ctx, bson.M{"_id": commentID, "likes": userID}).Decode(&existingComment)
	if err == nil {
		// User already liked this comment
		return int32(len(existingComment.Likes)), nil
	}

	// Add userId to `likes` array and remove from `unlikes`
	update := bson.M{
		"$pull": bson.M{
			"unlikes": userID, // Remove user from `unlikes` if they previously disliked
		},
		"$addToSet": bson.M{
			"likes": userID, // Ensure userId is only added once
		},
	}

	// Get updated document
	var updatedComment struct {
		Likes []primitive.ObjectID `bson:"likes"`
	}
	err = r.commentCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": commentID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedComment)

	if err != nil {
		return 0, err
	}

	// Compute total likes safely
	return int32(len(updatedComment.Likes)), nil
}

func (r *socialRepository) UnlikeComment(ctx context.Context, commentID, userID primitive.ObjectID) (int32, error) {

	// Check if user already unliked the comment
	var existingComment struct {
		Unlikes []primitive.ObjectID `bson:"unlikes"`
	}
	err := r.commentCollection.FindOne(ctx, bson.M{"_id": commentID, "unlikes": userID}).Decode(&existingComment)
	if err == nil {
		// User already unliked this comment
		return int32(len(existingComment.Unlikes)), nil
	}

	// Add userId to `unlikes` array and remove from `likes`
	update := bson.M{
		"$pull": bson.M{
			"likes": userID, // Remove user from `likes` if they previously liked
		},
		"$addToSet": bson.M{
			"unlikes": userID, // Ensure userId is only added once
		},
	}

	// Get updated document
	var updatedComment struct {
		Unlikes []primitive.ObjectID `bson:"unlikes"`
	}
	err = r.commentCollection.FindOneAndUpdate(
		ctx,
		bson.M{"_id": commentID},
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updatedComment)

	if err != nil {
		return 0, err
	}

	// Compute total unlikes safely
	return int32(len(updatedComment.Unlikes)), nil
}

func (r *socialRepository) CreateReply(ctx context.Context, reply *model.Comment) (primitive.ObjectID, error) {
	if reply.ID.IsZero() {
		reply.ID = primitive.NewObjectID() // Ensure the reply has a valid ObjectID
	}
	result, err := r.commentCollection.InsertOne(ctx, reply)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *socialRepository) ListReplies(ctx context.Context, parentCommentID string) ([]*model.Comment, error) {
	cursor, err := r.commentCollection.Find(ctx, bson.M{"parentId": parentCommentID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var replies []*model.Comment
	for cursor.Next(ctx) {
		var reply model.Comment
		if err := cursor.Decode(&reply); err == nil {
			replies = append(replies, &reply)
		}
	}
	return replies, cursor.Err()
}

func (r *socialRepository) ListComments(ctx context.Context, pageSize int64, nextCursor string) ([]*model.Comment, string, error) {

	findOptions := options.Find()
	findOptions.SetLimit(pageSize)
	findOptions.SetSort(bson.M{"_id": 1})

	filter := bson.M{}
	if nextCursor != "" {
		lastID, err := primitive.ObjectIDFromHex(nextCursor)
		if err != nil {
			return nil, "", err
		}
		filter["_id"] = bson.M{"$gt": lastID}
	}

	cursor, err := r.commentCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, "", err
	}
	defer cursor.Close(ctx)

	var comments []*model.Comment
	var lastCommentID primitive.ObjectID

	for cursor.Next(ctx) {
		var comment model.Comment
		if err := cursor.Decode(&comment); err != nil {
			continue
		}
		comments = append(comments, &comment)
		lastCommentID = comment.ID
	}

	nextCursor = ""
	if len(comments) == int(pageSize) {
		nextCursor = lastCommentID.Hex()
	}

	return comments, nextCursor, nil
}
