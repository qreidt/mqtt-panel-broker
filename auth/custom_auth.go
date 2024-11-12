// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2022 mochi-mqtt, mochi-co
// SPDX-FileContributor: mochi-co

package auth

import (
	"broker-manager/services"
	"bytes"

	"github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
)

// CustomAuth validates credentials with external services
type CustomAuth struct {
	mqtt.HookBase
}

// ID returns the ID of the hook.
func (h *CustomAuth) ID() string {
	return "custom-auth"
}

// Provides indicates which hook methods this hook provides.
func (h *CustomAuth) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate,
		mqtt.OnACLCheck,
	}, []byte{b})
}

// OnConnectAuthenticate returns true/allowed for all requests.
func (h *CustomAuth) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	h.Log.Info("Authenticating",
		"username", string(pk.Connect.Username),
		"remote", cl.Net.Remote)

	return services.AuthServiceInstance.Authenticate(
		cl.ID, string(pk.Connect.Username), string(pk.Connect.Password))
}

// OnACLCheck returns true/allowed for all checks.
func (h *CustomAuth) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	h.Log.Info("Authenticating ACL",
		"client", cl.ID,
		"username", string(cl.Properties.Username),
		"topic", topic)

	return true
}
