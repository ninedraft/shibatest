package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RubyDate,
		FullTimestamp:   true,
	})

	http.HandleFunc("/shiba", func(writer http.ResponseWriter, request *http.Request) {
		logrus.Infof("uploading shiba")
		fmt.Fprintf(writer, "%s\nWow, much fast connection, %s :3", shiba, request.RemoteAddr)
	})

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := upgrader.Upgrade(writer, request, request.Header)
		if err != nil {
			logrus.WithError(err).Errorf("websocket upgrade failed")
			http.Error(writer, "websocket upgrade failed", http.StatusInternalServerError)
			return
		}

		defer conn.Close()
		for {
			conn.SetWriteDeadline(time.Now().Add(4 * time.Second))
			newWow := wow()
			logrus.Infof("sending wow: %s", newWow)
			err := conn.WriteMessage(websocket.TextMessage, []byte(newWow))
			if err != nil {
				logrus.WithError(err).Errorf("error while sending data by ws to %q", conn.RemoteAddr())
			}
			time.Sleep(4 * time.Second)
		}
	})
	go func() {
		for range time.Tick(3 * time.Second) {
			logrus.Info(wow())
		}
	}()
	logrus.Infof("serving shiba")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logrus.Fatal(err)
	}
}
