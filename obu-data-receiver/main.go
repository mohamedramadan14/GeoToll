package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mohamedramadan14/roads-fees-system/types"
	"github.com/mohamedramadan14/roads-fees-system/utilities"
	"github.com/sirupsen/logrus"
)

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  DataProducer
}

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p     DataProducer
		err   error
		topic = "obu-data"
	)
	p, err = NewKafkaProducer(topic)
	if err != nil {
		logrus.Errorf("Failed to create Kafka producer: %v", err)
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		prod:  p,
	}, nil
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("Error while upgrading connection: %v", err)
	}
	dr.conn = conn
	go dr.readMessagesContinously()
}

func (dr *DataReceiver) produceDataToKafka(message types.OBUData) error {
	return dr.prod.ProduceData(message)
}

func (dr *DataReceiver) readMessagesContinously() {
	logrus.Infof("Starting to read messages from WebSocket connection...")
	for {
		var message types.OBUData
		if err := dr.conn.ReadJSON(&message); err != nil {
			logrus.Errorf("Error while reading message: %v", err)
			continue
		}

		if err := dr.produceDataToKafka(message); err != nil {
			logrus.Errorf("Error while producing message to Kafka: %v", err)
			continue
		}
	}
}

func main() {

	utilities.InitLogger()
	dataReceiver, err := NewDataReceiver()
	if err != nil {
		logrus.Errorf("Failed to create DataReceiver: %v", err)
		return
	}

	logrus.Info("OBU data receiver Microservice is starting...")
	// defer dataReceiver.kafkaClient.Close()
	defer dataReceiver.prod.Close()

	http.HandleFunc("/ws", dataReceiver.handleWS)
	logrus.Error(http.ListenAndServe(":8085", nil))
	logrus.Info("OBU data receiver Microservice is running on :8085")
}
