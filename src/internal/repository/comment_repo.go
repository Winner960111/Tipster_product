package repository

import (
	"context"
	"src/internal/errors"
	"src/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
		return errors.ToRpcError(err)
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
func (r *socialRepository) ListTipComments(ctx context.Context, tipID string) ([]*model.Comment, error) {
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
			"unlikes": userID,
		},
		"$addToSet": bson.M{
			"likes": userID,
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
		return 0, errors.ToRpcError(err)
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
			"likes": userID,
		},
		"$addToSet": bson.M{
			"unlikes": userID,
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
		return 0, errors.ToRpcError(err)
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
		return primitive.NilObjectID, errors.ToRpcError(err)
	}
	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *socialRepository) ListReplies(ctx context.Context, parentCommentID string) ([]*model.Comment, error) {
	cursor, err := r.commentCollection.Find(ctx, bson.M{"parentId": parentCommentID})
	if err != nil {
		return nil, errors.ToRpcError(err)
	}
	defer cursor.Close(ctx)

	var replies []*model.Comment
	for cursor.Next(ctx) {
		var reply model.Comment
		if err := cursor.Decode(&reply); err == nil {
			if parentCommentID != reply.ID.Hex() {
				replies = append(replies, &reply)
			}
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
