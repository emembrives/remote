package main

import (
	"log"
	"net/http"
	"time"

	protobuf "code.google.com/p/goprotobuf/proto"
	"github.com/sterops/tinkerings/android-remote/proto"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"bitbucket.org/gdamore/mangos/protocol/pair"
	transports "bitbucket.org/gdamore/mangos/transport/all"
)

var (
	commands    chan string
	hasFrontend bool
)

func main() {
	hasFrontend = false
	commands = make(chan string)
	go nanomsgListen(commands)
	websocketListen(commands)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func websocketListen(input <-chan string) {
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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("Upgrading connection to websocket")
	}
	for command := range commands {
		log.Print("Sending command ", command)
		conn.WriteJSON(command)
	}
}

func nanomsgListen(output chan<- string) {
	pair, err := pair.NewSocket()
	if err != nil {
		panic(err)
	}
	transports.AddTransports(pair)
	err = pair.Listen("tcp://0.0.0.0:6002")
	if err != nil {
		panic(err)
	}
	for {
		msg, err := pair.Recv()
		if err != nil {
			panic(err)
		}
		cmd := new(proto.Command)
		protobuf.Unmarshal(msg, cmd)
		if *cmd.Type == proto.Command_STATUS {
			response := new(proto.Response)
			response.FrontendConnected = protobuf.Bool(hasFrontend)
			data, err := protobuf.Marshal(response)
			if err != nil {
				panic(err)
			}
			if err = pair.Send(data); err != nil {
				panic(err)
			}
		}
	}
}
