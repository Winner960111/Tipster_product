package service

import (
	"context"
	"src/internal/errors"
	"src/internal/model"
	pb "src/protos/Tipster"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialServiceService) likeUsers(ctx context.Context, likeIDs, unlikeIDs []primitive.ObjectID) (map[string]*pb.UserDetail, map[string]*pb.UserDetail, error) {
	// Fetch user details in bulk
	likeUserDetails, err := s.repo.GetUserDetails(ctx, likeIDs)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch like user details", "error", err)
		return nil, nil, errors.ToRpcError(err)
	}

	unlikeUserDetails, err := s.repo.GetUserDetails(ctx, unlikeIDs)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch unlike user details", "error", err)
		return nil, nil, errors.ToRpcError(err)
	}

	// Map user details for quick lookup
	likeUserMap := map[string]*pb.UserDetail{}
	for _, user := range likeUserDetails {
		likeUserMap[user.ID.Hex()] = &pb.UserDetail{
			Id:       user.ID.Hex(),
			UserName: user.Username,
		}
	}

	unlikeUserMap := map[string]*pb.UserDetail{}
	for _, user := range unlikeUserDetails {
		unlikeUserMap[user.ID.Hex()] = &pb.UserDetail{
			Id:       user.ID.Hex(),
			UserName: user.Username,
		}
	}

	return likeUserMap, unlikeUserMap, nil
}

func (s *SocialServiceService) userTransformer(ctx context.Context, user *model.User) (*pb.CreateUserResponse_UserData, error) {
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
func (s *SocialServiceService) usersTransformer(ctx context.Context, users []*model.User, pageSize int64) (*pb.ListUserResponse_ListUsersData, error) {
	var pbUsers []*pb.CreateUserResponse_UserData
	var lastUserID primitive.ObjectID

	for _, user := range users {
		lastUserID = user.ID
		userData, err := s.userTransformer(ctx, user)
		if err != nil {
			return nil, errors.ToRpcError(err)
		}
		pbUsers = append(pbUsers, userData)
	}

	// Determine `NextCursor`
	nextCursor := ""
	if len(pbUsers) == int(pageSize) {
		nextCursor = lastUserID.Hex()
	}

	return &pb.ListUserResponse_ListUsersData{
		Users:      pbUsers,
		NextCursor: nextCursor,
	}, nil
}

func (s *SocialServiceService) likeTransformer(ctx context.Context, comments []*model.Comment) ([]*pb.CommentInfo, error) {
	var allLikeIDs, allUnlikeIDs []primitive.ObjectID
	for _, comment := range comments {
		allLikeIDs = append(allLikeIDs, comment.Likes...)
		allUnlikeIDs = append(allUnlikeIDs, comment.Unlikes...)
	}

	likeUserMap, unlikeUserMap, err := s.likeUsers(ctx, allLikeIDs, allUnlikeIDs)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}
	// Convert raw comment records into CommentInfo response
	var pbComments []*pb.CommentInfo
	for _, comment := range comments {
		likes := []*pb.UserDetail{}
		for _, likeID := range comment.Likes {
			if user, exists := likeUserMap[likeID.Hex()]; exists {
				likes = append(likes, user)
			}
		}

		unlikes := []*pb.UserDetail{}
		for _, unlikeID := range comment.Unlikes {
			if user, exists := unlikeUserMap[unlikeID.Hex()]; exists {
				unlikes = append(unlikes, user)
			}
		}

		pbComments = append(pbComments, &pb.CommentInfo{
			CommentId: comment.ID.Hex(),
			TipId:     comment.TipID,
			UserId:    comment.UserID,
			ParentId:  comment.ParentID,
			Content:   comment.Content,
			CreatedAt: timestamppb.New(comment.CreatedAt),
			UpdatedAt: timestamppb.New(comment.UpdatedAt),
			Likes:     likes,
			Unlikes:   unlikes,
		})
	}

	return pbComments, nil
}
func (s *SocialServiceService) replyTransformer(ctx context.Context, replies []*model.Comment) ([]*pb.ReplyInfo, error) {
	// Collect all like and unlike user IDs for batch query
	var allLikeIDs, allUnlikeIDs []primitive.ObjectID
	for _, reply := range replies {
		allLikeIDs = append(allLikeIDs, reply.Likes...)
		allUnlikeIDs = append(allUnlikeIDs, reply.Unlikes...)
	}

	likeUserMap, unlikeUserMap, err := s.likeUsers(ctx, allLikeIDs, allUnlikeIDs)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}

	// Convert raw reply records into ReplyInfo response
	var pbReplies []*pb.ReplyInfo
	for _, reply := range replies {
		likes := []*pb.UserDetail{}
		for _, likeID := range reply.Likes {
			if user, exists := likeUserMap[likeID.Hex()]; exists {
				likes = append(likes, user)
			}
		}

		unlikes := []*pb.UserDetail{}
		for _, unlikeID := range reply.Unlikes {
			if user, exists := unlikeUserMap[unlikeID.Hex()]; exists {
				unlikes = append(unlikes, user)
			}
		}

		pbReplies = append(pbReplies, &pb.ReplyInfo{
			ReplyId:         reply.ID.Hex(),
			ParentCommentId: reply.ParentID,
			UserId:          reply.UserID,
			Content:         reply.Content,
			DateCreated:     timestamppb.New(reply.CreatedAt),
			Likes:           likes,
			Unlikes:         unlikes,
		})
	}

	return pbReplies, nil
}

func (s *SocialServiceService) tipTransformer(ctx context.Context, tip *model.Tip) (*pb.TipData, error) {
	likeUserMap, unlikeUserMap, err := s.likeUsers(ctx, tip.Likes, tip.Unlikes)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}

	likes := []*pb.UserDetail{}
	for _, likeID := range tip.Likes {
		if user, exists := likeUserMap[likeID.Hex()]; exists {
			likes = append(likes, user)
		}
	}

	unlikes := []*pb.UserDetail{}
	for _, unlikeID := range tip.Unlikes {
		if user, exists := unlikeUserMap[unlikeID.Hex()]; exists {
			unlikes = append(unlikes, user)
		}
	}

	return &pb.TipData{
		TipId:     tip.ID.Hex(),
		TipsterId: tip.TipsterID,
		Title:     tip.Title,
		Content:   tip.Content,
		Tags:      tip.Tags,
		Likes:     likes,
		Unlikes:   unlikes,
		CreatedAt: timestamppb.New(tip.CreatedAt),
		UpdatedAt: timestamppb.New(tip.UpdatedAt),
		ShareType: tip.ShareType,
	}, nil
}

func (s *SocialServiceService) tipsTransformer(ctx context.Context, tips []*model.Tip, pageSize int64, lastTipID primitive.ObjectID) (*pb.ListTipsResponse_ListTipsData, error) {
	// Convert raw tip records into TipData response
	var pbTips []*pb.TipData
	for _, tip := range tips {
		pbTip, err := s.tipTransformer(ctx, tip)
		if err != nil {
			return nil, errors.ToRpcError(err)
		}
		pbTips = append(pbTips, pbTip)
	}

	nextCursor := ""
	if len(pbTips) == int(pageSize) {
		nextCursor = lastTipID.Hex()
	}

	return &pb.ListTipsResponse_ListTipsData{
		Tips:       pbTips,
		NextCursor: nextCursor,
	}, nil
}
