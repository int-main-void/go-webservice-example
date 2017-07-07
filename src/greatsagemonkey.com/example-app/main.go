/*
 *
 *
 *
 */
package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"greatsagemonkey.com/example-app/server"
)

const HEARTBEAT_SLEEP_INTERVAL = 30000 * time.Millisecond

const CONFIG_FILENAME_KEY = "CONFIG_FILENAME"
const RUNTIME_STAGE_KEY = "RUNTIME_STAGE"

const LISTENING_PORT_KEY = "LISTENING_PORT"
const SERVER_CERT_FILE_KEY = "SERVER_CERT_FILE"
const SERVER_KEY_FILE_KEY = "SERVER_KEY_FILE"

const CONFIG_NAME = "example-webservice"
const VERSION = "v1"

func main() {
	startTime := time.Now()

	// set up logging
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)

	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	if os.Getenv("LOG_LEVEL") == "DEBUG" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// run main program
	listeningPort := os.Getenv(LISTENING_PORT_KEY)
	_, err := strconv.Atoi(listeningPort)
	if err != nil {
		logrus.Error("invalid listening port: ", listeningPort)
		os.Exit(1)
	}
	serverCert := os.Getenv(SERVER_CERT_FILE_KEY)
	serverKey := os.Getenv(SERVER_KEY_FILE_KEY)
	go server.StartServer(listeningPort, serverCert, serverKey)

	// enter heartbeat loop
	for true {
		logrus.Info("system has been running for ", time.Since(startTime))
		time.Sleep(HEARTBEAT_SLEEP_INTERVAL)
	}
}
