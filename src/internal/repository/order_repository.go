package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"src/internal/model"
)

type OrderRepository interface {
	FindByUserID(ctx context.Context, userID string, pageIndex, pageSize int32) ([]model.Order, int64, error)

	CreateOrder(ctx context.Context, order *model.Order) (string, error)
	FindByID(ctx context.Context, orderID string) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID string, newStatus string) error
	UpdateOrderPayment(ctx context.Context, orderID string, payment *model.PaymentInfoModel) error
	UpdateOrderPaymentResult(ctx context.Context, orderID string, platformStatus string, orderStatus int32) error
	FindByPlatformId(ctx context.Context, platform int32, paymentId1, paymentId2, paymentId3 string) (*model.Order, error)
	UpdateOrderPeriod(ctx context.Context, orderID string, periodStart, periodEnd int64) error
}

type orderRepository struct {
	collection *mongo.Collection
	logger     log.Logger
}

func NewOrderRepository(db *mongo.Database, logger log.Logger) OrderRepository {
	collection := db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Log(log.LevelError, "msg", "Failed to connect to collection orders", "error", err)
	} else {
		logger.Log(log.LevelInfo, "msg", "Successfully connected to collection orders")
	}

	return &orderRepository{
		collection: collection,
		logger:     logger,
	}
}

// FindByUserID queries orders by user ID with pagination
func (r *orderRepository) FindByUserID(ctx context.Context, userID string, pageIndex, pageSize int32) ([]model.Order, int64, error) {
	filter := bson.M{"user_id": userID}

	skip := int64(pageIndex * pageSize)
	limit := int64(pageSize)

	// Example: Sort by CreatedAt in descending order
	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to find orders", "error", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var orders []model.Order
	if err := cursor.All(ctx, &orders); err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to decode orders", "error", err)
		return nil, 0, err
	}

	// Calculate total count
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to count documents", "error", err)
		return nil, 0, err
	}

	return orders, totalCount, nil
}

func (r *orderRepository) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	if order.ID == "" {
		order.ID = primitive.NewObjectID().Hex()
	}

	_, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to create order", "error", err)
		return "", err
	}

	r.logger.Log(log.LevelInfo, "msg", "Order created successfully", "order_id", order.ID)
	return order.ID, nil
}

func (r *orderRepository) FindByID(ctx context.Context, orderID string) (*model.Order, error) {
	filter := bson.M{"orderId": orderID}

	var order model.Order
	if err := r.collection.FindOne(ctx, filter).Decode(&order); err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to find order", "order_id", orderID, "error", err)
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID string, newStatus string) error {
	filter := bson.M{"orderId": orderID}
	update := bson.M{"$set": bson.M{"status": newStatus, "dateUpdated": time.Now().Unix()}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to update order status", "order_id", orderID, "error", err)
		return err
	}

	if result.ModifiedCount == 0 {
		r.logger.Log(log.LevelWarn, "msg", "No order was updated", "order_id", orderID)
	} else {
		r.logger.Log(log.LevelInfo, "msg", "Order status updated successfully", "order_id", orderID)
	}

	return nil
}

func (r *orderRepository) UpdateOrderPayment(ctx context.Context, orderID string, payment *model.PaymentInfoModel) error {
	filter := bson.M{"orderId": orderID}
	update := bson.M{"$set": bson.M{"payment": payment, "dateUpdated": time.Now().Unix()}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to update order payment", "order_id", orderID, "error", err)
		return err
	}

	if result.ModifiedCount == 0 {
		r.logger.Log(log.LevelWarn, "msg", "No order was updated", "order_id", orderID)
	} else {
		r.logger.Log(log.LevelInfo, "msg", "Order payment updated successfully", "order_id", orderID)
	}

	return nil
}

func (r *orderRepository) UpdateOrderPaymentResult(ctx context.Context, orderID string, platformStatus string, orderStatus int32) error {
	filter := bson.M{"orderId": orderID}

	// Prepare update document
	update := bson.M{"$set": bson.M{"dateUpdated": time.Now().Unix()}}

	// Only update order status if it's not "none"
	if orderStatus != -1 {
		update["$set"].(bson.M)["status"] = orderStatus
	}

	// Only update payment status if it's not "none"
	if platformStatus != "none" {
		update["$set"].(bson.M)["payment.status"] = platformStatus
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to update order payment result", "order_id", orderID, "error", err)
		return err
	}

	if result.ModifiedCount == 0 {
		r.logger.Log(log.LevelWarn, "msg", "No order was updated", "order_id", orderID)
	} else {
		r.logger.Log(log.LevelInfo, "msg", "Order payment result updated successfully", "order_id", orderID)
	}

	return nil
}

// FindByPlatformId searches for an order by platform and optional payment IDs
func (r *orderRepository) FindByPlatformId(ctx context.Context, platform int32, paymentId1, paymentId2, paymentId3 string) (*model.Order, error) {
	// Build a filter with platform
	filter := bson.M{"payment.platform": platform}

	// Add non-empty payment IDs to the filter
	if paymentId1 != "" {
		filter["payment.paymentId1"] = paymentId1
	}
	if paymentId2 != "" {
		filter["payment.paymentId2"] = paymentId2
	}
	if paymentId3 != "" {
		filter["payment.paymentId3"] = paymentId3
	}

	r.logger.Log(log.LevelInfo, "msg", "Finding order by platform",
		"platform", platform,
		"paymentId1", paymentId1, "paymentId2", paymentId2, "paymentId3", paymentId3, "filter", filter)

	// First, check if there are multiple matches
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to count orders by platform", "platform", platform, "error", err)
		return nil, err
	}

	// If multiple orders found, return a custom error
	if count > 1 {
		r.logger.Log(log.LevelWarn, "msg", "Multiple orders found for the same platform and payment IDs",
			"platform", platform, "count", count)
		return nil, fmt.Errorf("multiple orders found: %d", count)
	}

	// If no orders found, return no documents error
	if count == 0 {
		r.logger.Log(log.LevelInfo, "msg", "No orders found for the given platform and payment IDs", "platform", platform)
		return nil, mongo.ErrNoDocuments
	}

	// If we get here, exactly one order was found
	var order model.Order
	if err := r.collection.FindOne(ctx, filter).Decode(&order); err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to find order by platform", "platform", platform, "error", err)
		return nil, err
	}

	r.logger.Log(log.LevelInfo, "msg", "Successfully found order by platform", "platform", platform, "order_id", order.OrderId)
	return &order, nil
}

// UpdateOrderPeriod updates order's CurrentPeriodStart and CurrentPeriodEnd
func (r *orderRepository) UpdateOrderPeriod(ctx context.Context, orderID string, periodStart, periodEnd int64) error {
	filter := bson.M{"orderId": orderID}
	update := bson.M{"$set": bson.M{
		"currentPeriodStart": periodStart,
		"currentPeriodEnd":   periodEnd,
		"dateUpdated":        time.Now().Unix(),
	}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to update order period", "order_id", orderID, "error", err)
		return err
	}

	if result.ModifiedCount == 0 {
		r.logger.Log(log.LevelWarn, "msg", "No order was updated", "order_id", orderID)
	} else {
		r.logger.Log(log.LevelInfo, "msg", "Order period updated successfully", "order_id", orderID)
	}

	return nil
}
