package service

import (
	"context"

	"src/internal/errors"
	pb "src/protos/Tipster"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *SocialServiceService) CreateTip(ctx context.Context, req *pb.CreateTipRequest) (*pb.CreateTipResponse, error) {
	data, err := s.biz.CreateTip(ctx, req)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to create tip", "error", err)
		return &pb.CreateTipResponse{
			Code: CodeError,
			Msg:  "Failed to create tip",
		}, nil
	}
	return &pb.CreateTipResponse{
		Code: CodeOk,
		Msg:  "Tip created successfully",
		Data: data,
	}, nil
}

func (s *SocialServiceService) GetTip(ctx context.Context, req *pb.GetTipRequest) (*pb.GetTipResponse, error) {
	if req.TipId == "" {
		return &pb.GetTipResponse{
			Code: CodeError,
			Msg:  "Tip ID is required",
		}, nil
	}

	objID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.GetTipResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tip ID format",
		}, nil
	}

	tip, err := s.biz.GetTip(ctx, objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.ToRpcError(err)
		}
		s.logger.Log(log.LevelError, "failed to fetch tip", "error", err)
		return nil, errors.ToRpcError(err)
	}

	data, err := s.tipTransformer(ctx, tip)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}

	return &pb.GetTipResponse{
		Code: CodeOk,
		Msg:  "Tip fetched successfully",
		Data: data,
	}, nil
}

func (s *SocialServiceService) UpdateTip(ctx context.Context, req *pb.UpdateTipRequest) (*pb.UpdateTipResponse, error) {
	// Validate tip ID
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.UpdateTipResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tip ID format",
		}, nil
	}
	_, err = s.biz.UpdateTip(ctx, tipID, req)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UpdateTipResponse{
				Code: CodeNotFound,
				Msg:  "Tip not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to update tip", "error", err)
		return &pb.UpdateTipResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}

	return &pb.UpdateTipResponse{
		Code: CodeOk,
		Msg:  "Tip updated successfully",
	}, nil
}

func (s *SocialServiceService) DeleteTip(ctx context.Context, req *pb.DeleteTipRequest) (*pb.DeleteTipResponse, error) {
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.DeleteTipResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tip ID format",
		}, nil
	}

	err = s.biz.DeleteTip(ctx, tipID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.DeleteTipResponse{
				Code: CodeNotFound,
				Msg:  "Tip not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to delete tip", "error", err)
		return &pb.DeleteTipResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}

	return &pb.DeleteTipResponse{
		Code: CodeOk,
		Msg:  "Tip deleted successfully",
	}, nil
}

func (s *SocialServiceService) ListTips(ctx context.Context, req *pb.ListTipsRequest) (*pb.ListTipsResponse, error) {
	// Default page size
	pageSize := int64(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10 // Default to 10 items per page
	}
	// Fetch tips
	tips, lastTipID, err := s.biz.ListTips(ctx, req)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch tips", "error", err)
		return &pb.ListTipsResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}
	data, err := s.tipsTransformer(ctx, tips, pageSize, lastTipID)
	if err != nil {
		return nil, errors.ToRpcError(err)
	}

	return &pb.ListTipsResponse{
		Code: CodeOk,
		Msg:  "Tips retrieved successfully",
		Data: data,
	}, nil
}

func (s *SocialServiceService) ShareTip(ctx context.Context, req *pb.ShareTipRequest) (*pb.ShareTipResponse, error) {
	// Validate tip ID
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.ShareTipResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tip ID format",
		}, nil
	}
	err = s.biz.ShareTip(ctx, tipID, req.ShareType)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to update tip share type", "error", err)
		return &pb.ShareTipResponse{
			Code: CodeError,
			Msg:  "Failed to update tip share type",
		}, nil
	}

	return &pb.ShareTipResponse{
		Code: CodeOk,
		Msg:  "Tip shared successfully",
	}, nil
}

func (s *SocialServiceService) LikeTip(ctx context.Context, req *pb.LikeTipRequest) (*pb.LikeTipResponse, error) {
	if req.TipId == "" || req.UserId == "" {
		return &pb.LikeTipResponse{
			Code: CodeInvalidID,
			Msg:  "Tip ID and User ID are required",
		}, nil
	}
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.LikeTipResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tip ID format",
		}, nil
	}
	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.LikeTipResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid user ID format",
		}, nil
	}

	totalLikes, err := s.biz.LikeTip(ctx, tipID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.LikeTipResponse{
				Code: CodeNotFound,
				Msg:  "Tip not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to like tip", "error", err)
		return &pb.LikeTipResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}

	return &pb.LikeTipResponse{
		Code: CodeOk,
		Msg:  "Tip liked successfully",
		Data: &pb.LikeTipResponse_LikeTipData{
			TotalLikes: totalLikes,
			UserLiked:  true,
		},
	}, nil
}

func (s *SocialServiceService) UnlikeTip(ctx context.Context, req *pb.UnlikeTipRequest) (*pb.UnlikeTipResponse, error) {
	if req.TipId == "" || req.UserId == "" {
		return &pb.UnlikeTipResponse{
			Code: CodeInvalid,
			Msg:  "Tip ID and User ID are required",
		}, nil
	}
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.UnlikeTipResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid tip ID format",
		}, nil
	}

	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.UnlikeTipResponse{
			Code: CodeInvalidID,
			Msg:  "Invalid user ID format",
		}, nil
	}

	totalUnlikes, err := s.biz.UnlikeTip(ctx, tipID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UnlikeTipResponse{
				Code: CodeNotFound,
				Msg:  "Tip not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to unlike tip", "error", err)
		return &pb.UnlikeTipResponse{
			Code: CodeError,
			Msg:  "Database error",
		}, nil
	}

	return &pb.UnlikeTipResponse{
		Code: CodeOk,
		Msg:  "Tip unliked successfully",
		Data: &pb.UnlikeTipResponse_UnLikeTipData{
			TotalUnLikes: totalUnlikes,
			UserUnLiked:  true,
		},
	}, nil
}
