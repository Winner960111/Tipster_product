package repository

import (
	"context"

	"src/internal/model"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionRepository defines transaction repository interface
type TransactionRepository interface {
	// CreateTransaction creates a new transaction record
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	// CreateTransactionLog creates a new transaction log record
	CreateTransactionLog(ctx context.Context, transactionLog *model.TransactionLog) error
	// FindByUserID finds transaction records by user ID
	FindByUserID(ctx context.Context, userID string, pageIndex, pageSize int32) ([]model.Transaction, int64, error)
	// FindByTransactionID finds transaction record by transaction ID
	FindByTransactionID(ctx context.Context, transactionID string) (*model.Transaction, error)
	// UpdateTransaction updates transaction record
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) error
	// AddLogReferenceToTransaction adds a log reference to a transaction
	AddLogReferenceToTransaction(ctx context.Context, transactionID, logID string) error
	// GetTransactionLogs retrieves all logs for a transaction
	GetTransactionLogs(ctx context.Context, transactionID string) ([]model.TransactionLog, error)
}

// transactionRepository implements TransactionRepository interface
type transactionRepository struct {
	transactionCollection *mongo.Collection
	logCollection         *mongo.Collection
	logger                log.Logger
}

// NewTransactionRepository creates a new TransactionRepository instance
func NewTransactionRepository(db *mongo.Database, logger log.Logger) TransactionRepository {
	transactionCollection := db.Collection("transactions")
	logCollection := db.Collection("transaction_logs")

	// Create indexes for transactions collection
	transactionIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "transactionId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "orderId", Value: 1}},
		},
	}

	_, err := transactionCollection.Indexes().CreateMany(context.Background(), transactionIndexes)
	if err != nil {
		logger.Log(log.LevelError, "msg", "Failed to create transaction indexes", "error", err)
	}

	// Create indexes for transaction_logs collection
	logIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "logId", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "transactionId", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "orderId", Value: 1}},
		},
	}

	_, err = logCollection.Indexes().CreateMany(context.Background(), logIndexes)
	if err != nil {
		logger.Log(log.LevelError, "msg", "Failed to create transaction log indexes", "error", err)
	}

	return &transactionRepository{
		transactionCollection: transactionCollection,
		logCollection:         logCollection,
		logger:                logger,
	}
}

// CreateTransaction implements TransactionRepository's CreateTransaction method
func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction *model.Transaction) error {
	_, err := r.transactionCollection.InsertOne(ctx, transaction)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to create transaction", "transaction_id", transaction.TransactionId, "error", err)
		return err
	}

	r.logger.Log(log.LevelInfo, "msg", "Transaction created successfully", "transaction_id", transaction.TransactionId)
	return nil
}

// CreateTransactionLog implements TransactionRepository's CreateTransactionLog method
func (r *transactionRepository) CreateTransactionLog(ctx context.Context, transactionLog *model.TransactionLog) error {
	// Insert the transaction log
	_, err := r.logCollection.InsertOne(ctx, transactionLog)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to create transaction log", "transaction_id", transactionLog.TransactionId, "log_id", transactionLog.LogId, "error", err)
		return err
	}

	// Add the log reference to the transaction
	err = r.AddLogReferenceToTransaction(ctx, transactionLog.TransactionId, transactionLog.LogId)
	if err != nil {
		r.logger.Log(log.LevelWarn, "msg", "Created log but failed to update transaction reference", "transaction_id", transactionLog.TransactionId, "log_id", transactionLog.LogId, "error", err)
		// Continue despite the error, as the log itself was created successfully
	}

	r.logger.Log(log.LevelInfo, "msg", "Transaction log created successfully", "transaction_id", transactionLog.TransactionId, "log_id", transactionLog.LogId)
	return nil
}

// FindByUserID finds transaction records by user ID
func (r *transactionRepository) FindByUserID(ctx context.Context, userID string, pageIndex, pageSize int32) ([]model.Transaction, int64, error) {
	// Currently using empty implementation because our Transaction model doesn't have UserID field
	// Based on actual requirements, you should implement proper query logic
	// For example, you might need to query orders for this user first, then query transactions related to these orders

	r.logger.Log(log.LevelInfo, "msg", "Finding transactions by user ID", "user_id", userID)

	// Return empty list and count 0
	return []model.Transaction{}, 0, nil
}

// FindByTransactionID finds transaction record by transaction ID
func (r *transactionRepository) FindByTransactionID(ctx context.Context, transactionID string) (*model.Transaction, error) {
	var transaction model.Transaction
	filter := bson.M{"transactionId": transactionID}

	err := r.transactionCollection.FindOne(ctx, filter).Decode(&transaction)
	if err != nil {
		r.logger.Log(log.LevelWarn, "msg", "Transaction not found", "transaction_id", transactionID, "error", err)
		return nil, err
	}

	r.logger.Log(log.LevelInfo, "msg", "Transaction found", "transaction_id", transactionID)
	return &transaction, nil
}

// UpdateTransaction updates transaction record
func (r *transactionRepository) UpdateTransaction(ctx context.Context, transaction *model.Transaction) error {
	filter := bson.M{"transactionId": transaction.TransactionId}
	update := bson.M{"$set": transaction}

	result, err := r.transactionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to update transaction", "transaction_id", transaction.TransactionId, "error", err)
		return err
	}

	if result.MatchedCount == 0 {
		r.logger.Log(log.LevelWarn, "msg", "No transaction found to update", "transaction_id", transaction.TransactionId)
		return mongo.ErrNoDocuments
	}

	r.logger.Log(log.LevelInfo, "msg", "Transaction updated successfully", "transaction_id", transaction.TransactionId)
	return nil
}

// AddLogReferenceToTransaction adds a log reference to a transaction
func (r *transactionRepository) AddLogReferenceToTransaction(ctx context.Context, transactionID, logID string) error {
	// Use $addToSet to avoid duplicates
	filter := bson.M{"transactionId": transactionID}
	update := bson.M{"$addToSet": bson.M{"logIds": logID}}

	result, err := r.transactionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to add log reference to transaction", "transaction_id", transactionID, "log_id", logID, "error", err)
		return err
	}

	if result.MatchedCount == 0 {
		r.logger.Log(log.LevelWarn, "msg", "No transaction found to add log reference", "transaction_id", transactionID, "log_id", logID)
		return mongo.ErrNoDocuments
	}

	r.logger.Log(log.LevelInfo, "msg", "Log reference added to transaction successfully", "transaction_id", transactionID, "log_id", logID)
	return nil
}

// GetTransactionLogs retrieves all logs for a transaction
func (r *transactionRepository) GetTransactionLogs(ctx context.Context, transactionID string) ([]model.TransactionLog, error) {
	var logs []model.TransactionLog
	filter := bson.M{"transactionId": transactionID}

	// Sort by date created descending to get newest logs first
	opts := options.Find().SetSort(bson.D{{Key: "dateCreated", Value: -1}})

	cursor, err := r.logCollection.Find(ctx, filter, opts)
	if err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to retrieve transaction logs", "transaction_id", transactionID, "error", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &logs); err != nil {
		r.logger.Log(log.LevelError, "msg", "Failed to decode transaction logs", "transaction_id", transactionID, "error", err)
		return nil, err
	}

	r.logger.Log(log.LevelInfo, "msg", "Retrieved transaction logs successfully", "transaction_id", transactionID, "log_count", len(logs))
	return logs, nil
}
