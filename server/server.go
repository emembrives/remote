package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sterops/tinkerings/android-remote/proto"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"bitbucket.org/gdamore/mangos"
	"bitbucket.org/gdamore/mangos/protocol/pair"
)

var (
	commands chan string
)

func main() {
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
	pair.Listen("tcp://0.0.0.0:6002")
}
