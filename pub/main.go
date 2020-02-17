package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var shouldFail bool
var counter int
var logger *log.Logger

func statusHandler(w http.ResponseWriter, r *http.Request) {
	txt := "ok"
	if shouldFail {
		w.WriteHeader(400)
		txt = "fail"
	} else {
		w.WriteHeader(200)
	}
	fmt.Fprintf(w, txt)
}

func makeItFail(w http.ResponseWriter, r *http.Request) {
	shouldFail = true
	w.WriteHeader(200)
	fmt.Fprintf(w, "next statusHandler call will fail")
}

func mustRead(env string) string {
	s := os.Getenv(env)
	if s == "" {
		panic(errors.New(env + " is not set"))
	}
	return s
}

func startFiringLiveSignals(c mqtt.Client, interval int, topic, clientID, auth string) {
	for {
		<-time.After(time.Duration(interval) * time.Second)
		go func() {
			text := fmt.Sprintf("%s::%s::live-signal-%d", clientID, auth, counter)
			counter++
			logger.Println(text)
			token := c.Publish(topic, 0, false, text)
			token.Wait()
		}()
	}
}

func main() {
	port := mustRead("PORT")
	clientID := mustRead("PUBLISHER_ID")
	topic := mustRead("PUBLISH_TOPIC")
	brokerURL := mustRead("BROKER_URL")
	auth := mustRead("AUTH")
	i := mustRead("PUBLISH_INTERVAL")

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

	opts := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID(clientID)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer c.Disconnect(250)

	interval, err := strconv.Atoi(i)
	if err != nil {
		panic(err)
	}
	go startFiringLiveSignals(c, interval, topic, clientID, auth)

	allowHeaders := handlers.AllowedHeaders([]string{"Content-Type"})
	allowOrigins := handlers.AllowedOrigins([]string{"*"})
	allowMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE"})

	r := mux.NewRouter()
	r.HandleFunc("/status", statusHandler).Methods("GET")
	r.HandleFunc("/make_it_fail", makeItFail).Methods("PUT")

	logger.Printf("Running server at port: %s\n", port)
	log.Fatalf("%v", http.ListenAndServe(":"+port, handlers.CORS(allowOrigins, allowHeaders, allowMethods)(r)))
}
