package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/mohamedramadan14/roads-fees-system/types"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

const kafkaTopic = "obu-data"
const kafkaBroker = "localhost:9092"

var options = []kgo.Opt{
	kgo.SeedBrokers(kafkaBroker),
	kgo.AllowAutoTopicCreation(),
	kgo.RetryBackoffFn(func(attempt int) time.Duration {
		return time.Second * time.Duration(attempt)
	}),
	kgo.RecordRetries(5),
}

type DataProducer interface {
	ProduceData(data types.OBUData) error
	Close() error
}

type KafkaProducer struct {
	producer *kgo.Client
	topic    string
}

func NewKafkaProducer(topic string) (DataProducer, error) {
	client, err := kgo.NewClient(options...)
	if err != nil {
		logrus.Errorf("Failed to create Kafka client: %v", err)
		return nil, err
	} else {
		logrus.Info("Kafka client created successfully")
	}
	if topic == "" {
		return &KafkaProducer{producer: client, topic: kafkaTopic}, nil
	}
	return &KafkaProducer{producer: client, topic: topic}, nil
}

func (kp *KafkaProducer) ProduceData(data types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("Error while marshalling message to JSON: %v", err)
		return err
	}
	record := &kgo.Record{
		Topic: kp.topic,
		Value: b,
	}
	err = kp.producer.ProduceSync(context.Background(), record).FirstErr()
	if err != nil {
		logrus.Errorf("Error while producing message to Kafka: %v", err)
		return err
	}

	partition := record.Partition
	offset := record.Offset
	logrus.Infof("Message delivered successfully to Kafka - Topic: %s, Partition: %d, Offset: %d", kafkaTopic, partition, offset)

	return nil
}

func (kp *KafkaProducer) Close() error {
	kp.producer.Close()
	return nil
}
