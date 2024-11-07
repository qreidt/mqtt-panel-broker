package hooks

import (
	"broker-manager/websockets"
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"time"
)

type ClientUnsubscribedEvent struct {
	ID        string `json:"id"`
	TopicName string `json:"topic_name"`
	Timestamp uint64 `json:"timestamp"`
}

type OnUnsubscribed struct {
	mqtt.HookBase
}

// ID returns the ID of the hook.
func (h *OnUnsubscribed) ID() string {
	return "on-unsubscribed"
}

// Provides indicates which hook methods this hook provides.
func (h *OnUnsubscribed) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnUnsubscribed,
	}, []byte{b})
}

// OnUnsubscribed Intercepts the disconnected client and generates an event to be sent on the websocket
func (h *OnUnsubscribed) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	if len(pk.Filters) == 0 {
		return
	}

	filter := pk.Filters[0]
	event := ClientUnsubscribedEvent{
		ID:        cl.ID,
		TopicName: filter.Filter,
		Timestamp: uint64(time.Now().UnixMilli()),
	}

	h.Log.Info("Client unsubscribed from a topic", "event", event)
	websockets.SendMessage(websockets.MqttClientUnsubscribed, event)
}
