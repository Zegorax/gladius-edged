package networkd

import (
	"github.com/gladiusio/gladius-networkd/networkd/p2p/handler"
	"github.com/gladiusio/gladius-networkd/networkd/server/contserver"
	"github.com/gladiusio/gladius-networkd/networkd/state"
	log "github.com/sirupsen/logrus"

	"github.com/gladiusio/gladius-utils/config"
	"github.com/gladiusio/gladius-utils/init/manager"
)

// SetupAndRun runs the networkd as a service
func SetupAndRun() {
	// Define some variables
	name, displayName, description :=
		"GladiusNetworkDaemon",
		"Gladius Network (Edge) Daemon",
		"Gladius Network (Edge) Daemon"

	// Run the function "run" in newtworkd as a service
	manager.RunService(name, displayName, description, Run)
}

// Run - Start a web server
func Run() {
	log.Info("Loading config")

	// Setup config handling
	config.SetupConfig("gladius-networkd", config.NetworkDaemonDefaults())

	// Setup logging level
	switch loglevel := config.GetString("LogLevel"); loglevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	log.Info("Starting...")

	// Create a p2p handler
	controldBase := config.GetString("ControldProtocol") + "://" + config.GetString("ControldHostname") + ":" + config.GetString("ControldPort") + "/api/p2p"
	// TODO: Get seed node from the blockchain
	p2pHandler := handler.New(controldBase, config.GetString("P2PSeedNodeAddress"))
	go p2pHandler.Connect()

	// Create new thread safe state of the networkd
	s := state.New(p2pHandler)

	// Create a content server
	cs := contserver.New(s)
	cs.Start()
	defer cs.Stop()

	log.Info("Started HTTP server.")

	// Block forever
	select {}
}
