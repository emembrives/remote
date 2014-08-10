package main

import (
	"log"
	"net/http"
	"time"

	protobuf "code.google.com/p/goprotobuf/proto"
	"github.com/sterops/tinkerings/android-remote/proto"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/streadway/amqp"
)

func main() {
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
	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Println(err)
		return
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		log.Println(err)
		return
	}

	defer ch.Close()

	q, err := ch.QueueDeclare(
		"websocket", // name
		false,    // durable
		false,    // delete when usused
		true,     // exclusive
		false,    // noWait
		nil,      // arguments
	)
	if err != nil {
		log.Println(err)
		return
	}

	msgs, err := ch.Consume(q.Name, "", true,
		true,
		false,
		false,
		nil)
	if err != nil {
		log.Println(err)
		return
	}

	for d := range msgs {
		log.Printf("Received a message: %s", d.Body)
		cmd := new(proto.Command)
		protobuf.Unmarshal(d.Body, cmd)
		if *cmd.Type == proto.Command_COMMAND {
			websocketConn.WriteJSON(*cmd.Command)
		}
	}
}
