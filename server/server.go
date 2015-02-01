package main

import (
	"log"
	"net/http"
	"time"

	protobuf "code.google.com/p/goprotobuf/proto"
	"github.com/sterops/tinkerings/android-remote/proto"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/gdamore/mangos"
	"github.com/gdamore/mangos/protocol/pull"
	"github.com/gdamore/mangos/transport/tcp"

)

const (
	nanomsgUrl = "tcp://0.0.0.0:7001"
)

var (
	messages = make(chan *mangos.Message)
)

func main() {
	go setupMessageServer()
	websocketListen()
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func websocketListen() {
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

func setupMessageServer() {
	var sock mangos.Socket
	var err error
	if sock, err = pull.NewSocket(); err != nil {
		log.Println("can't get new pull socket: %s", err)
		return
	}

	sock.AddTransport(tcp.NewTransport())
	if err = sock.Listen(nanomsgUrl); err != nil {
		log.Println("can't listen on pull socket: %s", err.Error())
		return
	}

	for {
		msg, err := sock.RecvMsg()
		log.Println("Received message: %d", msg)
		if err != nil {
			log.Println(err)
		} else {
			messages <- msg
		}
	}
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
			log.Printf("Received a message: %s", d.Body)
			cmd := new(proto.Command)
			protobuf.Unmarshal(d.Body, cmd)
			if *cmd.Type == proto.Command_COMMAND {
				websocketConn.WriteJSON(*cmd.Command)
			}
		case err = <-websocketErrors:
			if err != nil {
				return
			}
		}
	}
	log.Println("Closing websocket and amqp connections")
}
