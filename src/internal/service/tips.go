package service

import (
	"context"
	"strings"
	"time"

	"src/internal/model"
	pb "src/protos/Tipster"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *SocialServiceService) CreateTip(ctx context.Context, req *pb.CreateTipRequest) (*pb.CreateTipResponse, error) {
	currentTime := time.Now().UTC()

	tip := &model.Tip{
		TipsterID: req.TipsterId,
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		ShareType: req.ShareType,
		Likes:     []primitive.ObjectID{},
		Unlikes:   []primitive.ObjectID{},
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	tipID, err := s.repo.CreateTip(ctx, tip)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to create tip", "error", err)
		return &pb.CreateTipResponse{
			Code: "COMM0501",
			Msg:  "Failed to create tip",
		}, nil
	}

	return &pb.CreateTipResponse{
		Code: "COMM0000",
		Msg:  "Tip created successfully",
		Data: &pb.TipData{
			TipId:     tipID.Hex(),
			TipsterId: req.TipsterId,
			Title:     req.Title,
			Content:   req.Content,
			Tags:      req.Tags,
			CreatedAt: timestamppb.New(currentTime),
			UpdatedAt: timestamppb.New(currentTime),
			ShareType: req.ShareType,
		},
	}, nil
}
func (s *SocialServiceService) GetTip(ctx context.Context, req *pb.GetTipRequest) (*pb.GetTipResponse, error) {
	if req.TipId == "" {
		return &pb.GetTipResponse{
			Code: "COMM0101",
			Msg:  "Tip ID is required",
		}, nil
	}

	objID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.GetTipResponse{
			Code: "COMM0101",
			Msg:  "Invalid tip ID format",
		}, nil
	}

	tip, err := s.repo.GetTip(ctx, objID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, status.Errorf(codes.NotFound, "Tip not found")
		}
		s.logger.Log(log.LevelError, "failed to fetch tip", "error", err)
		return nil, status.Errorf(codes.Internal, "Error fetching tip: %v", err)
	}

	// Get likes details
	likesDetails, err := s.repo.GetUserDetails(ctx, tip.Likes)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch likes details", "error", err)
		return nil, status.Errorf(codes.Internal, "Error fetching likes details: %v", err)
	}

	// Get unlikes details
	unLikesDetails, err := s.repo.GetUserDetails(ctx, tip.Unlikes)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch unlikes details", "error", err)
		return nil, status.Errorf(codes.Internal, "Error fetching unlikes details: %v", err)
	}

	pbLikes := make([]*pb.UserDetail, len(likesDetails))
	for i, l := range likesDetails {
		pbLikes[i] = &pb.UserDetail{
			Id:       strings.TrimSpace(l.ID.Hex()),
			UserName: strings.TrimSpace(l.Username),
		}
	}

	pbUnlikes := make([]*pb.UserDetail, len(unLikesDetails))
	for i, ul := range unLikesDetails {
		pbUnlikes[i] = &pb.UserDetail{
			Id:       strings.TrimSpace(ul.ID.Hex()),
			UserName: strings.TrimSpace(ul.Username),
		}
	}

	return &pb.GetTipResponse{
		Code: "COMM0000",
		Msg:  "Tip fetched successfully",
		Data: &pb.TipData{
			TipId:     strings.TrimSpace(tip.ID.Hex()),
			TipsterId: tip.TipsterID,
			Title:     tip.Title,
			Content:   tip.Content,
			Tags:      tip.Tags,
			Likes:     pbLikes,
			Unlikes:   pbUnlikes,
			CreatedAt: timestamppb.New(tip.CreatedAt),
			UpdatedAt: timestamppb.New(tip.UpdatedAt),
			ShareType: tip.ShareType,
		},
	}, nil
}
func (s *SocialServiceService) UpdateTip(ctx context.Context, req *pb.UpdateTipRequest) (*pb.UpdateTipResponse, error) {
	// Validate tip ID
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.UpdateTipResponse{
			Code: "COMM0101",
			Msg:  "Invalid tip ID format",
		}, nil
	}

	currentTime := time.Now().UTC()

	updates := bson.M{
		"$set": bson.M{
			"title":     req.Title,
			"content":   req.Content,
			"tags":      req.Tags,
			"shareType": req.ShareType,
			"updatedAt": currentTime,
		},
	}

	_, err = s.repo.UpdateTip(ctx, tipID, updates)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UpdateTipResponse{
				Code: "COMM0300",
				Msg:  "Tip not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to update tip", "error", err)
		return &pb.UpdateTipResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	return &pb.UpdateTipResponse{
		Code: "COMM0000",
		Msg:  "Tip updated successfully",
	}, nil
}
func (s *SocialServiceService) DeleteTip(ctx context.Context, req *pb.DeleteTipRequest) (*pb.DeleteTipResponse, error) {
	// Validate tip ID
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.DeleteTipResponse{
			Code: "COMM0101",
			Msg:  "Invalid tip ID format",
		}, nil
	}

	// Attempt to delete the tip
	err = s.repo.DeleteTip(ctx, tipID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.DeleteTipResponse{
				Code: "COMM0300",
				Msg:  "Tip not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to delete tip", "error", err)
		return &pb.DeleteTipResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	return &pb.DeleteTipResponse{
		Code: "COMM0000",
		Msg:  "Tip deleted successfully",
	}, nil
}
func (s *SocialServiceService) ListTips(ctx context.Context, req *pb.ListTipsRequest) (*pb.ListTipsResponse, error) {
	// Default page size
	pageSize := int64(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10 // Default to 10 items per page
	}

	// Build filter
	filter := bson.M{}
	if req.TipsterId != "" {
		filter["tipsterId"] = req.TipsterId
	}

	// Fetch tips
	tips, lastTipID, err := s.repo.ListTips(ctx, filter, pageSize, req.NextCursor)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to fetch tips", "error", err)
		return &pb.ListTipsResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	// Collect all like and unlike user IDs for batch query
	var allLikeIDs, allUnlikeIDs []primitive.ObjectID
	for _, tip := range tips {
		allLikeIDs = append(allLikeIDs, tip.Likes...)
		allUnlikeIDs = append(allUnlikeIDs, tip.Unlikes...)
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

	// Convert raw tip records into TipData response
	var pbTips []*pb.TipData
	for _, tip := range tips {
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

		pbTips = append(pbTips, &pb.TipData{
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
		})
	}

	// Determine `NextCursor`
	nextCursor := ""
	if len(pbTips) == int(pageSize) {
		nextCursor = lastTipID.Hex()
	}

	return &pb.ListTipsResponse{
		Code: "COMM0000",
		Msg:  "Tips retrieved successfully",
		Data: &pb.ListTipsResponse_ListTipsData{
			Tips:       pbTips,
			NextCursor: nextCursor,
		},
	}, nil
}
func (s *SocialServiceService) ShareTip(ctx context.Context, req *pb.ShareTipRequest) (*pb.ShareTipResponse, error) {
	// Validate tip ID
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.ShareTipResponse{
			Code: "COMM0101",
			Msg:  "Invalid tip ID format",
		}, nil
	}

	currentTime := time.Now().UTC()

	// Attempt to update the tip's share type
	err = s.repo.ShareTip(ctx, tipID, req.ShareType, currentTime)
	if err != nil {
		s.logger.Log(log.LevelError, "failed to update tip share type", "error", err)
		return &pb.ShareTipResponse{
			Code: "COMM0501",
			Msg:  "Failed to update tip share type",
		}, nil
	}

	return &pb.ShareTipResponse{
		Code: "COMM0000",
		Msg:  "Tip shared successfully",
	}, nil
}
func (s *SocialServiceService) LikeTip(ctx context.Context, req *pb.LikeTipRequest) (*pb.LikeTipResponse, error) {
	if req.TipId == "" || req.UserId == "" {
		return &pb.LikeTipResponse{
			Code: "COMM0101",
			Msg:  "Tip ID and User ID are required",
		}, nil
	}

	// Convert to ObjectID
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.LikeTipResponse{
			Code: "COMM0101",
			Msg:  "Invalid tip ID format",
		}, nil
	}

	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.LikeTipResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to like the tip
	totalLikes, err := s.repo.LikeTip(ctx, tipID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.LikeTipResponse{
				Code: "COMM0300",
				Msg:  "Tip not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to like tip", "error", err)
		return &pb.LikeTipResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	return &pb.LikeTipResponse{
		Code: "COMM0000",
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
			Code: "COMM0101",
			Msg:  "Tip ID and User ID are required",
		}, nil
	}

	// Convert to ObjectID
	tipID, err := primitive.ObjectIDFromHex(req.TipId)
	if err != nil {
		return &pb.UnlikeTipResponse{
			Code: "COMM0101",
			Msg:  "Invalid tip ID format",
		}, nil
	}

	userID, err := primitive.ObjectIDFromHex(req.UserId)
	if err != nil {
		return &pb.UnlikeTipResponse{
			Code: "COMM0101",
			Msg:  "Invalid user ID format",
		}, nil
	}

	// Attempt to unlike the tip
	totalUnlikes, err := s.repo.UnlikeTip(ctx, tipID, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &pb.UnlikeTipResponse{
				Code: "COMM0300",
				Msg:  "Tip not found",
			}, nil
		}
		s.logger.Log(log.LevelError, "failed to unlike tip", "error", err)
		return &pb.UnlikeTipResponse{
			Code: "COMM0501",
			Msg:  "Database error",
		}, nil
	}

	return &pb.UnlikeTipResponse{
		Code: "COMM0000",
		Msg:  "Tip unliked successfully",
		Data: &pb.UnlikeTipResponse_UnLikeTipData{
			TotalUnLikes: totalUnlikes,
			UserUnLiked:  true,
		},
	}, nil
}
