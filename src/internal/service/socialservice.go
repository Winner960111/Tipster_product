package service

import (
	"context"

	pb "src/protos/Tipster"
)

type SocialServiceService struct {
	pb.UnimplementedSocialServiceServer
	// repo repository.SocialServiceRepository
}

func NewSocialServiceService() *SocialServiceService {
	return &SocialServiceService{}
}

func (s *SocialServiceService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	return &pb.CreateUserResponse{}, nil
}
func (s *SocialServiceService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return &pb.GetUserResponse{}, nil
}
func (s *SocialServiceService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	return &pb.UpdateUserResponse{}, nil
}
func (s *SocialServiceService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return &pb.DeleteUserResponse{}, nil
}
func (s *SocialServiceService) ListUsers(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserResponse, error) {
	return &pb.ListUserResponse{}, nil
}
func (s *SocialServiceService) FollowTipster(ctx context.Context, req *pb.FollowTipsterRequest) (*pb.FollowTipsterResponse, error) {
	return &pb.FollowTipsterResponse{}, nil
}
func (s *SocialServiceService) UnfollowTipster(ctx context.Context, req *pb.UnFollowTipsterRequest) (*pb.UnfollowTipsterResponse, error) {
	return &pb.UnfollowTipsterResponse{}, nil
}
func (s *SocialServiceService) CreateTip(ctx context.Context, req *pb.CreateTipRequest) (*pb.CreateTipResponse, error) {
	return &pb.CreateTipResponse{}, nil
}
func (s *SocialServiceService) GetTip(ctx context.Context, req *pb.GetTipRequest) (*pb.GetTipResponse, error) {
	return &pb.GetTipResponse{}, nil
}
func (s *SocialServiceService) UpdateTip(ctx context.Context, req *pb.UpdateTipRequest) (*pb.UpdateTipResponse, error) {
	return &pb.UpdateTipResponse{}, nil
}
func (s *SocialServiceService) DeleteTip(ctx context.Context, req *pb.DeleteTipRequest) (*pb.DeleteTipResponse, error) {
	return &pb.DeleteTipResponse{}, nil
}
func (s *SocialServiceService) ListTips(ctx context.Context, req *pb.ListTipsRequest) (*pb.ListTipsResponse, error) {
	return &pb.ListTipsResponse{}, nil
}
func (s *SocialServiceService) ShareTip(ctx context.Context, req *pb.ShareTipRequest) (*pb.ShareTipResponse, error) {
	return &pb.ShareTipResponse{}, nil
}
func (s *SocialServiceService) LikeTip(ctx context.Context, req *pb.LikeTipRequest) (*pb.LikeTipResponse, error) {
	return &pb.LikeTipResponse{}, nil
}
func (s *SocialServiceService) UnlikeTip(ctx context.Context, req *pb.UnlikeTipRequest) (*pb.UnlikeTipResponse, error) {
	return &pb.UnlikeTipResponse{}, nil
}
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
