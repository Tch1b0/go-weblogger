package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/Tch1b0/go-weblogger/pkg/chans"
)

//go:embed index.html
var rootHTML string

type WebInterfaceWriter struct {
	msgChan *chans.SuperChannel[[]byte]
}

func NewWebInterfaceWriter(msgChan *chans.SuperChannel[[]byte]) WebInterfaceWriter {
	return WebInterfaceWriter{
		msgChan: msgChan,
	}
}

func (w WebInterfaceWriter) Write(p []byte) (n int, err error) {
	w.msgChan.Sender <- p

	return len(p), nil
}

type WebInterface struct {
	Host    string
	Port    int
	msgChan *chans.SuperChannel[[]byte]
	Writer  WebInterfaceWriter
}

func NewWebInterface(host string, port int) WebInterface {
	sc := chans.NewSuperChannel[[]byte]()
	go sc.Start()
	return WebInterface{
		Host:    host,
		Port:    port,
		msgChan: sc,
		Writer:  WebInterfaceWriter{msgChan: sc},
	}
}

func (wi *WebInterface) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, rootHTML)
	})

	// Add websocket functionality
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024 * 100,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}

	mux.HandleFunc("/io/stream", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error upgrading to websocket:", err)
			return
		}
		defer conn.Close()

		recv := chans.NewDebouncedChan[[]byte]()
		go recv.Debounce()
		wi.msgChan.AddReceiver(recv.Sender)

		// Keep the connection alive without transferring data
		for {
			messages, ok := <-recv.Receiver

			if ok {
				for i := range messages {
					msg := messages[len(messages)-i-1]
					if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
						fmt.Println("Couldnt send message")
						// connection was most likely closed
						wi.msgChan.RemoveReceiver(recv.Sender)
						return
					}
				}
			} else {
				wi.msgChan.RemoveReceiver(recv.Sender)
				return
			}
		}

	})

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", wi.Host, wi.Port),
		Handler: mux,
	}

	log.Printf("Starting web server on http://%s:%d\n", wi.Host, wi.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server error:", err)
	}
}
