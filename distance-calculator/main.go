package main

import (
	"github.com/mohamedramadan14/roads-fees-system/utilities"
	"github.com/sirupsen/logrus"
)

func main() {
	utilities.InitLogger()
	logrus.Info("Starting distance calculator service...")

	svc := NewDistanceService()
	svcWithLogging := NewLogMiddleware(svc)

	c, err := NewKafkaConsumer(kafkaTopic, svcWithLogging)
	if err != nil {
		logrus.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer c.consumer.Close()
	c.Start()
}
