package model

import (
	"time"

	pb "src/protos/YM.Transaction"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// Transaction represents transaction record
type Transaction struct {
	TransactionId string   `bson:"transactionId"`
	Platform      int32    `bson:"platform"`
	PlatformId    string   `bson:"platformId"`
	OrderId       string   `bson:"orderId"`
	Status        string   `bson:"status"`
	PaymentLink   string   `bson:"paymentLink"`
	PeriodStart   int64    `bson:"periodStart"`      // Unix timestamp in milliseconds
	PeriodEnd     int64    `bson:"periodEnd"`        // Unix timestamp in milliseconds
	DateCreated   int64    `bson:"dateCreated"`      // Unix timestamp in milliseconds
	LogIds        []string `bson:"logIds,omitempty"` // References to transaction_logs documents
}

// ToProto converts Transaction to protobuf TransactionInfo
func (t *Transaction) ToProto() *pb.TransactionInfo {
	return &pb.TransactionInfo{
		TransactionId: t.TransactionId,
		Platform:      t.Platform,
		PlatformId:    t.PlatformId,
		OrderId:       t.OrderId,
		Status:        t.Status,
		PaymentLink:   t.PaymentLink,
		PeriodStart:   timestamppb.New(time.UnixMilli(t.PeriodStart)),
		PeriodEnd:     timestamppb.New(time.UnixMilli(t.PeriodEnd)),
		DateCreated:   timestamppb.New(time.UnixMilli(t.DateCreated)),
	}
}

// FromProto creates Transaction from protobuf TransactionInfo
func TransactionFromProto(info *pb.TransactionInfo) *Transaction {
	t := &Transaction{
		TransactionId: info.TransactionId,
		Platform:      info.Platform,
		PlatformId:    info.PlatformId,
		OrderId:       info.OrderId,
		Status:        info.Status,
		PaymentLink:   info.PaymentLink,
		DateCreated:   time.Now().Unix(),
		LogIds:        []string{}, // Initialize empty array for log references
	}

	if info.PeriodStart != nil {
		t.PeriodStart = info.PeriodStart.AsTime().Unix()
	}

	if info.PeriodEnd != nil {
		t.PeriodEnd = info.PeriodEnd.AsTime().Unix()
	}

	return t
}

// TransactionLog represents transaction log record
type TransactionLog struct {
	LogId         string `bson:"logId"` // Unique ID for each log entry
	TransactionId string `bson:"transactionId"`
	Platform      int32  `bson:"platform"`
	PlatformId    string `bson:"platformId"`
	OrderId       string `bson:"orderId"`
	Status        string `bson:"status"`
	PaymentLink   string `bson:"paymentLink"`
	PeriodStart   int64  `bson:"periodStart"` // Unix timestamp in milliseconds
	PeriodEnd     int64  `bson:"periodEnd"`   // Unix timestamp in milliseconds
	DateCreated   int64  `bson:"dateCreated"` // Unix timestamp in milliseconds
}

// FromTransaction creates TransactionLog from Transaction
func TransactionLogFromTransaction(t *Transaction, snowflakeID string) *TransactionLog {
	// Generate a unique log ID using timestamp and snowflake ID
	logId := snowflakeID

	return &TransactionLog{
		LogId:         logId,
		TransactionId: t.TransactionId,
		Platform:      t.Platform,
		PlatformId:    t.PlatformId,
		OrderId:       t.OrderId,
		Status:        t.Status,
		PaymentLink:   t.PaymentLink,
		PeriodStart:   t.PeriodStart,
		PeriodEnd:     t.PeriodEnd,
		DateCreated:   time.Now().Unix(),
	}
}
