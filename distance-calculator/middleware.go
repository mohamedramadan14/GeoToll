package main

import (
	"time"

	"github.com/mohamedramadan14/roads-fees-system/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next DistanceServicer
}

func NewLogMiddleware(next DistanceServicer) *LogMiddleware {
	return &LogMiddleware{next: next}
}

func (m *LogMiddleware) CalculateDistance(data1, data2 types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID1":   data1.OBUID,
			"took":     time.Since(start).String(),
			"err":      err,
			"distance": dist,
		}).Info("Calculating distance")
	}(time.Now())

	return m.next.CalculateDistance(data1, data2)
}
