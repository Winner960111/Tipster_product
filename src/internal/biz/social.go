package biz

import (
	"context"
	"fmt"
	"time"

	"src/internal/errors"
	"src/internal/model"
	"src/internal/repository"
	pb "src/protos/Tipster"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SocialService struct {
	Repo repository.SocialRepository
}

func (s *SocialService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) *pb.CreateUserResponse_UserData {
	currentTime := time.Now().UTC()

	user := &model.User{
		ID:        primitive.NewObjectID(),
		Username:  req.UserName,
		Password:  req.Password,
		Email:     req.Email,
		Tags:      req.Tags,
		Following: []primitive.ObjectID{},
		Followers: []primitive.ObjectID{},
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	userID, err := s.Repo.CreateUser(ctx, user)
	if err != nil {
		fmt.Println(errors.ToRpcError(err))
	}
	if userID == "exist" {
		return &pb.CreateUserResponse_UserData{
			UserId: "",
			Email:  req.Email,
		}
	}
	return &pb.CreateUserResponse_UserData{
		UserId:    userID,
		UserName:  req.UserName,
		Email:     req.Email,
		Tags:      req.Tags,
		CreatedAt: timestamppb.New(currentTime),
		UpdatedAt: timestamppb.New(currentTime),
	}
}

func (s *SocialService) GetUser(ctx context.Context, userID primitive.ObjectID) (*model.User, error) {
	user, err := s.Repo.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}
	return user, nil
}

func (s *SocialService) UpdateUser(ctx context.Context, userID primitive.ObjectID, req *pb.UpdateUserRequest) error {
	// Create updates struct
	updates := &model.UserUpdates{
		Username:  req.UserName,
		Email:     req.Email,
		Tags:      req.Tags,
		UpdatedAt: time.Now().UTC(),
	}
	err := s.Repo.UpdateUser(ctx, userID, updates)
	if err != nil {
		return errors.ToRpcError(err)
	}
	return nil
}

func (s *SocialService) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	err := s.Repo.DeleteUser(ctx, userID)
	if err != nil {
		return errors.ToRpcError(err)
	}
	return nil
}
func (s *SocialService) ListUsers(ctx context.Context, req *pb.ListUserRequest) ([]*model.User, error) {
	filter := bson.M{}
	if req.NextCursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(req.NextCursor)
		if err != nil {
			return nil, errors.ToRpcError(err)
		}
		filter["_id"] = bson.M{"$gt": cursorID}
	}

	users, err := s.Repo.ListUsers(ctx, filter, int64(req.PageSize))
	if err != nil {
		return nil, errors.ToRpcError(err)
	}
	return users, nil
}

func (s *SocialService) FollowTipster(ctx context.Context, userID, tipsterID primitive.ObjectID) (*pb.FollowTipsterResponse_FollowData, error) {
	err := s.Repo.FollowTipster(ctx, userID, tipsterID)
	if err != nil {
		return nil, err
	}
	return &pb.FollowTipsterResponse_FollowData{
		IsFollowing: true,
		TipsterId:   tipsterID.Hex(),
		UserId:      userID.Hex(),
	}, nil
}

func (s *SocialService) UnfollowTipster(ctx context.Context, userID, tipsterID primitive.ObjectID) (*pb.UnfollowTipsterResponse_UnfollowData, error) {
	err := s.Repo.UnfollowTipster(ctx, userID, tipsterID)
	if err != nil {
		return nil, err
	}
	return &pb.UnfollowTipsterResponse_UnfollowData{
		IsFollowing: false,
		TipsterId:   tipsterID.Hex(),
		UserId:      userID.Hex(),
	}, nil
}
