package model

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	pb "src/protos/YM.Transaction"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// OrderModel represents the MongoDB document structure for orders
type Order struct {
	ID                 string               `bson:"_id,omitempty"`
	OrderId            string               `bson:"orderId"`
	OrderCode          string               `bson:"orderCode"`
	MemberId           string               `bson:"memberId"`
	PackageCode        string               `bson:"packageCode"`
	Status             int32                `bson:"status"`
	LocalPrice         primitive.Decimal128 `bson:"localPrice"`
	LocalCurrency      string               `bson:"localCurrency"`
	SubscriptionType   int32                `bson:"subscriptionType"`
	BuyIn              []MerchandiseModel   `bson:"buyIn"`
	DateCreated        int64                `bson:"dateCreated"` // Unix timestamp in milliseconds
	DateUpdated        int64                `bson:"dateUpdated"` // Unix timestamp in milliseconds
	Payment            *PaymentInfoModel    `bson:"payment,omitempty"`
	CurrentPeriodStart *int64               `bson:"currentPeriodStart,omitempty"` // Unix timestamp in milliseconds
	CurrentPeriodEnd   *int64               `bson:"currentPeriodEnd,omitempty"`   // Unix timestamp in milliseconds
}

type MerchandiseModel struct {
	Type   int32 `bson:"type"`
	Amount int32 `bson:"amount"`
}

type PaymentInfoModel struct {
	Platform    int32  `bson:"platform"`
	PaymentId1  string `bson:"paymentId1,omitempty"`
	PaymentId2  string `bson:"paymentId2,omitempty"`
	PaymentId3  string `bson:"paymentId3,omitempty"`
	PaymentLink string `bson:"paymentLink,omitempty"`
	Status      string `bson:"status,omitempty"`
}

// OrderRepo handles the data access layer for orders
type OrderRepo struct {
	collection *mongo.Collection
}

// NewOrderRepo creates a new OrderRepo instance
func NewOrderRepo(db *mongo.Database) *OrderRepo {
	collection := db.Collection("orders")

	// Create indexes
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "orderId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "orderCode", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "payment.platform", Value: 1},
				{Key: "payment.paymentId1", Value: 1},
				{Key: "payment.paymentId2", Value: 1},
				{Key: "payment.paymentId3", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetPartialFilterExpression(
				bson.D{{Key: "payment", Value: bson.D{{Key: "$exists", Value: true}}}},
			),
		},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		panic(err)
	}

	return &OrderRepo{
		collection: collection,
	}
}

// CreateOrder creates a new order in MongoDB
func (r *OrderRepo) CreateOrder(ctx context.Context, order *pb.OrderInfo) error {
	decimalPrice, err := primitive.ParseDecimal128(order.LocalPrice)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	model := &Order{
		OrderId:          order.OrderId,
		OrderCode:        order.OrderCode,
		MemberId:         order.MemberId,
		PackageCode:      order.PackageCode,
		Status:           order.Status,
		LocalPrice:       decimalPrice,
		LocalCurrency:    order.LocalCurrency,
		SubscriptionType: order.SubscriptionType,
		BuyIn:            make([]MerchandiseModel, len(order.BuyIn)),
		DateCreated:      now,
		DateUpdated:      now,
	}

	// Convert BuyIn items
	for i, item := range order.BuyIn {
		model.BuyIn[i] = MerchandiseModel{
			Type:   item.Type,
			Amount: item.Amount,
		}
	}

	_, err = r.collection.InsertOne(ctx, model)
	return err
}

// ToProto converts OrderModel to protobuf OrderInfo
func (m *Order) ToProto() *pb.OrderInfo {
	order := &pb.OrderInfo{
		OrderId:          m.OrderId,
		OrderCode:        m.OrderCode,
		MemberId:         m.MemberId,
		PackageCode:      m.PackageCode,
		Status:           m.Status,
		LocalPrice:       m.LocalPrice.String(),
		LocalCurrency:    m.LocalCurrency,
		SubscriptionType: m.SubscriptionType,
		BuyIn:            make([]*pb.Merchandise, len(m.BuyIn)),
		DateCreated:      timestamppb.New(time.UnixMilli(m.DateCreated)),
		DateUpdated:      timestamppb.New(time.UnixMilli(m.DateUpdated)),
	}

	// Convert BuyIn items
	for i, item := range m.BuyIn {
		order.BuyIn[i] = &pb.Merchandise{
			Type:   item.Type,
			Amount: item.Amount,
		}
	}

	// Convert Payment if exists
	if m.Payment != nil {
		order.Payment = &pb.PaymentInfo{
			Platform:    m.Payment.Platform,
			PaymentId1:  m.Payment.PaymentId1,
			PaymentId2:  m.Payment.PaymentId2,
			PaymentId3:  m.Payment.PaymentId3,
			PaymentLink: m.Payment.PaymentLink,
			Status:      m.Payment.Status,
		}
	}

	// Convert period times if they exist
	if m.CurrentPeriodStart != nil {
		order.CurrentPeriodStart = timestamppb.New(time.UnixMilli(*m.CurrentPeriodStart))
	}
	if m.CurrentPeriodEnd != nil {
		order.CurrentPeriodEnd = timestamppb.New(time.UnixMilli(*m.CurrentPeriodEnd))
	}

	return order
}
