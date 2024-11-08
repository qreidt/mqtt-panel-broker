package hooks

import (
	"broker-manager/websockets"
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"time"
)

type PacketType byte

//goland:noinspection SpellCheckingInspection
const (
	CONNECT     PacketType = iota + 1 // The MQTT client says "Hello, I want to connect!"
	CONNACK                           // MQTT Broker replies, "Hello, you are connected!"
	PUBLISH                           // The MQTT client or Broker sends a message.
	PUBACK                            // Receiver says, "Got your message!"
	PUBREC                            // Receiver says, "Received your message, will process."
	PUBREL                            // Sender says, "Please release the message."
	PUBCOMP                           // Receiver says, "Message processing complete."
	SUBSCRIBE                         // The MQTT client says, "I want to receive updates about this topic."
	SUBACK                            // The MQTT Broker replies, "You are subscribed."
	UNSUBSCRIBE                       // The MQTT client says, "I don't want updates about this topic anymore."
	UNSUBACK                          // The MQTT Broker replies, "You are unsubscribed."
	PINGREQ                           // The MQTT client says, "Are you there?"
	PINGRESP                          // The MQTT Broker replies, "Yes, I'm here."
	DISCONNECT                        // The MQTT client says, "Goodbye, I'm disconnecting."
)

type ClientPacketProcessedEvent struct {
	ID           string     `json:"id"`
	PacketId     uint16     `json:"packet_id"`
	PacketType   PacketType `json:"packet_type"`
	PacketLength uint       `json:"packet_length"`
	Timestamp    uint64     `json:"timestamp"`
}

type OnPacketProcessed struct {
	mqtt.HookBase
}

// ID returns the ID of the hook.
func (h *OnPacketProcessed) ID() string {
	return "on-packet-processed"
}

// Provides indicates which hook methods this hook provides.
func (h *OnPacketProcessed) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnPacketProcessed,
		mqtt.OnPacketSent,
	}, []byte{b})
}

func (h *OnPacketProcessed) OnPacketSent(cl *mqtt.Client, pk packets.Packet, b []byte) {
	event := ClientPacketProcessedEvent{
		ID:           cl.ID,
		PacketId:     pk.PacketID,
		PacketType:   PacketType(pk.FixedHeader.Type),
		PacketLength: uint(pk.FixedHeader.Remaining),
		Timestamp:    uint64(time.Now().UnixMilli()),
	}

	h.Log.Info("Packet Processed", "event", event)
	websockets.SendMessage(websockets.MqttPacketProcessed, event)
}

// OnPacketProcessed Intercepts the disconnected client and generates an event to be sent on the websocket
func (h *OnPacketProcessed) OnPacketProcessed(cl *mqtt.Client, pk packets.Packet, err error) {
	event := ClientPacketProcessedEvent{
		ID:           cl.ID,
		PacketId:     pk.PacketID,
		PacketType:   PacketType(pk.FixedHeader.Type),
		PacketLength: uint(pk.FixedHeader.Remaining),
		Timestamp:    uint64(time.Now().UnixMilli()),
	}

	h.Log.Info("Packet Processed", "event", event)
	websockets.SendMessage(websockets.MqttPacketProcessed, event)
}
