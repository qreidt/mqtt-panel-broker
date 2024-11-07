package hooks

import (
	"broker-manager/websockets"
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"time"
)

type ClientPublishedEvent struct {
	ID        string `json:"id"`
	TopicName string `json:"topic_name"`
	Payload   string `json:"payload"`
	QoS       uint8  `json:"qos"`
	Retain    bool   `json:"retain"`
	Timestamp uint64 `json:"timestamp"`
}

type OnPublished struct {
	mqtt.HookBase
}

// ID returns the ID of the hook.
func (h *OnPublished) ID() string {
	return "on-published"
}

// Provides indicates which hook methods this hook provides.
func (h *OnPublished) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnPublished,
	}, []byte{b})
}

// OnPublished Intercepts the disconnected client and generates an event to be sent on the websocket
func (h *OnPublished) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	event := ClientPublishedEvent{
		ID:        cl.ID,
		TopicName: pk.TopicName,
		Payload:   string(pk.Payload),
		QoS:       pk.FixedHeader.Qos,
		Retain:    pk.FixedHeader.Retain,
		Timestamp: uint64(time.Now().UnixMilli()),
	}

	h.Log.Info("Client published", "event", event)
	websockets.SendMessage(websockets.MqttClientPublished, event)
}
