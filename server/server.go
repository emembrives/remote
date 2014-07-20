package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	commands chan string
)

func main() {
	commands = make(chan string)
	go socketListen(commands)
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

func socketListen(output chan<- string) {
	server, err := net.Listen("tcp", ":6002")
	if err != nil {
		panic(err)
	}
	defer server.Close()

	for {
		// Listen for an incoming connection.
		conn, err := server.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
			continue
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn, output)
	}
}

func handleRequest(conn net.Conn, output chan<- string) {
	log.Print("New connection")
	reader := bufio.NewReader(conn)
	for {
		message, err := readFromClient(reader)
		if err == io.EOF {
			log.Print("Connection interrupted")
			return
		} else if err != nil {
			log.Print("Error while reading from socket: ", err)
			return
		}
		output <- message
	}
}

func readFromClient(reader *bufio.Reader) (string, error) {
	data, err := reader.ReadSlice('\n')
	message := strings.Trim(string(data), " \n")
	return message, err
}
