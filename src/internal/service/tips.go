package service

import (
	"context"

	pb "src/protos/Tipster"
)

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
