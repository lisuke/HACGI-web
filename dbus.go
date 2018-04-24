package main

import (
	"fmt"
	"github.com/godbus/dbus"
)

func ServiceChanged() {
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
		fmt.Println(typeof(v.Body[0]))
		fmt.Println(typeof(v.Body[1]))
		fmt.Println(typeof(v.Body[2]))
		fmt.Println(typeof(v.Body[3]))
		fmt.Println(typeof(v.Body[4]))
	}
}
func typeof(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
