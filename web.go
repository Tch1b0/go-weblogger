package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

//go:embed index.html
var rootHTML string

var msgChansMutex *sync.Mutex = new(sync.Mutex)

func RemoveFromArr[T comparable](arr *[]T, value T) {
	passedObj := false
	for i := 0; i < len(*arr)-1; i++ {
		if (*arr)[i] == value {
			passedObj = true
		}

		if passedObj {
			(*arr)[i] = (*arr)[i+1]
		}
	}

	*arr = (*arr)[:len(*arr)-1]
}

type WebInterfaceWriter struct {
	msgChans *[](chan []byte)
}

func NewWebInterfaceWriter(msgChans *[](chan []byte)) WebInterfaceWriter {
	return WebInterfaceWriter{
		msgChans: msgChans,
	}
}

func (w WebInterfaceWriter) Write(p []byte) (n int, err error) {
	msgChansMutex.Lock()

	data := append([]byte(nil), p...)

	var queuedRemoval [](chan []byte)

	for _, c := range *w.msgChans {
		select {
		case c <- data:
			// successfully wrote to c
		default:
			queuedRemoval = append(queuedRemoval, c)
		}
	}

	for _, c := range queuedRemoval {
		RemoveFromArr(w.msgChans, c)
	}
	msgChansMutex.Unlock()

	return len(p), nil
}

type WebInterface struct {
	Host     string
	Port     int
	msgChans *[](chan []byte)
	Writer   WebInterfaceWriter
}

func NewWebInterface(host string, port int) WebInterface {
	var msgChans [](chan []byte)

	return WebInterface{
		Host:     host,
		Port:     port,
		msgChans: &msgChans,
		Writer:   WebInterfaceWriter{msgChans: &msgChans},
	}
}

func (wi *WebInterface) Serve() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, rootHTML)
	}))

	http.Handle("/io/stream", websocket.Handler(func(ws *websocket.Conn) {
		recv := make(chan []byte, 100)
		fmt.Println("Adding to Receivers: ", recv)
		msgChansMutex.Lock()
		*wi.msgChans = append((*wi.msgChans), recv)
		msgChansMutex.Unlock()

		for {
			message, ok := <-recv

			if ok {
				if err := websocket.Message.Send(ws, string(message)); err != nil {
					// connection was most likely closed
					msgChansMutex.Lock()
					RemoveFromArr(wi.msgChans, recv)
					msgChansMutex.Unlock()
					return
				}
			} else {
				msgChansMutex.Lock()
				RemoveFromArr(wi.msgChans, recv)
				msgChansMutex.Unlock()
				return
			}
		}

	}))

	log.Printf("Starting web server on http://%s:%d\n", wi.Host, wi.Port)

	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", wi.Port), nil)
}
