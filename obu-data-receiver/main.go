package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mohamedramadan14/roads-fees-system/types"
)

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
}

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		conn:  nil,
	}
}
func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error while upgrading connection: ", err)
	}
	dr.conn = conn
	go dr.readMessagesContinously()
}

func (dr *DataReceiver) readMessagesContinously() {
	fmt.Printf("Websockt Connected to %s\n", dr.conn.RemoteAddr())
	for {
		var message types.OBUData
		if err := dr.conn.ReadJSON(&message); err != nil {
			log.Printf("Error while reading message: %v", err)
			continue
		}
		dr.msgch <- message
		fmt.Printf("Received OBU Data from [%d] :: <Lat %.8f, Long %.8f> \n", message.OBUID, message.Lat, message.Long)
	}
}
func main() {
	fmt.Println("This is the OBU data receiver Microservice.")
	dataReceiver := NewDataReceiver()
	http.HandleFunc("/ws", dataReceiver.handleWS)
	log.Fatal(http.ListenAndServe(":8085", nil))
	fmt.Println("Server started Successfully on :8085")
}
