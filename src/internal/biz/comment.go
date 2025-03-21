package biz

import (
	"context"
	"src/internal/errors"
	"src/internal/model"
	pb "src/protos/Tipster"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialService) CreateComment(ctx context.Context, req *pb.CommentOnTipRequest) (*pb.CommentInfo, error) {
	currentTime := time.Now().UTC()
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

	commentID, err := s.Repo.CreateComment(ctx, comment)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}

	return &pb.CommentInfo{
		CommentId: commentID.Hex(),
		TipId:     req.TipId,
		UserId:    req.UserId,
		ParentId:  parentID,
		Content:   req.Content,
		CreatedAt: timestamppb.New(currentTime),
		UpdatedAt: timestamppb.New(currentTime),
		Likes:     []*pb.UserDetail{},
		Unlikes:   []*pb.UserDetail{},
	}, nil
}

func (s *SocialService) UpdateComment(ctx context.Context, commentID primitive.ObjectID, req *pb.UpdateCommentRequest) error {
	currentTime := time.Now().UTC()
	return s.Repo.UpdateComment(ctx, commentID, req.Content, currentTime)
}

func (s *SocialService) DeleteComment(ctx context.Context, commentID primitive.ObjectID) error {
	return s.Repo.DeleteComment(ctx, commentID)
}

func (s *SocialService) ListTipComments(ctx context.Context, tipID string) ([]*model.Comment, error) {
	comments, err := s.Repo.ListTipComments(ctx, tipID)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}

	return comments, nil
}

func (s *SocialService) LikeComment(ctx context.Context, commentID primitive.ObjectID, userID primitive.ObjectID) (int32, error) {
	return s.Repo.LikeComment(ctx, commentID, userID)
}

func (s *SocialService) UnlikeComment(ctx context.Context, commentID primitive.ObjectID, userID primitive.ObjectID) (int32, error) {
	return s.Repo.UnlikeComment(ctx, commentID, userID)
}

func (s *SocialService) CreateReply(ctx context.Context, req *pb.ReplyCommentRequest) (*pb.ReplyCommentResponse_ReplyData, error) {
	// Create new reply document
	currentTime := time.Now().UTC()
	reply := &model.Comment{
		UserID:    req.UserId,
		ParentID:  req.ParentCommentId,
		Content:   req.Content,
		Likes:     []primitive.ObjectID{},
		Unlikes:   []primitive.ObjectID{},
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	replyID, err := s.Repo.CreateReply(ctx, reply)

	if err != nil {
		return nil, errors.ToRpcError(err)
	}
	return &pb.ReplyCommentResponse_ReplyData{
		ReplyId:         replyID.Hex(),
		ParentCommentId: req.ParentCommentId,
		UserId:          req.UserId,
		Content:         req.Content,
		DateCreated:     timestamppb.New(currentTime),
	}, nil
}

func (s *SocialService) ListReplies(ctx context.Context, parentCommentID string) ([]*model.Comment, error) {
	return s.Repo.ListReplies(ctx, parentCommentID)
}

func (s *SocialService) ListComments(ctx context.Context, pageSize int64, nextCursor string) ([]*model.Comment, string, error) {
	return s.Repo.ListComments(ctx, pageSize, nextCursor)
}
