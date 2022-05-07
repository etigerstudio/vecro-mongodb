package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

func loggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next BaseService) BaseService {
		return logmw{logger, next}
	}
}

type logmw struct {
	logger log.Logger
	BaseService
}

func (mw logmw) Execute() (result string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	result, err = mw.BaseService.Execute()
	return
}