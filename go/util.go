package main

import (
	"fmt"
	"isucondition/pkg/kutil/wait"
	"go.uber.org/zap"
	"time"
)

func RetryWithTimeout(f func() error, interval, timeout time.Duration, logger *zap.SugaredLogger) error  {
	var times int
	err := wait.PollImmediate(interval, timeout, func() (done bool, err error) {
		if err := f(); err != nil {
			times++
			logger.Error(fmt.Sprintf("failed %v times: ", times), err)
			return false, nil
		}
		return true, nil
	})
	return err
}


// GetValidPagination pagination
func GetValidPagination(total, offset, limit int) (startIndex, endIndex int) {
	// no pagination
	if limit == 0 {
		return 0, total
	}

	// out of range
	if limit < 0 || offset < 0 || offset > total {
		return 0, 0
	}

	startIndex = offset
	endIndex = startIndex + limit

	if endIndex > total {
		endIndex = total
	}

	return startIndex, endIndex
}