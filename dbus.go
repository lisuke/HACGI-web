package main

import (
	"fmt"
	"github.com/godbus/dbus"
)

func dbus_init(hub *Hub) {
	go ServiceChanged(hub)
	// go getAllServices(hub, "")
}

func ServiceChanged(hub *Hub) {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		"type='signal',interface='com.HACGI.convergence.ServiceManager',member='ServicesChanged'",
	)

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)
	for v := range c {
		json_str := fmt.Sprintf(`{"resource":"ServiceChanged","Message": {"ServiceType": "%s","ServiceBus": "%s","ObjectPath": "%s","ChangeType": "%s","has_new": "%t"}`, v.Body[0], v.Body[1], v.Body[2], v.Body[3], v.Body[4])
		go responseToAll(hub, json_str)
	}
}

func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

func getAllServices(hub *Hub, clientId string) {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}
	v, _ := conn.Object("com.HACGI.convergence", "/com/HACGI/convergence/ServiceManager").GetProperty("com.HACGI.convergence.ServiceManager.ServicesJSON")
	fmt.Println(typeof(v.Value()))
	if typeof(v.Value()) == "string" {
		json_str := fmt.Sprintf(`{"resource":"getAllServices","Message": %s}`, v.Value())
		go responseToClientId(hub, json_str, clientId)
	}
}
