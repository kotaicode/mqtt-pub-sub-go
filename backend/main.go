package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ping struct {
	PingedAt      time.Time `json:"pinged_at"`
	Authenticated bool      `json:"authenticated"`
}

var auth string
var logger *log.Logger
var history = map[string][]ping{}

var publishHalder mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	a := strings.Split(string(msg.Payload()), "::")
	clientID := a[0]
	authenticated := a[1] == auth
	txt := a[2]
	if _, ok := history[clientID]; !ok {
		history[clientID] = []ping{}
	}
	history[clientID] = append(history[clientID], ping{PingedAt: time.Now(), Authenticated: authenticated})
	logger.Printf("TOPIC: %s, CLIENT: %s, MSG: %s\n", msg.Topic(), clientID, txt)
}

func mustRead(env string) string {
	s := os.Getenv(env)
	if s == "" {
		panic(errors.New(env + " is not set"))
	}
	return s
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "ok")
}

func getHistory(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(&history)
}

func main() {
	topic := mustRead("SUBSCRIBE_TOPIC")
	port := mustRead("PORT")
	brokerURL := mustRead("BROKER_URL")
	auth = mustRead("AUTH")

	logger = log.New(os.Stdout, "", log.LstdFlags)
	logPath := os.Getenv("LOG_PATH")
	if logPath != "" {
		if f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			logger.Println(err)
			os.Exit(1)
		} else {
			mw := io.MultiWriter(os.Stdout, f)
			logger = log.New(mw, "", log.LstdFlags)
			defer f.Close()
		}
	}

	opts := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID("backend")
	opts.SetDefaultPublishHandler(publishHalder)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Printf("Unable to setup MQTT client: %v", token.Error())
		os.Exit(1)
	}
	defer c.Disconnect(250)

	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		logger.Printf("Unable to subscribe: %v", token.Error())
		os.Exit(1)
	}

	allowHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	allowOrigins := handlers.AllowedOrigins([]string{"*"})
	allowMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler).Methods("GET")
	r.HandleFunc("/history", getHistory).Methods("GET")

	logger.Printf("Running server at port: %s\n", port)
	logger.Fatalf("%v", http.ListenAndServe(":"+port, handlers.CORS(allowOrigins, allowHeaders, allowMethods)(r)))
}
