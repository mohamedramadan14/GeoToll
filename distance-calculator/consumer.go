package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mohamedramadan14/roads-fees-system/aggregator/client"
	"github.com/mohamedramadan14/roads-fees-system/types"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

const kafkaBroker = "localhost:9092"
const kafkaTopic = "obu-data"

var options = []kgo.Opt{
	kgo.SeedBrokers(kafkaBroker),
}

var pointsStore = make(map[int]types.OBUData)

type KafkaConsumer struct {
	consumer        *kgo.Client
	service         DistanceServicer
	aggregateClient *client.Client
}

func NewKafkaConsumer(topic string, svc DistanceServicer, client *client.Client) (*KafkaConsumer, error) {
	c, err := kgo.NewClient(options...)
	if topic != "" {
		c.AddConsumeTopics(topic)
	}
	if err != nil {
		logrus.Fatal("Failed to create Kafka consumer: ", err)
		return nil, err
	}

	logrus.Info("Distance calculator service started")

	// create topic if not exists
	logrus.Info("Waiting for OBU messages...")

	return &KafkaConsumer{
		consumer:        c,
		service:         svc,
		aggregateClient: client,
	}, nil
}

func (kc *KafkaConsumer) ConsumeMessageLoop(ctx context.Context) {
	for {
		fetches := kc.consumer.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				logrus.Error("Error fetching topic messages: ", err)
				continue
			}
		}
		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			logrus.Infof("Processing OBU message record: %+v\n", string(record.Value))
			var data types.OBUData
			if err := json.Unmarshal(record.Value, &data); err != nil {
				logrus.Error("Failed to unmarshal OBU data: ", err)
				continue
			}
			exitstData, ok := pointsStore[data.OBUID]
			if ok {
				logrus.Infof("OBU ID: %d already exists, updating data", data.OBUID)
				distance, err := kc.service.CalculateDistance(exitstData, data)
				if err != nil {
					logrus.Errorf("Failed to calculate distance of OBU ID: %d due to %v", data.OBUID, err)
					continue
				}
				logrus.Infof("Distance calculated for OBU ID: %d is %f", data.OBUID, distance)
				invoiceData := types.Distance{
					Value: distance,
					OBUID: data.OBUID,
					Unix:  time.Now().UnixNano(),
				}

				if err = kc.aggregateClient.AggregateInvoice(invoiceData); err != nil {
					logrus.Errorf("Failed to aggregate invoice for OBU ID: %d due to %v", data.OBUID, err)
				} else {
					logrus.Infof("Invoice aggregated for OBU ID: %d", data.OBUID)
				}
			}
			pointsStore[data.OBUID] = data
		}
	}
}

func (kf *KafkaConsumer) Start() {
	logrus.Info("Starting Kafka consumer...")
	kf.ConsumeMessageLoop(context.Background())
	logrus.Info("Kafka consumer started")
}
