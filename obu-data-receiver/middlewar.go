package main

import (
	"time"

	"github.com/mohamedramadan14/roads-fees-system/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) *LogMiddleware {
	return &LogMiddleware{next: next}
}

func (l *LogMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(start).String(),
		}).Info("Producing OBU data to Kafka")
	}(time.Now())
	return l.next.ProduceData(data)
}

func (l *LogMiddleware) Close() error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start).String(),
		}).Info("Closing Kafka producer")
	}(time.Now())
	err := l.next.Close()
	if err != nil {
		logrus.Errorf("Error while closing Kafka producer: %v", err)
	}
	return err
}
