package main

import (
	"encoding/json"
	"fmt"
	"github.com/bitly/go-simplejson"
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
		json_str := fmt.Sprintf(`{"resource":"ServiceChanged","Message": {"serviceType": "%s","serviceName": "%s","objectPath": "%s","changeType": "%s","has_new": "%t"}}`, v.Body[0], v.Body[1], v.Body[2], v.Body[3], v.Body[4])
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

func serviceInvoke(hub *Hub, js *simplejson.Json) {
	// clientId, _ := js.GetPath("ClientId").String()
	resource, _ := js.GetPath("Message", "data", "resource").String()
	objectPath, _ := js.GetPath("Message", "data", "objectPath").String()
	ifaceName, _ := js.GetPath("Message", "data", "ifaceName").String()
	method, _ := js.GetPath("Message", "data", "method").String()
	argsJsonStr, _ := js.GetPath("Message", "data", "args").String()

	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	obj := conn.Object(resource, dbus.ObjectPath(objectPath))

	var ret string

	obj.Call(resource+".CallControlInterface", 0, ifaceName, method, argsJsonStr).Store(&ret)

	json_str := fmt.Sprintf(`{"resource":"serviceInvoke","Message": {"resource":"%s","ifaceName":"%s","method":"%s","ret":%s}}`, resource, ifaceName, method, ret)
	// go responseToClientId(hub, json_str, clientId)
	go responseToAll(hub, json_str)
}

func getAllStatus(hub *Hub, js *simplejson.Json) {
	clientId, _ := js.GetPath("ClientId").String()
	Service, _ := js.GetPath("Message", "data", "Service").String()
	objectPath, _ := js.GetPath("Message", "data", "objectPath").String()

	conn, err := dbus.SessionBus()
	if err != nil {
		panic(err)
	}

	obj := conn.Object(Service, dbus.ObjectPath(objectPath))

	var ifaces []string
	obj.Call(Service+".getAllInterfaces", 0).Store(&ifaces)
	iface_json, _ := json.Marshal(&ifaces)

	var status string
	obj.Call(Service+".getStatus", 0).Store(&status)

	json_str := fmt.Sprintf(`{"resource":"getAllStatus","Message": {"interfaces":%s,"status":%s}}`, iface_json, status)
	fmt.Println("json_str:", json_str)
	go responseToClientId(hub, json_str, clientId)
}
