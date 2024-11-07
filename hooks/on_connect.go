package hooks

import (
	"broker-manager/websockets"
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"time"
)

type ClientConnectedEvent struct {
	ID              string `json:"id"`
	ProtocolVersion uint8  `json:"protocol_version"`
	Username        string `json:"username"`
	Remote          string `json:"remote"`

	QoS       uint8  `json:"qos"`
	KeepAlive uint16 `json:"keep_alive"`
	Timestamp uint64 `json:"timestamp"`
}

// OnConnect intercepts new connections
type OnConnect struct {
	mqtt.HookBase
}

// ID returns the ID of the hook.
func (h *OnConnect) ID() string {
	return "on-connect"
}

// Provides indicates which hook methods this hook provides.
func (h *OnConnect) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
	}, []byte{b})
}

// OnConnect Intercepts the new connection and generates an event to be sent on the websocket
func (h *OnConnect) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	event := ClientConnectedEvent{
		ID:              cl.ID,
		ProtocolVersion: pk.ProtocolVersion,
		Username:        string(pk.Connect.Username),
		Remote:          cl.Net.Remote,
		QoS:             pk.Properties.MaximumQos,
		KeepAlive:       pk.Connect.Keepalive,
		Timestamp:       uint64(time.Now().UnixMilli()),
	}

	h.Log.Info("New connection", "event", event)
	websockets.SendMessage(websockets.MqttClientConnected, event)
	return nil
}
