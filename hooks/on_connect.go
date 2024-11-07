package hooks

import (
	"broker-manager/websockets"
	"bytes"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

type ClientConnectedEvent struct {
	ID              string `json:"id"`
	ProtocolVersion uint8  `json:"protocol_version"`
	Username        string `json:"username"`
	Remote          string `json:"remote"`

	QoS       uint8  `json:"qos"`
	KeepAlive uint16 `json:"keep_alive"`
}

// Options contains the configuration.
type Options struct {
	Reverb *websockets.ReverbConn
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

func (h *OnConnect) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	event := ClientConnectedEvent{
		ID:              cl.ID,
		ProtocolVersion: pk.ProtocolVersion,
		Username:        string(pk.Connect.Username),
		Remote:          cl.Net.Remote,
		QoS:             pk.Properties.MaximumQos,
		KeepAlive:       pk.Connect.Keepalive,
	}

	h.Log.Info("New connection", "event", event)
	websockets.SendMessage(websockets.MqttClientConnected, event)
	return nil
}
