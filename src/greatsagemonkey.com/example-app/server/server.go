package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"greatsagemonkey.com/example-app/db"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func sockittomeHandler(resp http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(resp, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if err = conn.WriteMessage(messageType, p); err != nil {
			logrus.Error(err)
			return
		}
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, cowboy"))
}

func widgetsHandler(resp http.ResponseWriter, req *http.Request) {
	widgets, _ := db.GetAppDb().GetWidgets()
	widgetsBytes, _ := json.Marshal(widgets)
	resp.Write(widgetsBytes)
}

func notFoundHandler(resp http.ResponseWriter, req *http.Request) {
	logrus.Println("request for invalid endpoint: ", req.URL.String())
	http.Error(resp, "go fish", http.StatusNotFound)
}

func authMiddleware(resp http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if req.Header.Get("IMPORTANTSTUFF") == "" {
		http.Error(resp, "you can't do that on television", http.StatusUnauthorized)
		return
	}
	next(resp, req)
}

func setupRoutesAndServe(listeningPort string, serverCert string, serverKey string) error {

	router := mux.NewRouter().PathPrefix("/example-webservice/v1").Subrouter()
	router.HandleFunc("/widgets", widgetsHandler).Methods("GET")
	router.HandleFunc("/hello", helloHandler).Methods("GET")
	router.HandleFunc("/sockittome", sockittomeHandler)
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	http.Handle("/", router)

	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(authMiddleware))
	n.UseHandler(router)

	//return http.ListenAndServeTLS(":" + listeningPort, serverCert, serverKey, nil)
	return http.ListenAndServe(":"+listeningPort, n)

}

func StartServer(listeningPort string, certFile string, certKey string) {
	error := setupRoutesAndServe(listeningPort, certFile, certKey)
	if error != nil {
		log.Println("Error setting up server: ", error)
		os.Exit(1)
	}
}
