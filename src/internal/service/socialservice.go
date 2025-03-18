package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"src/internal/model"
	"src/internal/repository"
	pb "src/protos/Tipster"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SocialServiceService struct {
	pb.UnimplementedSocialServiceServer
	repo   repository.SocialRepository
	logger log.Logger
}

func NewSocialServiceService(repo repository.SocialRepository, logger log.Logger) *SocialServiceService {
	return &SocialServiceService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SocialServiceService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	currentTime := time.Now().UTC()

	// Check if email already exists
	existingUser := s.repo.CheckEmail(ctx, req.Email)
	if existingUser.Err() != mongo.ErrNoDocuments {
		return &pb.CreateUserResponse{
			Code: "COMM0400",
			Msg:  "Email already exists",
		}, nil
	}

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

	userID, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Fatal("MongoDB Insertion Error:", err)
	}
	fmt.Println("Successfully created user")

	return &pb.CreateUserResponse{
		Code: "COMM0000",
		Msg:  "User created successfully",
		Data: &pb.CreateUserResponse_UserData{
			UserId:    userID,
			UserName:  req.UserName,
			Email:     req.Email,
			Tags:      req.Tags,
			CreatedAt: timestamppb.New(currentTime),
			UpdatedAt: timestamppb.New(currentTime),
		},
	}, nil
}
func (s *SocialServiceService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.GetUserResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	user, err := s.repo.GetUser(ctx, objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "User not found")
		}
		return nil, status.Errorf(codes.Internal, "Error fetching user: %v", err)
	}

	// Get follower details
	followerDetails, err := s.repo.GetUserDetails(ctx, user.Followers)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error fetching followers: %v", err)
	}

	// Get following details
	followingDetails, err := s.repo.GetUserDetails(ctx, user.Following)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error fetching followings: %v", err)
	}

	// Convert model.UserDetail to pb.UserDetail
	pbFollowers := make([]*pb.UserDetail, len(followerDetails))
	for i, f := range followerDetails {
		pbFollowers[i] = &pb.UserDetail{
			Id:       strings.TrimSpace(f.ID.Hex()),
			UserName: strings.TrimSpace(f.Username),
		}
	}

	pbFollowings := make([]*pb.UserDetail, len(followingDetails))
	for i, f := range followingDetails {
		pbFollowings[i] = &pb.UserDetail{
			Id:       strings.TrimSpace(f.ID.Hex()),
			UserName: strings.TrimSpace(f.Username),
		}
	}

	return &pb.GetUserResponse{
		Code: "COMM0000",
		Msg:  "User found successfully",
		Data: &pb.CreateUserResponse_UserData{
			UserId:     user.ID.Hex(),
			UserName:   user.Username,
			Email:      user.Email,
			Tags:       user.Tags,
			CreatedAt:  timestamppb.New(user.CreatedAt),
			UpdatedAt:  timestamppb.New(user.UpdatedAt),
			Followers:  pbFollowers,
			Followings: pbFollowings,
		},
	}, nil
}
func (s *SocialServiceService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// Validate user ID
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.UpdateUserResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Create updates struct
	updates := &model.UserUpdates{
		Username:  req.UserName,
		Email:     req.Email,
		Tags:      req.Tags,
		UpdatedAt: time.Now().UTC(),
	}

	// Attempt to update the user
	err = s.repo.UpdateUser(ctx, objID, updates)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UpdateUserResponse{
				Code: "COMM0201",
				Msg:  "User not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to update user", "error", err)
		return &pb.UpdateUserResponse{
			Code: "COMM0501",
			Msg:  "Failed to update user",
		}, nil
	}

	return &pb.UpdateUserResponse{
		Code: "COMM0000",
		Msg:  "User updated successfully",
	}, nil
}
func (s *SocialServiceService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	// Validate user ID
	objID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.DeleteUserResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to delete the user
	err = s.repo.DeleteUser(ctx, objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.DeleteUserResponse{
				Code: "COMM0300",
				Msg:  "User not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to delete user", "error", err)
		return &pb.DeleteUserResponse{
			Code: "COMM0501",
			Msg:  "Database error while deleting user",
		}, nil
	}

	s.logger.Log(log.LevelInfo, "user deleted successfully", "user_id", req.UserId)
	return &pb.DeleteUserResponse{
		Code: "COMM0000",
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
				Code: "COMM0101",
				Msg:  "Invalid cursor format",
			}, nil
		}
		filter["_id"] = bson.M{"$gt": cursorID}
	}

	// Fetch users
	users, err := s.repo.ListUsers(ctx, filter, int64(req.PageSize))
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch users", "error", err)
		return &pb.ListUserResponse{
			Code: "COMM0501",
			Msg:  "Failed to fetch users",
		}, nil
	}

	var pbUsers []*pb.CreateUserResponse_UserData
	var lastUserID primitive.ObjectID

	for _, user := range users {
		lastUserID = user.ID

		// Get followers' details
		followerDetails, err := s.repo.GetUserDetails(ctx, user.Followers)
		if err != nil {
			s.logger.Log(log.LevelError, "failed to fetch follower details", "error", err)
			continue
		}

		// Get following details
		followingDetails, err := s.repo.GetUserDetails(ctx, user.Following)
		if err != nil {
			s.logger.Log(log.LevelError, "failed to fetch following details", "error", err)
			continue
		}

		pbFollowers := make([]*pb.UserDetail, len(followerDetails))
		for i, f := range followerDetails {
			pbFollowers[i] = &pb.UserDetail{
				Id:       strings.TrimSpace(f.ID.Hex()),
				UserName: strings.TrimSpace(f.Username),
			}
		}

		pbFollowings := make([]*pb.UserDetail, len(followingDetails))
		for i, f := range followingDetails {
			pbFollowings[i] = &pb.UserDetail{
				Id:       strings.TrimSpace(f.ID.Hex()),
				UserName: strings.TrimSpace(f.Username),
			}
		}

		userData := &pb.CreateUserResponse_UserData{
			UserId:     user.ID.Hex(),
			UserName:   user.Username,
			Email:      user.Email,
			Tags:       user.Tags,
			CreatedAt:  timestamppb.New(user.CreatedAt),
			UpdatedAt:  timestamppb.New(user.UpdatedAt),
			Followers:  pbFollowers,
			Followings: pbFollowings,
		}
		pbUsers = append(pbUsers, userData)
	}

	// Determine `NextCursor`
	nextCursor := ""
	if len(pbUsers) == int(req.PageSize) {
		nextCursor = lastUserID.Hex()
	}

	return &pb.ListUserResponse{
		Code: "COMM0000",
		Msg:  "Users retrieved successfully",
		Data: &pb.ListUserResponse_ListUsersData{
			Users:      pbUsers,
			NextCursor: nextCursor,
		},
	}, nil
}
func (s *SocialServiceService) FollowTipster(ctx context.Context, req *pb.FollowTipsterRequest) (*pb.FollowTipsterResponse, error) {
	// Validate IDs
	tipsterID, err := primitive.ObjectIDFromHex(req.TipsterId)
	if err != nil {
		return &pb.FollowTipsterResponse{
			Code: "COMM0101",
			Msg:  "Invalid tipster ID format",
		}, nil
	}
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.FollowTipsterResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to follow the tipster
	err = s.repo.FollowTipster(ctx, userID, tipsterID)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to update follow relationship", "error", err)
		return &pb.FollowTipsterResponse{
			Code: "COMM0501",
			Msg:  "Failed to update follow relationship",
		}, nil
	}

	return &pb.FollowTipsterResponse{
		Code: "COMM0000",
		Msg:  "Successfully followed tipster",
		Data: &pb.FollowTipsterResponse_FollowData{
			IsFollowing: true,
			TipsterId:   req.TipsterId,
			UserId:      req.UserId,
		},
	}, nil
}
func (s *SocialServiceService) UnfollowTipster(ctx context.Context, req *pb.UnFollowTipsterRequest) (*pb.UnfollowTipsterResponse, error) {
	// Validate IDs
	tipsterID, err := primitive.ObjectIDFromHex(req.TipsterId)
	if err != nil {
		return &pb.UnfollowTipsterResponse{
			Code: "COMM0101",
			Msg:  "Invalid tipster ID format",
		}, nil
	}
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.UnfollowTipsterResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to unfollow the tipster
	err = s.repo.UnfollowTipster(ctx, userID, tipsterID)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to update unfollow relationship", "error", err)
		return &pb.UnfollowTipsterResponse{
			Code: "COMM0501",
			Msg:  "Failed to update unfollow relationship",
		}, nil
	}

	return &pb.UnfollowTipsterResponse{
		Code: "COMM0000",
		Msg:  "Successfully unfollowed tipster",
		Data: &pb.UnfollowTipsterResponse_UnfollowData{
			IsFollowing: false,
			TipsterId:   req.TipsterId,
			UserId:      req.UserId,
		},
	}, nil
}
