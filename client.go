package main

import (
	"fmt"

	"github.com/panjf2000/gnet/v2"
)

type CHandler struct {
	gnet.BuiltinEventEngine
}

func (ev *CHandler) OnOpen(_ gnet.Conn) ([]byte, gnet.Action) {
	Log.Debugln("Authelia connection opened")
	return nil, gnet.None
}

func (ev *CHandler) OnClose(_ gnet.Conn, _ error) gnet.Action {
	Log.Debugln("Authelia connection closed")
	return gnet.Close
}

func (ev *CHandler) OnTraffic(c gnet.Conn) gnet.Action {
	data, err := c.Next(-1)
	if err != nil {
		Log.Errorf("error while reading response: %v", err)
	}

	ctx := c.Context().(SubReq)

	copy(ctx.data, data)
	ctx.notif <- len(data)

	return gnet.None
}

func PingConnection(c gnet.Conn) error {
	_, err := c.Write([]byte("GET /api/authz/auth-request HTTP/1.1\r\nX-Forwarded-For: 0.0.0.0\r\nX-Original-URL: " + Conf.External + "\r\nX-Original-Method: GET\r\n\r\n"))
	if err != nil {
		return fmt.Errorf("error while writing to connection: %w", err)
	}

	return nil
}

func CreateConnections(C *gnet.Client) error {
	for i := 0; i < Conf.Authelia.Connections; i++ {
		data, notif := make([]byte, 2048), make(chan int)

		c, err := C.DialContext("tcp4", Conf.Authelia.Address, SubReq{
			data:  data,
			notif: notif,
		})
		if err != nil {
			return fmt.Errorf("error while connecting to Authelia: %w", err)
		}

		if err = PingConnection(c); err != nil {
			return err
		}

		<-notif
		Conns <- c
	}

	return nil
}
