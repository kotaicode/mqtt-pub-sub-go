package main

import (
	"errors"
	"fmt"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func mustRead(env string) string {
	s := os.Getenv(env)
	if s == "" {
		panic(errors.New(env + " is not set"))
	}
	return s
}

func main() {
	clientID := mustRead("PUBLISHER_ID")
	topic := mustRead("PUBLISH_TOPIC")
	brokerURL := mustRead("BROKER_URL")
	auth := mustRead("AUTH")

	opts := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID(clientID)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	text := fmt.Sprintf("%s::%s::live-signal", clientID, auth)
	fmt.Println(text)
	token := c.Publish(topic, 0, false, text)
	token.Wait()

	defer c.Disconnect(250)
}
