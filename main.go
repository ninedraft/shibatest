package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

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
		websocket.Handler(func(conn *websocket.Conn) {
			logrus.Infof("starting websoket stream to %v", conn.RemoteAddr())
			defer conn.Close()
			for {
				conn.SetDeadline(time.Now().Add(4 * time.Second))
				newWow := wow()
				logrus.Infof("sending wow: %s", newWow)
				_, err := conn.Write([]byte(newWow))
				if err != nil {
					logrus.WithError(err).Errorf("error while sending data by ws to %q", conn.RemoteAddr())
				}
				time.Sleep(4 * time.Second)
			}
		}).ServeHTTP(writer, request)
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
