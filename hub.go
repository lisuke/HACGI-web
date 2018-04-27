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
			fmt.Println(reqType)
			switch reqType {
			case "invoke-remote":
				go serviceInvoke(h, js)
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

func responseToAll(h *Hub, data string) {
	tmp_str := `{"Message":{"reqType":"response","data":` + data + `}}`
	fmt.Println("send to: all ", tmp_str)
	tmp_byte := []byte(tmp_str)
	h.message <- tmp_byte
}

func responseToClientId(h *Hub, data string, clientId string) {
	tmp_str := `{"Message":{"reqType":"response","data":` + data + `}}`
	fmt.Println("send to: ", clientId, tmp_str)
	tmp_byte := []byte(tmp_str)
	client := h.getClient(clientId)
	client.send <- tmp_byte
}

func query(h *Hub, js *simplejson.Json) {
	// data, _ := js.GetPath("Message", "data").String()
	// fmt.Println(reqType, data)
	resource, _ := js.GetPath("Message", "data", "resource").String()
	clientId, _ := js.GetPath("ClientId").String()
	switch resource {
	case "getAllServices":
		go getAllServices(h, clientId)
	case "getAllStatus":
		go getAllStatus(h, js)
	}
	// go response(h, `"hello world"`)
}

func (h *Hub) getClient(ClientId string) *Client {
	for client := range h.clients {
		if client.cid == ClientId {
			return client
		}
	}
	return nil
}
