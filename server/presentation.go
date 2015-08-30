package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func websocketListen(server *ZMQServer) {
	r := mux.NewRouter()
	h := &http.Server{
		Addr:           ":6001",
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/conn", websocketHandler)
	h.ListenAndServe()
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	websocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	} else {
		defer websocketConn.Close()
		log.Println("Upgrading connection to websocket")
	}

	websocketErrors := make(chan error)
	go func(conn *websocket.Conn, errChan chan error) {
		for {
			var message interface{}
			err := conn.ReadJSON(&message)
			if err != nil {
				errChan <- err
				conn.Close()
				return
			}
		}
	}(websocketConn, websocketErrors)

	for {
		select {
		case d := <-messages:
			log.Printf("Received a message: %s", d)
			websocketConn.WriteJSON(d)
		case err = <-websocketErrors:
			if err != nil {
				return
			}
		}
	}
	log.Println("Closing websocket and amqp connections")
}
