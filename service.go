package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BaseService interface {
	Execute() (string, error)
}

type baseService struct {
	dbCollection *mongo.Collection

	dbReadOps  int
	dbWriteOps int
}

const itemsCount = 100000

func (svc baseService) Execute() (string, error) {
	var result struct {
		Value float64
	}
	filter := bson.D{{"id", rand.Intn(itemsCount)}}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := svc.dbCollection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return "id not found", nil
	} else if err != nil {
		return "failed to find in the collection", nil
	}

	return fmt.Sprintf("%f", result.Value), nil
}

type ServiceMiddleware func(BaseService) BaseService
