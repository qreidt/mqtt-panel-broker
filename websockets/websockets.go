package websockets

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

type ReverbConn websocket.Conn

var WebsocketConn *websocket.Conn

type EventType string

const (
	MqttPacketReceived     EventType = "MqttPacketReceived"
	MqttClientConnected              = "MqttClientConnected"
	MqttClientDisconnected           = "MqttClientDisconnected"
)

func Init() {
	reverbHost := flag.String("reverb-host", "127.0.0.1:8080", "reverb service address")
	reverbAppKey := flag.String("reverb-app-key", "8jtblx730rmylh68ipdx", "reverb service address")

	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *reverbHost, Path: "/app/" + *reverbAppKey}
	log.Printf("connecting to %s", u.String())

	var err error
	if WebsocketConn, _, err = websocket.DefaultDialer.Dial(u.String(), nil); err != nil {
		log.Fatal("dial:", err)
	}
}

func Close() error {
	return WebsocketConn.Close()
}

func SendMessage(eventType EventType, data any) {
	message, _ := json.Marshal(map[string]interface{}{
		"channel": "mqtt",
		"event":   eventType,
		"data":    data,
	})

	if err := WebsocketConn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println("write:", err)
		return
	}
}

func sendPing() {
	message := []byte("{\"event\": \"pusher:ping\", \"data\": {}}")
	if err := WebsocketConn.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Println("ping:", err)
		return
	}
}
