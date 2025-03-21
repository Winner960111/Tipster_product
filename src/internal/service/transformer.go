package service

import (
	"context"
	"src/internal/model"
	pb "src/protos/Tipster"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialServiceService) Transformer(ctx context.Context, user *model.User) (*pb.CreateUserResponse_UserData, error) {

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

	return &pb.CreateUserResponse_UserData{
		UserId:     user.ID.Hex(),
		UserName:   user.Username,
		Email:      user.Email,
		Tags:       user.Tags,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
		Followers:  pbFollowers,
		Followings: pbFollowings,
	}, nil
}
