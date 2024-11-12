package main

import (
	"broker-manager/auth"
	"broker-manager/hooks"
	"broker-manager/services"
	"broker-manager/websockets"
	"flag"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var server *mqtt.Server

func main() {
	websockets.Init()
	services.AuthServiceInit()

	//goland:noinspection GoUnhandledErrorResult
	defer websockets.Close()

	// Create signals channel to run server until interrupted
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// Create the new MQTT Server.
	server = mqtt.New(nil)

	setupHooks()
	setupListeners()

	// Start Server
	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Run server until interrupted
	<-done

	// Cleanup
	server.Log.Warn("caught signal, stopping...")
	_ = server.Close()
	server.Log.Info("mochi mqtt shutdown complete")
}

func setupHooks() {
	// Allow all connections. ToDo setup authentication
	_ = server.AddHook(new(auth.CustomAuth), nil)

	// Setup intercept hooks
	_ = server.AddHook(new(hooks.OnConnect), nil)
	_ = server.AddHook(new(hooks.OnDisconnect), nil)
	_ = server.AddHook(new(hooks.OnSubscribed), nil)
	_ = server.AddHook(new(hooks.OnUnsubscribed), nil)
	_ = server.AddHook(new(hooks.OnPublished), nil)

	_ = server.AddHook(new(hooks.OnPacketProcessed), nil)
}

func setupListeners() {
	tcpAddr := flag.String("tcp", ":1883", "network address for TCP listener")
	wsAddr := flag.String("ws", ":1882", "network address for Websocket listener")
	infoAddr := flag.String("info", ":8080", "network address for web info dashboard listener")
	flag.Parse()

	// Create a TCP listener on a standard port.
	tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: *tcpAddr})

	// Create WebSocket Listener
	ws := listeners.NewWebsocket(listeners.Config{
		ID:      "ws1",
		Address: *wsAddr,
	})

	// Create HTTP status port
	stats := listeners.NewHTTPStats(
		listeners.Config{
			ID:      "info",
			Address: *infoAddr,
		},
		server.Info,
	)

	// Listen to tcp
	if err := server.AddListener(tcp); err != nil {
		log.Fatal(err)
	}

	// Listen to ws
	if err := server.AddListener(ws); err != nil {
		log.Fatal(err)
	}

	// Listen to HTTP
	if err := server.AddListener(stats); err != nil {
		log.Fatal(err)
	}
}
