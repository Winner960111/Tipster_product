package service

import (
	"context"

	pb "src/protos/Tipster"
)

func (s *SocialServiceService) CommentOnTip(ctx context.Context, req *pb.CommentOnTipRequest) (*pb.CommentOnTipResponse, error) {
	return &pb.CommentOnTipResponse{}, nil
}
func (s *SocialServiceService) UpdateComment(ctx context.Context, req *pb.UpdateCommentRequest) (*pb.UpdateCommentResponse, error) {
	return &pb.UpdateCommentResponse{}, nil
}
func (s *SocialServiceService) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	return &pb.DeleteCommentResponse{}, nil
}
func (s *SocialServiceService) ListTipComments(ctx context.Context, req *pb.ListTipCommentsRequest) (*pb.ListTipCommentsResponse, error) {
	return &pb.ListTipCommentsResponse{}, nil
}
func (s *SocialServiceService) LikeComment(ctx context.Context, req *pb.LikeCommentRequest) (*pb.LikeCommentResponse, error) {
	return &pb.LikeCommentResponse{}, nil
}
func (s *SocialServiceService) UnlikeComment(ctx context.Context, req *pb.UnlikeCommentRequest) (*pb.UnlikeCommentResponse, error) {
	return &pb.UnlikeCommentResponse{}, nil
}
func (s *SocialServiceService) ReplyComment(ctx context.Context, req *pb.ReplyCommentRequest) (*pb.ReplyCommentResponse, error) {
	return &pb.ReplyCommentResponse{}, nil
}
func (s *SocialServiceService) ListCommentReplies(ctx context.Context, req *pb.ListCommentRepliesRequest) (*pb.ListCommentRepliesResponse, error) {
	return &pb.ListCommentRepliesResponse{}, nil
}
func (s *SocialServiceService) ListComments(ctx context.Context, req *pb.ListCommentsRequest) (*pb.ListCommentsResponse, error) {
	return &pb.ListCommentsResponse{}, nil
}
func (s *SocialServiceService) ListFollowingFeed(ctx context.Context, req *pb.ListFollowingFeedRequest) (*pb.ListFollowingFeedResponse, error) {
	return &pb.ListFollowingFeedResponse{}, nil
}
