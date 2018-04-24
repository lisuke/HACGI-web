package main

import (
	// "encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
)

type Hub struct {
	clients    map[*Client]bool
	message    chan []byte
	register   chan *Client
	unregister chan *Client
}

type InvokeServiceMethod struct {
	ServiceName   string `json:"service_name"`
	InterfaceName string `json:"interface_name"`
	kwargs        string `json:"kwargs"`
}

type HubMessage struct {
	ClientId string `json:"ClientId"`
	Message  string `json:"Message"`
}

type TransferMessage struct {
	ReqType string `json:"reqType,string"`
	Data    string `json:"data,string"`
}

func newHub() *Hub {
	return &Hub{
		message:    make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
			}
		case message := <-h.message:

			js, _ := simplejson.NewJson(message)
			reqType, _ := js.GetPath("Message", "reqType").String()

			switch reqType {
			case "invoke-remote":
				invoke(h, js)
			case "query":
				query(h, js)
			case "response":
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

func response(h *Hub, data string) {
	tmp := []byte(`{"Message":{"reqType":"response","data":` + data + `}}`)
	h.message <- tmp
}

func invoke(h *Hub, js *simplejson.Json) {
	reqType, _ := js.GetPath("Message", "reqType").String()
	data, _ := js.GetPath("Message", "data").String()
	fmt.Println(reqType, data)
	go response(h, `"hello world"`)
}

func query(h *Hub, js *simplejson.Json) {

}
