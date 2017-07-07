package server

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

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
