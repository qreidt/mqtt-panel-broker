package hooks

import (
	"broker-manager/websockets"
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
	"time"
)

type ClientDisconnectedEvent struct {
	ID        string `json:"id"`
	Timestamp uint64 `json:"timestamp"`
}

type OnDisconnect struct {
	mqtt.HookBase
}

// ID returns the ID of the hook.
func (h *OnDisconnect) ID() string {
	return "on-disconnect"
}

// Provides indicates which hook methods this hook provides.
func (h *OnDisconnect) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnDisconnect,
	}, []byte{b})
}

// OnDisconnect Intercepts the disconnected client and generates an event to be sent on the websocket
func (h *OnDisconnect) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	event := ClientDisconnectedEvent{
		ID:        cl.ID,
		Timestamp: uint64(time.Now().UnixMilli()),
	}

	h.Log.Info("Client Disconnected", "event", event)
	websockets.SendMessage(websockets.MqttClientDisconnected, event)
}
