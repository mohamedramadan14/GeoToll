package main

import (
	"github.com/mohamedramadan14/roads-fees-system/aggregator/client"
	"github.com/mohamedramadan14/roads-fees-system/utilities"
	"github.com/sirupsen/logrus"
)

func main() {
	utilities.InitLogger()
	logrus.Info("Starting distance calculator service...")

	svc := NewDistanceService()
	svcWithLogging := NewLogMiddleware(svc)
	aggregatorClient := client.NewClient("http://localhost:3100/aggregate")
	c, err := NewKafkaConsumer(kafkaTopic, svcWithLogging, aggregatorClient)

	if err != nil {
		logrus.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer c.consumer.Close()
	c.Start()
}
