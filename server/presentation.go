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

type websocketHandler struct {
	server *ZMQServer
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
	handler := &websocketHandler{server}
	r.Handle("/conn", handler)
	h.ListenAndServe()
}

// Serves the websocket connection
func (h *websocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	websocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	service := &websocketService{
		conn:     websocketConn,
		url:      r.Referer(),
		server:   h.server,
		messages: make(chan string),
	}

	h.server.Services[service.url] = service
	go service.Run()
}

type websocketService struct {
	conn     *websocket.Conn
	url      string
	server   *ZMQServer
	messages chan string
}

func (s *websocketService) Name() string {
	return s.url
}

func (s *websocketService) Endpoints() []string {
	return []string{"next", "prev"}
}

func (s *websocketService) ReadEndpoint(name string) string {
	return ""
}

func (s *websocketService) WriteEndpoint(name, value string) string {
	if value != "" {
		s.messages <- name
	}
	return value
}

func (s *websocketService) Run() {
	defer s.conn.Close()
	defer s.Quit()
	log.Println("Upgrading connection to websocket")

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
	}(s.conn, websocketErrors)

	for {
		select {
		case d := <-s.messages:
			log.Printf("Received a message: %s", d)
			s.conn.WriteJSON(d)
		case err := <-websocketErrors:
			if err != nil {
				log.Println("Closing websocket connection")
				return
			}
		}
	}
}

func (s *websocketService) Quit() {
	delete(s.server.Services, s.url)
}
