package biz

import (
	"context"
	"src/internal/errors"
	"src/internal/model"
	pb "src/protos/Tipster"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialService) CreateTip(ctx context.Context, req *pb.CreateTipRequest) (*pb.TipData, error) {
	currentTime := time.Now().UTC()

	tip := &model.Tip{
		TipsterID: req.TipsterId,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		ShareType: req.ShareType,
		Likes:     []primitive.ObjectID{},
		Unlikes:   []primitive.ObjectID{},
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	tipID, err := s.Repo.CreateTip(ctx, tip)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}

	return &pb.TipData{
		TipId:     tipID.Hex(),
		TipsterId: req.TipsterId,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		CreatedAt: timestamppb.New(currentTime),
		UpdatedAt: timestamppb.New(currentTime),
		ShareType: req.ShareType,
	}, nil
}

func (s *SocialService) GetTip(ctx context.Context, tipID primitive.ObjectID) (*model.Tip, error) {
	tip, err := s.Repo.GetTip(ctx, tipID)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}
	return tip, nil
}

func (s *SocialService) UpdateTip(ctx context.Context, tipID primitive.ObjectID, req *pb.UpdateTipRequest) (*model.Tip, error) {
	currentTime := time.Now().UTC()
	updates := bson.M{
		"$set": bson.M{
			"title":     req.Title,
			"content":   req.Content,
			"tags":      req.Tags,
			"shareType": req.ShareType,
			"updatedAt": currentTime,
		},
	}
	tip, err := s.Repo.UpdateTip(ctx, tipID, updates)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}
	return tip, nil
}

func (s *SocialService) DeleteTip(ctx context.Context, tipID primitive.ObjectID) error {
	err := s.Repo.DeleteTip(ctx, tipID)
	if err != nil {
		return errors.ToRpcError(err)
	}
	return nil
}

func (s *SocialService) ListTips(ctx context.Context, req *pb.ListTipsRequest) ([]*model.Tip, primitive.ObjectID, error) {

	filter := bson.M{}
	if req.TipsterId != "" {
		filter["tipsterId"] = req.TipsterId
	}

	tips, lastTipID, err := s.Repo.ListTips(ctx, filter, int64(req.PageSize), req.NextCursor)
	if err != nil {
		return nil, primitive.NilObjectID, errors.ToRpcError(err)
	}
	return tips, lastTipID, nil
}

func (s *SocialService) ShareTip(ctx context.Context, tipID primitive.ObjectID, shareType string) error {
	currentTime := time.Now().UTC()
	err := s.Repo.ShareTip(ctx, tipID, shareType, currentTime)
	if err != nil {
		return errors.ToRpcError(err)
	}
	return nil
}

func (s *SocialService) LikeTip(ctx context.Context, tipID primitive.ObjectID, userID primitive.ObjectID) (int32, error) {
	totalLikes, err := s.Repo.LikeTip(ctx, tipID, userID)
	if err != nil {
		return 0, errors.ToRpcError(err)
	}
	return totalLikes, nil
}

func (s *SocialService) UnlikeTip(ctx context.Context, tipID primitive.ObjectID, userID primitive.ObjectID) (int32, error) {
	totalUnlikes, err := s.Repo.UnlikeTip(ctx, tipID, userID)
	if err != nil {
		return 0, errors.ToRpcError(err)
	}
	return totalUnlikes, nil
}
