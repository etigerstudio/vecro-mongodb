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
	calls []endpoint.Endpoint  // Downstream endpoints to be called on
	isSynchronous bool  // Whether to call services synchronously or asynchronously.  TODO: NOT IMPLEMENTED YET!

	delayTime int
	delayJitter int
	cpuLoad int
	ioLoad int
}

func (svc baseService) Execute() error {
	// Establish the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Simulate Stress
	// TODO: Re-determine timings of stress
	stress(svc.delayTime, svc.delayJitter, svc.cpuLoad, svc.ioLoad)

	// Call downstream services
	for _, ep := range svc.calls {
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
