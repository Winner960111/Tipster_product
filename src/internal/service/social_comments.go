package service

import (
	"context"

	pb "src/protos/Tipster"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *SocialServiceService) CommentOnTip(ctx context.Context, req *pb.CommentOnTipRequest) (*pb.CommentOnTipResponse, error) {
	_, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.CommentOnTipResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tip ID format",
		}, nil
	}

	data, err := s.biz.CreateComment(ctx, req)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to create comment", "error", err)
		return &pb.CommentOnTipResponse{
			Code: CodeError,
			Msg:  "Failed to create comment",
		}, nil
	}

	return &pb.CommentOnTipResponse{
		Code: CodeOk,
		Msg:  "Comment created successfully",
		Data: data,
	}, nil
}
func (s *SocialServiceService) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	// Validate comment ID
	commentID, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return &pb.UpdateCommentResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid comment ID format",
		}, nil
	}
	// Attempt to update the comment
	err = s.biz.UpdateComment(ctx, commentID, req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UpdateCommentResponse{
				Code: CodeNotFound,
				Msg:  "Comment not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to update comment", "error", err)
		return &pb.UpdateCommentResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}

	return &pb.UpdateCommentResponse{
		Code: CodeOk,
		Msg:  "Comment updated successfully",
	}, nil
}
func (s *SocialServiceService) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	// Validate comment ID
	commentID, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return &pb.DeleteCommentResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid comment ID format",
		}, nil
	}

	// Attempt to delete the comment
	err = s.biz.DeleteComment(ctx, commentID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.DeleteCommentResponse{
				Code: CodeNotFound,
				Msg:  "Comment not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to delete comment", "error", err)
		return &pb.DeleteCommentResponse{
			Code: CodeInvalid,
			Msg:  "Failed to delete comment",
		}, nil
	}

	return &pb.DeleteCommentResponse{
		Code: CodeOk,
		Msg:  "Comment deleted successfully",
	}, nil
}
func (s *SocialServiceService) ListTipComments(ctx context.Context, req *pb.ListTipCommentsRequest) (*pb.ListTipCommentsResponse, error) {
	_, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.ListTipCommentsResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tip ID format",
		}, nil
	}
	// Fetch comments
	comments, err := s.biz.ListTipComments(ctx, req.TipId)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch comments", "error", err)
		return &pb.ListTipCommentsResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}
	pbComments, err := s.likeTransformer(ctx, comments)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to transform comments", "error", err)
		return &pb.ListTipCommentsResponse{
			Code: CodeError,
			Msg:  "Failed to transform comments",
		}, nil
	}
	return &pb.ListTipCommentsResponse{
		Code:     CodeOk,
		Msg:      "Comments retrieved successfully",
		Comments: pbComments,
	}, nil
}
func (s *SocialServiceService) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.LikeCommentResponse, error) {
	if req.CommentId == "" || req.UserId == "" {
		return &pb.LikeCommentResponse{
			Code: CodeInvalid,
			Msg:  "Comment ID and User ID are required",
		}, nil
	}
	commentID, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return &pb.LikeCommentResponse{
			Code: CodeInvalid,
			Msg:  "Invalid comment ID format",
		}, nil
	}

	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.LikeCommentResponse{
			Code: CodeInvalid,
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to like the comment
	totalLikes, err := s.biz.LikeComment(ctx, commentID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.LikeCommentResponse{
				Code: CodeNotFound,
				Msg:  "Comment not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to like comment", "error", err)
		return &pb.LikeCommentResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}

	return &pb.LikeCommentResponse{
		Code: CodeOk,
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
			Code: CodeInvalid,
			Msg:  "Comment ID and User ID are required",
		}, nil
	}

	commentID, err := primitive.ObjectIDFromHex(req.CommentId)
	if err != nil {
		return &pb.UnlikeCommentResponse{
			Code: CodeInvalid,
			Msg:  "Invalid comment ID format",
		}, nil
	}

	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.UnlikeCommentResponse{
			Code: CodeInvalid,
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to unlike the comment
	totalUnlikes, err := s.biz.UnlikeComment(ctx, commentID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UnlikeCommentResponse{
				Code: CodeNotFound,
				Msg:  "Comment not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to unlike comment", "error", err)
		return &pb.UnlikeCommentResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}

	return &pb.UnlikeCommentResponse{
		Code: CodeOk,
		Msg:  "Comment unliked successfully",
		Data: &pb.UnlikeCommentResponse_UnlikeCommentData{
			TotalUnLikes: totalUnlikes,
			UserUnLiked:  true,
		},
	}, nil
}
func (s *SocialServiceService) ReplyComment(ctx context.Context, req *pb.ReplyCommentRequest) (*pb.ReplyCommentResponse, error) {

	_, err := primitive.ObjectIDFromHex(req.ParentCommentId)
	if err != nil {
		return &pb.ReplyCommentResponse{
			Code: CodeInvalid,
			Msg:  "Invalid parent comment ID format",
		}, nil
	}
	_, err = primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.ReplyCommentResponse{
			Code: CodeInvalid,
			Msg:  "Invalid user ID format",
		}, nil
	}

	data, err := s.biz.CreateReply(ctx, req)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to create reply", "error", err)
		return &pb.ReplyCommentResponse{
			Code: CodeError,
			Msg:  "Failed to create reply",
		}, nil
	}
	return &pb.ReplyCommentResponse{
		Code: CodeOk,
		Msg:  "Reply created successfully",
		Data: data,
	}, nil
}
func (s *SocialServiceService) ListCommentReplies(ctx context.Context, req *pb.ListCommentRepliesRequest) (*pb.ListCommentRepliesResponse, error) {
	_, err := primitive.ObjectIDFromHex(req.ParentCommentId)
	if err != nil {
		return &pb.ListCommentRepliesResponse{
			Code: CodeInvalid,
			Msg:  "Invalid parent comment ID format",
		}, nil
	}

	replies, err := s.biz.ListReplies(ctx, req.ParentCommentId)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch replies", "error", err)
		return &pb.ListCommentRepliesResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}

	data, err := s.replyTransformer(ctx, replies)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to transform replies", "error", err)
		return &pb.ListCommentRepliesResponse{
			Code: CodeError,
			Msg:  "Failed to transform replies",
		}, nil
	}

	return &pb.ListCommentRepliesResponse{
		Code:    CodeOk,
		Msg:     "Replies retrieved successfully",
		Replies: data,
	}, nil
}
func (s *SocialServiceService) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	pageSize := int64(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	comments, nextCursor, err := s.biz.ListComments(ctx, pageSize, req.NextCursor)
	if err != nil {
		return &pb.ListCommentsResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, err
	}

	pbComments, err := s.likeTransformer(ctx, comments)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to transform comments", "error", err)
		return &pb.ListCommentsResponse{
			Code: CodeError,
			Msg:  "Failed to transform comments",
		}, nil
	}

	return &pb.ListCommentsResponse{
		Code:       CodeOk,
		Msg:        "Comments retrieved successfully",
		Comments:   pbComments,
		NextCursor: nextCursor,
	}, nil
}

// func (s *SocialServiceService) ListFollowingFeed(ctx context.Context, req *pb.ListFollowingFeedRequest) (*pb.ListFollowingFeedResponse, error) {
// 	return &pb.ListFollowingFeedResponse{}, nil
// }
