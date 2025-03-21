package errors

import (
	"github.com/go-kratos/kratos/v2/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrEmailAlreadyExists = errors.New(400, "EMAIL_ALREADY_EXISTS", "email already exists")

func ToRpcError(err error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return status.Errorf(codes.NotFound, "Record not found")
	}
	if errors.Is(err, ErrEmailAlreadyExists) {
		return status.Errorf(codes.AlreadyExists, "Email already exists")
	}
	return status.Errorf(codes.Internal, "Internal server error: %v", err)
}
