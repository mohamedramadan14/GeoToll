package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mohamedramadan14/roads-fees-system/types"
	"github.com/mohamedramadan14/roads-fees-system/utilities"
	"github.com/sirupsen/logrus"
)

const sendInterval = time.Second
const wsEndpoint = "ws://localhost:8085/ws"
const reconnectInterval = 3 * time.Second

var (
	sequence int = 0
)

func generateCoordinate() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func sendOBUData(conn *websocket.Conn, obuData types.OBUData) error {
	err := conn.WriteJSON(obuData)
	if err != nil {
		return fmt.Errorf("failed to send OBU data: %w", err)
	}
	return nil
}

func generateLocation() (float64, float64) {
	lat := generateCoordinate()
	long := generateCoordinate()
	return lat, long
}

func seedInit() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func websocketConnect() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	return conn, nil
}

func main() {
	utilities.InitLogger()
	seedInit()
	var conn *websocket.Conn
	for {
		if conn == nil {
			var err error
			conn, err = websocketConnect()
			if err != nil {
				logrus.Errorf("Failed to connect to WebSocket: %v", err)
				time.Sleep(reconnectInterval)
				continue
			}
			logrus.Infof("Successfully connected to WebSocket to %s", wsEndpoint)
		}
		lat, long := generateLocation()
		sequence++

		obuData := types.OBUData{
			OBUID: sequence,
			Lat:   lat,
			Long:  long,
		}

		err := sendOBUData(conn, obuData)
		if err != nil {
			logrus.Errorf("Error sending OBU data: %v", err)
			conn.Close()
			conn = nil
			logrus.Warnf("Connection closed. Will reconnect on next iteration.")
			continue
		}
		logrus.Infof("Sent OBU data: %+v", obuData)
		time.Sleep(sendInterval)
	}
}
