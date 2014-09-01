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
	//r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/", websocketHandler)
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

	if err = ch.ExchangeDeclare(
		"remote", // name
		"topic",  // type
		true,     // durable
		true,     // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	); err != nil {
		log.Println(err)
		return
	}

	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		true,  // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Println(err)
		return
	}

	if err = ch.QueueBind(
		q.Name,
		"webcontrol",
		"remote",
		false,
		nil); err != nil {
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
		case d := <-msgs:
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
