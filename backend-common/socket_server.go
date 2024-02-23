package common

import (
	socketio "github.com/googollee/go-socket.io"
)

func SetUpSocketServer(l *Logger) *socketio.Server {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext(s.ID())
		l.Message("connected: " + s.ID())
		return nil
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		l.Error("eError:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		l.Message("closed: " + s.ID() + " (" + reason + ")")
	})

	return server
}
