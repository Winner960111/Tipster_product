package service

import (
	"context"

	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"src/internal/model"
	"src/internal/repository"

	ymcommonpb "src/protos/YM.Common"
	ymtransactionpb "src/protos/YM.Transaction"

	"fmt"
	"strings"
)

type orderServer struct {
	ymtransactionpb.UnimplementedOrderServiceServer
	repo            repository.OrderRepository
	transactionRepo repository.TransactionRepository
	logger          log.Logger
	node            *snowflake.Node
}

func NewOrderServer(repo repository.OrderRepository, transactionRepo repository.TransactionRepository, logger log.Logger) ymtransactionpb.OrderServiceServer {
	node, err := snowflake.NewNode(1)
	if err != nil {
		return nil
	}
	return &orderServer{
		repo:            repo,
		transactionRepo: transactionRepo,
		logger:          logger,
		node:            node,
	}
}

func (o *orderServer) CreateOrder(ctx context.Context, req *ymtransactionpb.OrderInfo) (*ymtransactionpb.CreateOrderResponse, error) {
	o.logger.Log(log.LevelInfo, "msg", "Creating order", "member_id", req.MemberId)

	decimalPrice, err := primitive.ParseDecimal128(req.LocalPrice)
	if err != nil {
		o.logger.Log(log.LevelError, "msg", "Failed to parse LocalPrice", "error", err)
		return nil, status.Errorf(codes.InvalidArgument, "invalid LocalPrice format: %v", err)
	}

	now := time.Now()
	order := &model.Order{
		OrderId:          o.node.Generate().String(),
		OrderCode:        o.node.Generate().String(),
		MemberId:         req.MemberId,
		PackageCode:      req.PackageCode,
		Status:           req.Status,
		LocalPrice:       decimalPrice,
		LocalCurrency:    req.LocalCurrency,
		SubscriptionType: req.SubscriptionType,
		DateCreated:      now.Unix(),
		DateUpdated:      now.Unix(),
	}

	if len(req.BuyIn) > 0 {
		order.BuyIn = make([]model.MerchandiseModel, len(req.BuyIn))
		for i, item := range req.BuyIn {
			order.BuyIn[i] = model.MerchandiseModel{
				Type:   item.Type,
				Amount: item.Amount,
			}
		}
	}

	orderID, err := o.repo.CreateOrder(ctx, order)
	if err != nil {
		o.logger.Log(log.LevelError, "msg", "CreateOrder failed", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	o.logger.Log(log.LevelInfo, "msg", "Order created successfully", "order_id", orderID)
	return &ymtransactionpb.CreateOrderResponse{
		Code:      "COMM0000",
		OrderId:   order.OrderId,
		OrderCode: order.OrderCode,
	}, nil
}

func (o *orderServer) CreateOrderPayment(ctx context.Context, req *ymtransactionpb.CreateOrderPaymentRequest) (*ymcommonpb.ResponseBase, error) {
	o.logger.Log(log.LevelInfo, "msg", "Creating order payment", "order_id", req.OrderId)

	_, err := o.repo.FindByID(ctx, req.OrderId)
	if err != nil {
		o.logger.Log(log.LevelError, "msg", "Order not found", "order_id", req.OrderId, "error", err)
		return &ymcommonpb.ResponseBase{
			Code: "COMM0103", // order not found
		}, nil
	}

	// convert PaymentInfo from protobuf to model
	paymentInfo := &model.PaymentInfoModel{
		Platform:    req.Payment.Platform,
		PaymentId1:  req.Payment.PaymentId1,
		PaymentId2:  req.Payment.PaymentId2,
		PaymentId3:  req.Payment.PaymentId3,
		PaymentLink: req.Payment.PaymentLink,
		Status:      req.Payment.Status,
	}

	// update order payment
	err = o.repo.UpdateOrderPayment(ctx, req.OrderId, paymentInfo)
	if err != nil {
		o.logger.Log(log.LevelError, "msg", "Failed to update order payment", "order_id", req.OrderId, "error", err)
		return &ymcommonpb.ResponseBase{
			Code: "COMM9999", // general error
		}, nil
	}

	o.logger.Log(log.LevelInfo, "msg", "Order payment created successfully", "order_id", req.OrderId)
	return &ymcommonpb.ResponseBase{
		Code: "COMM0000", // success
	}, nil
}

// UpdatePayment updates specific payment fields of an order
func (o *orderServer) UpdatePayment(ctx context.Context, req *ymtransactionpb.UpdatePaymentRequest) (*ymcommonpb.ResponseBase, error) {
	o.logger.Log(log.LevelInfo, "msg", "Updating order payment", "order_id", req.OrderId)

	// Find the order to check if it exists
	order, err := o.repo.FindByID(ctx, req.OrderId)
	if err != nil {
		o.logger.Log(log.LevelError, "msg", "Order not found", "order_id", req.OrderId, "error", err)
		return &ymcommonpb.ResponseBase{
			Code: "COMM0103", // order not found
		}, nil
	}

	// Check if payment information already exists
	if order.Payment == nil {
		o.logger.Log(log.LevelWarn, "msg", "Payment information does not exist for this order", "order_id", req.OrderId)
		order.Payment = &model.PaymentInfoModel{}
	}

	// Only update non-empty fields
	if req.PaymentId1 != "" {
		order.Payment.PaymentId1 = req.PaymentId1
	}
	if req.PaymentId2 != "" {
		order.Payment.PaymentId2 = req.PaymentId2
	}
	if req.PaymentId3 != "" {
		order.Payment.PaymentId3 = req.PaymentId3
	}
	if req.PaymentLink != "" {
		order.Payment.PaymentLink = req.PaymentLink
	}

	// Update the payment information in the database
	err = o.repo.UpdateOrderPayment(ctx, req.OrderId, order.Payment)
	if err != nil {
		o.logger.Log(log.LevelError, "msg", "Failed to update order payment", "order_id", req.OrderId, "error", err)
		return &ymcommonpb.ResponseBase{
			Code: "COMM9999", // general error
		}, nil
	}

	o.logger.Log(log.LevelInfo, "msg", "Order payment updated successfully", "order_id", req.OrderId)
	return &ymcommonpb.ResponseBase{
		Code: "COMM0000", // success
	}, nil
}

// UpdatePaymentResult updates the payment status and order status
func (o *orderServer) UpdatePaymentResult(ctx context.Context, req *ymtransactionpb.UpdatePaymentResultRequest) (*ymcommonpb.ResponseBase, error) {
	o.logger.Log(log.LevelInfo, "msg", "Updating payment result", "order_id", req.OrderId, "platform_status", req.PlatformStatus, "order_status", req.OrderStatus)

	// Find the order to check if it exists
	order, err := o.repo.FindByID(ctx, req.OrderId)
	if err != nil {
		o.logger.Log(log.LevelError, "msg", "Order not found", "order_id", req.OrderId, "error", err)
		// According to requirements, we return COMM0103 even if order is not found
		return &ymcommonpb.ResponseBase{
			Code: "COMM0103",
		}, nil
	}

	// Check if payment information exists
	if order.Payment == nil && req.PlatformStatus != "none" {
		o.logger.Log(log.LevelWarn, "msg", "Payment information does not exist for this order but trying to update status", "order_id", req.OrderId)
	}

	// Determine if we need to update order status
	// We'll use -1 as a marker value to indicate "don't update"
	orderStatus := int32(-1)
	if req.OrderStatus != 0 && req.OrderStatus != -1 { // If OrderStatus is not 0 (default proto value) or -1
		orderStatus = req.OrderStatus
	}

	// Update the payment result in the database
	err = o.repo.UpdateOrderPaymentResult(ctx, req.OrderId, req.PlatformStatus, orderStatus)
	if err != nil {
		o.logger.Log(log.LevelError, "msg", "Failed to update payment result", "order_id", req.OrderId, "error", err)
		return &ymcommonpb.ResponseBase{
			Code: "COMM9999", // general error
		}, nil
	}

	o.logger.Log(log.LevelInfo, "msg", "Payment result updated successfully", "order_id", req.OrderId)
	return &ymcommonpb.ResponseBase{
		Code: "COMM0000", // success
	}, nil
}

// GetOrderByPlatformId fetches an order by platform and payment IDs
func (o *orderServer) GetOrderByPlatformId(ctx context.Context, req *ymtransactionpb.SearchOrderByPlatformId) (*ymtransactionpb.GetOrderResponse, error) {
	o.logger.Log(log.LevelInfo, "msg", "Searching order by platform",
		"platform", req.Platform,
		"paymentId1", req.PaymentId1,
		"paymentId2", req.PaymentId2,
		"paymentId3", req.PaymentId3)

	// Search for the order in the repository
	order, err := o.repo.FindByPlatformId(ctx, req.Platform, req.PaymentId1, req.PaymentId2, req.PaymentId3)
	if err != nil {
		// If order not found, return COMM0000 with no order data
		if err == mongo.ErrNoDocuments {
			o.logger.Log(log.LevelInfo, "msg", "No order found for the given platform and payment IDs",
				"platform", req.Platform,
				"paymentId1", req.PaymentId1,
				"paymentId2", req.PaymentId2,
				"paymentId3", req.PaymentId3)
			return &ymtransactionpb.GetOrderResponse{
				Code: "COMM0000",
			}, nil
		}

		// If multiple orders found, return COMM0104
		if strings.Contains(err.Error(), "multiple orders found") {
			o.logger.Log(log.LevelWarn, "msg", "Multiple orders found for the given platform and payment IDs",
				"platform", req.Platform,
				"paymentId1", req.PaymentId1,
				"paymentId2", req.PaymentId2,
				"paymentId3", req.PaymentId3,
				"error", err)
			return &ymtransactionpb.GetOrderResponse{
				Code: "COMM0104", // multiple orders found
			}, nil
		}

		// For other errors, return COMM9999
		o.logger.Log(log.LevelError, "msg", "Failed to find order by platform",
			"platform", req.Platform,
			"error", err)
		return &ymtransactionpb.GetOrderResponse{
			Code: "COMM9999", // general error
		}, nil
	}

	// Convert the model to protobuf
	orderInfo := order.ToProto()

	o.logger.Log(log.LevelInfo, "msg", "Successfully found order by platform",
		"platform", req.Platform,
		"order_id", order.OrderId)

	// Return the success response with order info
	return &ymtransactionpb.GetOrderResponse{
		Code:  "COMM0000", // success
		Order: orderInfo,
	}, nil
}

// GetOrderById fetches an order by its ID
func (o *orderServer) GetOrderById(ctx context.Context, req *ymcommonpb.PKeyString) (*ymtransactionpb.GetOrderResponse, error) {
	orderID := req.Value
	o.logger.Log(log.LevelInfo, "msg", "Getting order by ID", "order_id", orderID)

	// Find the order by ID
	order, err := o.repo.FindByID(ctx, orderID)
	if err != nil {
		// If order not found, return COMM0103
		if err == mongo.ErrNoDocuments {
			o.logger.Log(log.LevelInfo, "msg", "Order not found", "order_id", orderID)
			return &ymtransactionpb.GetOrderResponse{
				Code: "COMM0103", // order not found
			}, nil
		}

		// For other errors, return COMM9999
		o.logger.Log(log.LevelError, "msg", "Failed to get order", "order_id", orderID, "error", err)
		return &ymtransactionpb.GetOrderResponse{
			Code: "COMM9999", // general error
		}, nil
	}

	// Convert the model to protobuf
	orderInfo := order.ToProto()

	o.logger.Log(log.LevelInfo, "msg", "Successfully found order", "order_id", orderID)

	// Return the success response with order info
	return &ymtransactionpb.GetOrderResponse{
		Code:  "COMM0000", // success
		Order: orderInfo,
	}, nil
}

// AddTransactionRecord adds transaction record and updates order
func (o *orderServer) AddTransactionRecord(ctx context.Context, req *ymtransactionpb.TransactionInfo) (*ymcommonpb.ResponseBase, error) {
	o.logger.Log(log.LevelInfo, "msg", "Adding transaction record", "order_id", req.OrderId, "transaction_id", req.TransactionId)

	// 1. Check if order exists
	_, err := o.repo.FindByID(ctx, req.OrderId)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			o.logger.Log(log.LevelInfo, "msg", "Order not found", "order_id", req.OrderId)
			return &ymcommonpb.ResponseBase{
				Code: "COMM0103", // order not found
			}, nil
		}
		o.logger.Log(log.LevelError, "msg", "Failed to find order", "order_id", req.OrderId, "error", err)
		return &ymcommonpb.ResponseBase{
			Code: "COMM9999", // general error
		}, nil
	}

	// Process differently based on whether TransactionId is provided
	isNewTransaction := false
	if req.TransactionId == "" {
		// No TransactionId provided, generate new transaction ID and create new record
		isNewTransaction = true
		req.TransactionId = o.node.Generate().String()
		o.logger.Log(log.LevelInfo, "msg", "Generated new transaction ID", "transaction_id", req.TransactionId)
	} else {
		// TransactionId provided, check if the transaction exists
		_, err := o.transactionRepo.FindByTransactionID(ctx, req.TransactionId)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				o.logger.Log(log.LevelInfo, "msg", "Transaction not found", "transaction_id", req.TransactionId)
				return &ymcommonpb.ResponseBase{
					Code: "COMM0103", // transaction not found
				}, nil
			}
			o.logger.Log(log.LevelError, "msg", "Failed to find transaction", "transaction_id", req.TransactionId, "error", err)
			return &ymcommonpb.ResponseBase{
				Code: "COMM9999", // general error
			}, nil
		}
		o.logger.Log(log.LevelInfo, "msg", "Found existing transaction, will update", "transaction_id", req.TransactionId)
	}

	// Create or update transaction record
	transaction := model.TransactionFromProto(req)

	var operationErr error
	if isNewTransaction {
		// Create new transaction record
		operationErr = o.transactionRepo.CreateTransaction(ctx, transaction)
		o.logger.Log(log.LevelInfo, "msg", "Creating new transaction", "transaction_id", req.TransactionId)
	} else {
		// Update existing transaction record
		operationErr = o.transactionRepo.UpdateTransaction(ctx, transaction)
		o.logger.Log(log.LevelInfo, "msg", "Updating existing transaction", "transaction_id", req.TransactionId)
	}

	if operationErr != nil {
		o.logger.Log(log.LevelError, "msg", "Failed to process transaction", "transaction_id", req.TransactionId, "error", operationErr)
		return &ymcommonpb.ResponseBase{
			Code: "COMM9999", // general error
		}, nil
	}

	// Create transaction log (regardless of whether it's a new creation or update)
	snowflakeID := o.node.Generate().String()
	transactionLog := model.TransactionLogFromTransaction(transaction, snowflakeID)
	if err := o.transactionRepo.CreateTransactionLog(ctx, transactionLog); err != nil {
		o.logger.Log(log.LevelError, "msg", "Failed to create transaction log", "transaction_id", req.TransactionId, "error", err)
		// Continue execution because transaction record processing is complete, log creation failure should not stop the entire process
		o.logger.Log(log.LevelWarn, "msg", "Continuing despite transaction log creation failure")
	}

	// Update order's CurrentPeriodStart and CurrentPeriodEnd
	var periodStart, periodEnd int64
	if req.PeriodStart != nil {
		periodStart = req.PeriodStart.AsTime().Unix()
	}
	if req.PeriodEnd != nil {
		periodEnd = req.PeriodEnd.AsTime().Unix()
	}

	if err := o.repo.UpdateOrderPeriod(ctx, req.OrderId, periodStart, periodEnd); err != nil {
		o.logger.Log(log.LevelError, "msg", "Failed to update order period", "order_id", req.OrderId, "error", err)
		return &ymcommonpb.ResponseBase{
			Code: "COMM9999", // general error
		}, nil
	}

	operationType := "added"
	if !isNewTransaction {
		operationType = "updated"
	}

	o.logger.Log(log.LevelInfo, "msg", fmt.Sprintf("Transaction record %s successfully", operationType), "order_id", req.OrderId, "transaction_id", req.TransactionId)
	return &ymcommonpb.ResponseBase{
		Code: "COMM0000", // success
	}, nil
}
