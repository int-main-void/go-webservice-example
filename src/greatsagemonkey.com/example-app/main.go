/*
 *
 *
 *
 */
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"greatsagemonkey.com/example-app/config"
)

const HEARTBEAT_SLEEP_INTERVAL = 30000 * time.Millisecond

const CONFIG_FILENAME_KEY = "CONFIG_FILENAME"
const RUNTIME_STAGE_KEY = "RUNTIME_STAGE"

const LISTENING_PORT_KEY = "ListeningPort"
const SERVER_CERT_FILE_KEY = "ServerCertFile"
const SERVER_KEY_FILE_KEY = "ServerKeyFile"

const CONFIG_NAME = "example-webservice"
const VERSION = "v1"

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, cowboy"))
}

func notFound(resp http.ResponseWriter, req *http.Request) {
	logrus.Println("request for invalid endpoint: ", req.URL.String())
	http.Error(resp, "go fish", http.StatusNotFound)
}

func MyAuthMiddleware(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if req.Header.Get("IMPORTANTSTUFF") == "" {
		http.Error(resp, "you can't do that on television", http.StatusUnauthorized)
		return
	}
	next(resp, req)
}

func setupRoutesAndServe(listeningPort string, serverCert string, serverKey string) error {

	router := mux.NewRouter()
	router.HandleFunc("/example-webservice/v1/hello", hello).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(notFound)

	http.Handle("/", router)

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(MyAuthMiddleware))
	n.UseHandler(router)

	//return http.ListenAndServeTLS(":" + listeningPort, serverCert, serverKey, nil)
	return http.ListenAndServe(":"+listeningPort, n)

}

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
	error = setupRoutesAndServe(listeningPort, serverCert, serverKey)
	if error != nil {
		log.Println("Error setting up server: ", error)
		os.Exit(1)
	}

	// enter heartbeat loop
	for true {
		log.Println("system has been running for ", time.Since(startTime))
		time.Sleep(HEARTBEAT_SLEEP_INTERVAL)
	}
}
