package main

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

func instrumentingMiddleware(
	requestCount metrics.Counter,
	requestLatency metrics.Histogram,
) ServiceMiddleware {
	return func(next BaseService) BaseService {
		return instrmw{requestCount, requestLatency, next}
	}
}

type instrmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	BaseService
}

func (mw instrmw) Execute() (string, error) {
	defer func(begin time.Time) {
		mw.requestCount.Add(1)
		mw.requestLatency.Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.BaseService.Execute()
}
