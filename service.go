package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"
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

	lastResult := ""

	// Perform read operations
	for i := 0; i < svc.dbReadOps; i++ {
		id := rand.Intn(itemsCount)
		filter := bson.D{{"id", id}}
		fmt.Printf("db read id: %d\n", id)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := svc.dbCollection.FindOne(ctx, filter).Decode(&result)
		fmt.Printf("db read: %d\n", id)
		if err == mongo.ErrNoDocuments {
			return "id not found", nil
		} else if err != nil {
			return "failed to find in the collection", nil
		}

		lastResult = fmt.Sprintf("%f", result.Value)
		fmt.Printf("db read result: %s\n", lastResult)
	}

	// Perform write operations
	for i := 0; i < svc.dbWriteOps; i++ {
		document := bson.D{{"app", "vecro-sim"}, {"rand_value", rand.Float64()}}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		result, err := svc.dbCollection.InsertOne(ctx, document)
		if err != nil {
			return "failed to insert into the collection", nil
		}

		lastResult = fmt.Sprintf("%v", result.InsertedID)
		fmt.Printf("db insert result: %s\n", lastResult)
	}

	return lastResult, nil
}

type ServiceMiddleware func(BaseService) BaseService
