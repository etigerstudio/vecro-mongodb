package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"time"
)

type BaseService interface {
	Execute() error
}

type baseService struct{
	callees []endpoint.Endpoint  // Downstream endpoints to be called on
	isSynchronous bool  // Whether to call services synchronously or asynchronously.  TODO: NOT IMPLEMENTED YET!
	serviceType ServiceType
}

type ServiceType int

const (
	vanilla ServiceType = 0
	cpu ServiceType = 1
	io ServiceType = 2
)

func (svc baseService) Execute() error {
	// Establish the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Simulate Stress
	// TODO: Re-determine timings of stress
	stress(svc.serviceType)

	// Call downstream services
	for _, ep := range svc.callees {
		// TODO: Pass parameters and parse responses
		// TODO: Support async calling mode
		_, err := ep(ctx, nil)
		if err != nil {
			return err
		}
	}

	// Return result
	return nil
}

type ServiceMiddleware func(BaseService) BaseService
