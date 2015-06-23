package services

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/emembrives/tinkerings/android-remote/proto"
	protobuf "github.com/golang/protobuf/proto"
)

type PresentationServiceFactory struct {
	upgrader websocket.Upgrader
	conn     *websocket.Conn

	services map[string]*PresentationService
}

func NewPresentationServiceFactory() *PresentationServiceFactory {
	factory := &PresentationServiceFactory{}
	factory.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return factory
}

func (f *PresentationServiceFactory) websocketListen() {
	r := mux.NewRouter()
	h := &http.Server{
		Addr:           ":6001",
		Handler:        r,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.HandleFunc("/conn", func(w http.ResponseWriter, r *http.Request) { f.websocketHandler(w, r) })
	go h.ListenAndServe()
}

func (f *PresentationServiceFactory) websocketHandler(w http.ResponseWriter, r *http.Request) {
	websocketConn, err := f.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	serviceOrigin := r.URL.Host + "/" + r.URL.Path
	if _, ok := f.services[serviceOrigin]; ok {
		log.Println(serviceOrigin, "is already connected")
		websocketConn.Close()
		return
	}

	f.services[serviceOrigin] = NewPresentationService(serviceOrigin, websocketConn)
	go f.services[serviceOrigin].Run()
}

type PresentationService struct {
	origin   string
	conn     *websocket.Conn
	Messages chan string
}

func NewPresentationService(origin string, websocketConn *websocket.Conn) *PresentationService {
	ps := new(PresentationService)
	ps.origin = origin
	ps.conn = websocketConn
	ps.Messages = make(chan string)
	return ps
}

func (ps *PresentationService) Run() {
	go func(conn *websocket.Conn) {
		for {
			var message interface{}
			err := conn.ReadJSON(&message)
			if err != nil {
				log.Println("Reading error on", ps.origin)
				conn.Close()
				return
			}
		}
	}(ps.conn)

	for {
		select {
		case d := <-ps.Messages:
			log.Printf("Received a message: %s", d)
			ps.conn.WriteJSON(d)
		}
	}
}

func (s *PresentationService) GetDescription() *proto.ServiceDefinition {
	p := &proto.ServiceDefinition{}
	p.ServiceName = protobuf.String("presentation")
	p.Method = append(p.Method,
		&proto.MethodDefinition{
			Name: protobuf.String("Next"),
		})
	p.Method = append(p.Method,
		&proto.MethodDefinition{
			Name: protobuf.String("Previous"),
		})
	return p
}

func (s *PresentationService) MakeCall(call *proto.RPCRequest) ([]*proto.ObjectType, error) {

	return nil, errors.New("")
}
