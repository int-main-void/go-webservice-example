package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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

func helloPostHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, cowboy poster"))
}

func widgetsHandler(resp http.ResponseWriter, req *http.Request) {
	dbclient, err := db.GetAppDb()
	if err != nil {
		logrus.Error(err)
		http.Error(resp, "", http.StatusInternalServerError)
		return
	}
	widgets, err := dbclient.GetWidgets()
	if err != nil {
		logrus.Error(err)
		http.Error(resp, "", http.StatusInternalServerError)
		return
	}
	widgetsBytes, err := json.Marshal(widgets)
	if err != nil {
		logrus.Error(err)
		http.Error(resp, "", http.StatusInternalServerError)
		return
	}
	resp.Write(widgetsBytes)
}

func notFoundHandler(resp http.ResponseWriter, req *http.Request) {
	logrus.Println("request for invalid endpoint: ", req.URL.String())
	http.Error(resp, "go fish", http.StatusNotFound)
}
