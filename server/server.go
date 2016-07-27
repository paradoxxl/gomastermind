package server

import (
	"golang.org/x/net/websocket"
	"net/http"
)

type Server struct{}

func NewServer(){
	http.Handle("/echo", websocket.Handler(ConnectionHandler))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}


}

func ConnectionHandler(ws *websocket.Conn) {
	defer ws.Close()
	client := ClientHandler{ws:ws}
	client.Listen()
}
