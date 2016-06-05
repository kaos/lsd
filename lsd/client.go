package main

import (
	log "github.com/Sirupsen/logrus"
	"net"
	"runtime/debug"
)

type Client interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close()
	Log() *log.Entry
	HandleErrors()
}

type SocketClient struct {
	log  *log.Entry
	conn net.Conn
}

func NewSocketClient(conn net.Conn, logger *log.Entry) Client {
	return &SocketClient{
		conn: conn,
		log: logger.WithFields(
			log.Fields{
				"client": conn.RemoteAddr(),
			}),
	}
}

func (sc *SocketClient) Log() *log.Entry {
	return sc.log
}

func (sc *SocketClient) Read(p []byte) (n int, err error) {
	return sc.conn.Read(p)
}

func (sc *SocketClient) Write(p []byte) (n int, err error) {
	return sc.conn.Write(p)
}

func (sc *SocketClient) Close() {
	sc.conn.Close()
}

// should run in a deferred context, as it is using 'recover'
func (sc *SocketClient) HandleErrors() {
	err := recover()
	if err == nil {
		return
	}

	switch err := err.(type) {
	case *RequestError:
		if err.code >= 500 {
			sc.log.Debugf("panic: %v\n%s", err, debug.Stack())
		}
		sc.log.Warning(err)
		err.Format(sc)
	default:
		sc.log.Debugf("panic: %v\n%s", err, debug.Stack())
		sc.log.Error(err)
	}
}
