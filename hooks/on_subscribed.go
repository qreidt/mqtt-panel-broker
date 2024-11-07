package hooks

import (
	"broker-manager/websockets"
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"time"
)

type ClientSubscribedEvent struct {
	ID        string `json:"id"`
	TopicName string `json:"topic_name"`
	QoS       uint8  `json:"qos"`
	Timestamp uint64 `json:"timestamp"`
}

type OnSubscribed struct {
	mqtt.HookBase
}

// ID returns the ID of the hook.
func (h *OnSubscribed) ID() string {
	return "on-connect"
}

// Provides indicates which hook methods this hook provides.
func (h *OnSubscribed) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnSubscribed,
	}, []byte{b})
}

// OnSubscribed Intercepts the disconnected client and generates an event to be sent on the websocket
func (h *OnSubscribed) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	if len(pk.Filters) == 0 {
		return
	}

	filter := pk.Filters[0]
	event := ClientSubscribedEvent{
		ID:        cl.ID,
		TopicName: filter.Filter,
		QoS:       filter.Qos,
		Timestamp: uint64(time.Now().UnixMilli()),
	}

	h.Log.Info("Client subscribed to a topic", "event", event)
	websockets.SendMessage(websockets.MqttClientSubscribed, event)
}
