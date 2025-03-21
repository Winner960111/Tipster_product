package service

import (
	"context"

	"src/internal/biz"
	"src/internal/errors"
	"src/internal/repository"

	pb "src/protos/Tipster"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	CodeOk          = "COMM0000"
	CodeInvalidID   = "COMM0101"
	CodeInvalidData = "COMM0201"
	CodeEmailExist  = "COMM0400"
	CodeError       = "COMM0501"
)

type SocialServiceService struct {
	pb.UnimplementedSocialServiceServer
	repo   repository.SocialRepository
	logger log.Logger
	biz    *biz.SocialService
}

func NewSocialServiceService(repo repository.SocialRepository, logger log.Logger) *SocialServiceService {
	return &SocialServiceService{
		repo:   repo,
		logger: logger,
		biz:    &biz.SocialService{Repo: repo},
	}
}

func (s *SocialServiceService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	data := s.biz.CreateUser(ctx, req)
	if data.UserId == "" {
		return nil, errors.ErrEmailAlreadyExists
	}

	return &pb.CreateUserResponse{
		Code: CodeOk,
		Msg:  "User created successfully",
		Data: data,
	}, nil
}

func (s *SocialServiceService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.GetUserResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid user ID format",
		}, nil
	}

	user, err := s.biz.GetUser(ctx, objID)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}

	data, err := s.Transformer(ctx, user)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}
	return &pb.GetUserResponse{
		Code: CodeOk,
		Msg:  "User found successfully",
		Data: data,
	}, nil
}

func (s *SocialServiceService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// Validate user ID
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.UpdateUserResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid user ID format",
		}, nil
	}
	// Attempt to update the user
	err = s.biz.UpdateUser(ctx, objID, req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UpdateUserResponse{
				Code: CodeInvalidData,
				Msg:  "User not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to update user", "error", err)
		return &pb.UpdateUserResponse{
			Code: CodeError,
			Msg:  "Failed to update user",
		}, nil
	}

	return &pb.UpdateUserResponse{
		Code: CodeOk,
		Msg:  "User updated successfully",
	}, nil
}

func (s *SocialServiceService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Validate user ID
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.DeleteUserResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to delete the user
	err = s.biz.DeleteUser(ctx, objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.DeleteUserResponse{
				Code: CodeInvalidData,
				Msg:  "User not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to delete user", "error", err)
		return &pb.DeleteUserResponse{
			Code: CodeError,
			Msg:  "Database error while deleting user",
		}, nil
	}

	s.logger.Log(log.LevelInfo, "user deleted successfully", "user_id", req.UserId)
	return &pb.DeleteUserResponse{
		Code: CodeOk,
		Msg:  "User deleted successfully",
	}, nil
}

func (s *SocialServiceService) ListUsers(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	// Pagination settings
	filter := bson.M{}
	if req.NextCursor != "" {
		cursorID, err := primitive.ObjectIDFromHex(req.NextCursor)
		if err != nil {
			return &pb.ListUserResponse{
				Code: CodeInvalidID,
				Msg:  "Invalid cursor format",
			}, nil
		}
		filter["_id"] = bson.M{"$gt": cursorID}
	}

	// Fetch users
	users, err := s.biz.ListUsers(ctx, filter, int64(req.PageSize))
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch users", "error", err)
		return &pb.ListUserResponse{
			Code: CodeError,
			Msg:  "Failed to fetch users",
		}, nil
	}

	return &pb.ListUserResponse{
		Code: CodeOk,
		Msg:  "Users retrieved successfully",
		Data: &pb.ListUserResponse_ListUsersData{
			Users:      pbUsers,
			NextCursor: nextCursor,
		},
	}, nil
}
func (s *SocialServiceService) FollowTipster(ctx context.Context, req *pb.FollowTipsterRequest) (*pb.FollowTipsterResponse, error) {
	tipsterID, err := primitive.ObjectIDFromHex(req.TipsterId)
	if err != nil {
		return &pb.FollowTipsterResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tipster ID format",
		}, nil
	}
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.FollowTipsterResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to follow the tipster
	data, err := s.biz.FollowTipster(ctx, userID, tipsterID)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to update follow relationship", "error", err)
		return &pb.FollowTipsterResponse{
			Code: CodeError,
			Msg:  "Failed to update follow relationship",
		}, nil
	}
	return &pb.FollowTipsterResponse{
		Code: CodeOk,
		Msg:  "Successfully followed tipster",
		Data: data,
	}, nil
}

func (s *SocialServiceService) UnfollowTipster(ctx context.Context, req *pb.UnFollowTipsterRequest) (*pb.UnfollowTipsterResponse, error) {
	// Validate IDs
	tipsterID, err := primitive.ObjectIDFromHex(req.TipsterId)
	if err != nil {
		return &pb.UnfollowTipsterResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tipster ID format",
		}, nil
	}
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.UnfollowTipsterResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to unfollow the tipster
	data, err := s.biz.UnfollowTipster(ctx, userID, tipsterID)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to update unfollow relationship", "error", err)
		return &pb.UnfollowTipsterResponse{
			Code: CodeError,
			Msg:  "Failed to update unfollow relationship",
		}, nil
	}

	return &pb.UnfollowTipsterResponse{
		Code: CodeOk,
		Msg:  "Successfully unfollowed tipster",
		Data: data,
	}, nil
}
