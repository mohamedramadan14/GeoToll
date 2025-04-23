package main

import (
	"time"

	"github.com/mohamedramadan14/roads-fees-system/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{next: next}
}

func (lm *LogMiddleware) AggregateDistances(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": distance.OBUID,
			"value": distance.Value,
			"took":  time.Since(start).String(),
			"err":   err,
		}).Info("Aggregating distance")
	}(time.Now())

	return lm.next.AggregateDistances(distance)
}
