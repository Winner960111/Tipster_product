package service

import (
	"context"
	"fmt"
	"time"

	"src/internal/model"
	pb "src/protos/Tipster"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialServiceService) CommentOnTip(ctx context.Context, req *pb.CommentOnTipRequest) (*pb.CommentOnTipResponse, error) {
	// Validate tip ID
	_, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.CommentOnTipResponse{
			Code: "COMM0101",
			Msg:  "Invalid tip ID format",
		}, nil
	}

	currentTime := time.Now().UTC()

	// Generate new ObjectID for the comment
	newCommentID := primitive.NewObjectID()

	// If parentId is empty, use the new comment's ID as parentId
	parentID := req.ParentId
	if parentID == "" {
		parentID = newCommentID.Hex()
	}

	comment := &model.Comment{
		ID:        newCommentID,
		TipID:     req.TipId,
		UserID:    req.UserId,
		ParentID:  parentID,
		Content:   req.Content,
		Likes:     []primitive.ObjectID{},
		Unlikes:   []primitive.ObjectID{},
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	commentID, err := s.repo.CreateComment(ctx, comment)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to create comment", "error", err)
		return &pb.CommentOnTipResponse{
			Code: "COMM0501",
			Msg:  "Failed to create comment",
		}, nil
	}

	return &pb.CommentOnTipResponse{
		Code: "COMM0000",
		Msg:  "Comment created successfully",
		Data: &pb.CommentInfo{
			CommentId: commentID.Hex(),
			TipId:     req.TipId,
			UserId:    req.UserId,
			ParentId:  parentID,
			Content:   req.Content,
			CreatedAt: timestamppb.New(currentTime),
			UpdatedAt: timestamppb.New(currentTime),
			Likes:     []*pb.UserDetail{},
			Unlikes:   []*pb.UserDetail{},
		},
	}, nil
}
func (s *SocialServiceService) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	// Validate comment ID
	commentID, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return &pb.UpdateCommentResponse{
			Code: "COMM0101",
			Msg:  "Invalid comment ID format",
		}, nil
	}

	currentTime := time.Now().UTC()

	// Attempt to update the comment
	err = s.repo.UpdateComment(ctx, commentID, req.Content, currentTime)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UpdateCommentResponse{
				Code: "COMM0300",
				Msg:  "Comment not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to update comment", "error", err)
		return &pb.UpdateCommentResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	return &pb.UpdateCommentResponse{
		Code: "COMM0000",
		Msg:  "Comment updated successfully",
	}, nil
}
func (s *SocialServiceService) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	// Validate comment ID
	commentID, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return &pb.DeleteCommentResponse{
			Code: "COMM0101",
			Msg:  "Invalid comment ID format",
		}, nil
	}

	// Attempt to delete the comment
	err = s.repo.DeleteComment(ctx, commentID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.DeleteCommentResponse{
				Code: "COMM0103",
				Msg:  "Comment not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to delete comment", "error", err)
		return &pb.DeleteCommentResponse{
			Code: "COMM0102",
			Msg:  "Failed to delete comment",
		}, nil
	}

	return &pb.DeleteCommentResponse{
		Code: "COMM0000",
		Msg:  "Comment deleted successfully",
	}, nil
}
func (s *SocialServiceService) ListTipComments(ctx context.Context, req *pb.ListTipCommentsRequest) (*pb.ListTipCommentsResponse, error) {
	// Validate tip ID
	_, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.ListTipCommentsResponse{
			Code: "COMM0101",
			Msg:  "Invalid tip ID format",
		}, nil
	}

	// Fetch comments
	comments, err := s.repo.ListTip(ctx, req.TipId)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch comments", "error", err)
		return &pb.ListTipCommentsResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	// Collect all like and unlike user IDs for batch query
	var allLikeIDs, allUnlikeIDs []primitive.ObjectID
	for _, comment := range comments {
		allLikeIDs = append(allLikeIDs, comment.Likes...)
		allUnlikeIDs = append(allUnlikeIDs, comment.Unlikes...)
	}

	// Fetch user details in bulk
	likeUserDetails, err := s.repo.GetUserDetails(ctx, allLikeIDs)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch like user details", "error", err)
		return nil, status.Errorf(codes.Internal, "Error fetching like user details: %v", err)
	}

	unlikeUserDetails, err := s.repo.GetUserDetails(ctx, allUnlikeIDs)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch unlike user details", "error", err)
		return nil, status.Errorf(codes.Internal, "Error fetching unlike user details: %v", err)
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

	return &pb.ListTipCommentsResponse{
		Code:     "COMM0000",
		Msg:      "Comments retrieved successfully",
		Comments: pbComments,
	}, nil
}
func (s *SocialServiceService) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.LikeCommentResponse, error) {
	if req.CommentId == "" || req.UserId == "" {
		return &pb.LikeCommentResponse{
			Code: "COMM0101",
			Msg:  "Comment ID and User ID are required",
		}, nil
	}

	// Convert to ObjectID
	commentID, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return &pb.LikeCommentResponse{
			Code: "COMM0101",
			Msg:  "Invalid comment ID format",
		}, nil
	}

	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.LikeCommentResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to like the comment
	totalLikes, err := s.repo.LikeComment(ctx, commentID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.LikeCommentResponse{
				Code: "COMM0300",
				Msg:  "Comment not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to like comment", "error", err)
		return &pb.LikeCommentResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	return &pb.LikeCommentResponse{
		Code: "COMM0000",
		Msg:  "Comment liked successfully",
		Data: &pb.LikeCommentResponse_LikeCommentData{
			TotalLikes: totalLikes,
			UserLiked:  true,
		},
	}, nil
}
func (s *SocialServiceService) UnlikeComment(ctx context.Context, req *pb.UnlikeCommentRequest) (*pb.UnlikeCommentResponse, error) {
	if req.CommentId == "" || req.UserId == "" {
		return &pb.UnlikeCommentResponse{
			Code: "COMM0101",
			Msg:  "Comment ID and User ID are required",
		}, nil
	}

	// Convert to ObjectID
	commentID, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return &pb.UnlikeCommentResponse{
			Code: "COMM0101",
			Msg:  "Invalid comment ID format",
		}, nil
	}

	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.UnlikeCommentResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to unlike the comment
	totalUnlikes, err := s.repo.UnlikeComment(ctx, commentID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UnlikeCommentResponse{
				Code: "COMM0300",
				Msg:  "Comment not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to unlike comment", "error", err)
		return &pb.UnlikeCommentResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	return &pb.UnlikeCommentResponse{
		Code: "COMM0000",
		Msg:  "Comment unliked successfully",
		Data: &pb.UnlikeCommentResponse_UnlikeCommentData{
			TotalUnLikes: totalUnlikes,
			UserUnLiked:  true,
		},
	}, nil
}
func (s *SocialServiceService) ReplyComment(ctx context.Context, req *pb.ReplyCommentRequest) (*pb.ReplyCommentResponse, error) {

	// Validate parent comment ID
	_, err := primitive.ObjectIDFromHex(req.ParentCommentId)
	if err != nil {
		return &pb.ReplyCommentResponse{
			Code: "COMM0101",
			Msg:  "Invalid parent comment ID format",
		}, nil
	}
	_, err = primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.ReplyCommentResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	currentTime := time.Now().UTC()

	// Create new reply document
	reply := &model.Comment{
		UserID:    req.UserId,
		ParentID:  req.ParentCommentId,
		Content:   req.Content,
		Likes:     []primitive.ObjectID{},
		Unlikes:   []primitive.ObjectID{},
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	replyID, err := s.repo.CreateReply(ctx, reply)

	fmt.Println("replyID", replyID)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to create reply", "error", err)
		return &pb.ReplyCommentResponse{
			Code: "COMM0501",
			Msg:  "Failed to create reply",
		}, nil
	}

	return &pb.ReplyCommentResponse{
		Code: "COMM0000",
		Msg:  "Reply created successfully",
		Data: &pb.ReplyCommentResponse_ReplyData{
			ReplyId:         replyID.Hex(),
			ParentCommentId: req.ParentCommentId,
			UserId:          req.UserId,
			Content:         req.Content,
			DateCreated:     timestamppb.New(currentTime),
		},
	}, nil
}
func (s *SocialServiceService) ListCommentReplies(ctx context.Context, req *pb.ListCommentRepliesRequest) (*pb.ListCommentRepliesResponse, error) {
	// Validate parent comment ID
	_, err := primitive.ObjectIDFromHex(req.ParentCommentId)
	if err != nil {
		return &pb.ListCommentRepliesResponse{
			Code: "COMM0101",
			Msg:  "Invalid parent comment ID format",
		}, nil
	}

	// Fetch replies
	replies, err := s.repo.ListReplies(ctx, req.ParentCommentId)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch replies", "error", err)
		return &pb.ListCommentRepliesResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	// Collect all like and unlike user IDs for batch query
	var allLikeIDs, allUnlikeIDs []primitive.ObjectID
	for _, reply := range replies {
		allLikeIDs = append(allLikeIDs, reply.Likes...)
		allUnlikeIDs = append(allUnlikeIDs, reply.Unlikes...)
	}

	// Fetch user details in bulk
	likeUserDetails, err := s.repo.GetUserDetails(ctx, allLikeIDs)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch like user details", "error", err)
		return nil, status.Errorf(codes.Internal, "Error fetching like user details: %v", err)
	}

	unlikeUserDetails, err := s.repo.GetUserDetails(ctx, allUnlikeIDs)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch unlike user details", "error", err)
		return nil, status.Errorf(codes.Internal, "Error fetching unlike user details: %v", err)
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

	return &pb.ListCommentRepliesResponse{
		Code:    "COMM0000",
		Msg:     "Replies retrieved successfully",
		Replies: pbReplies,
	}, nil
}
func (s *SocialServiceService) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	pageSize := int64(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	comments, nextCursor, err := s.repo.ListComments(ctx, pageSize, req.NextCursor)
	if err != nil {
		return &pb.ListCommentsResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, err
	}

	var allLikeIDs, allUnlikeIDs []primitive.ObjectID
	for _, comment := range comments {
		allLikeIDs = append(allLikeIDs, comment.Likes...)
		allUnlikeIDs = append(allUnlikeIDs, comment.Unlikes...)
	}

	likeUserDetails, err := s.repo.GetUserDetails(ctx, allLikeIDs)
	if err != nil {
		return &pb.ListCommentsResponse{
			Code: "COMM0502",
			Msg:  "Error fetching like user details",
		}, err
	}

	unlikeUserDetails, err := s.repo.GetUserDetails(ctx, allUnlikeIDs)
	if err != nil {
		return &pb.ListCommentsResponse{
			Code: "COMM0503",
			Msg:  "Error fetching unlike user details",
		}, err
	}

	likeUserMap := make(map[string]*pb.UserDetail)
	for _, user := range likeUserDetails {
		likeUserMap[user.ID.Hex()] = &pb.UserDetail{
			Id:       user.ID.Hex(),
			UserName: user.Username,
		}
	}

	unlikeUserMap := make(map[string]*pb.UserDetail)
	for _, user := range unlikeUserDetails {
		unlikeUserMap[user.ID.Hex()] = &pb.UserDetail{
			Id:       user.ID.Hex(),
			UserName: user.Username,
		}
	}

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
	fmt.Println("nextCursor", nextCursor)
	return &pb.ListCommentsResponse{
		Code:       "COMM0000",
		Msg:        "Comments retrieved successfully",
		Comments:   pbComments,
		NextCursor: nextCursor,
	}, nil
}

// func (s *SocialServiceService) ListFollowingFeed(ctx context.Context, req *pb.ListFollowingFeedRequest) (*pb.ListFollowingFeedResponse, error) {
// 	return &pb.ListFollowingFeedResponse{}, nil
// }
