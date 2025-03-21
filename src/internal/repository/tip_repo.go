package repository

import (
	"context"
	"src/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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
