/*
 *
 *
 *
 */
package main

import (
	"log"
	"os"
	"time"

	"greatsagemonkey.com/example-app/config"
	"greatsagemonkey.com/example-app/server"
)

const HEARTBEAT_SLEEP_INTERVAL = 30000 * time.Millisecond

const CONFIG_FILENAME_KEY = "CONFIG_FILENAME"
const RUNTIME_STAGE_KEY = "RUNTIME_STAGE"

const LISTENING_PORT_KEY = "ListeningPort"
const SERVER_CERT_FILE_KEY = "ServerCertFile"
const SERVER_KEY_FILE_KEY = "ServerKeyFile"

const CONFIG_NAME = "example-webservice"
const VERSION = "v1"

func main() {
	startTime := time.Now()

	// set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	// set up configuration

	configFilename := os.Getenv(CONFIG_FILENAME_KEY)
	runtimeStage := os.Getenv(RUNTIME_STAGE_KEY)

	conf, error := config.NewConfig(configFilename, CONFIG_NAME, VERSION, runtimeStage)
	if error != nil {
		log.Println("FATAL: error setting up configuration: ", error)
		os.Exit(1)
	}
	log.Println(conf)

	// run main program
	listeningPort := conf[LISTENING_PORT_KEY]
	serverCert := conf[SERVER_CERT_FILE_KEY]
	serverKey := conf[SERVER_KEY_FILE_KEY]
	server.StartServer(listeningPort, serverCert, serverKey)

	// enter heartbeat loop
	for true {
		log.Println("system has been running for ", time.Since(startTime))
		time.Sleep(HEARTBEAT_SLEEP_INTERVAL)
	}
}
