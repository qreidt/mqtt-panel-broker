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

var WebsocketConn *websocket.Conn

type EventType string

const (
	MqttPacketReceived  EventType = "MqttPacketReceived"
	MqttClientConnected           = "MqttClientConnected"
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

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	WebsocketConn = c

	if err != nil {
		log.Fatal("dial:", err)
	}

	//goland:noinspection GoUnhandledErrorResult
	defer c.Close()
}

func SendMessage(eventType EventType, data map[string]interface{}) {
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
