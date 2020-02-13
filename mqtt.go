package main

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
	//log "github.com/sirupsen/logrus"
)

func connect(clid string, uri *url.URL) mqtt.Client {
	opts := createClientOptions(clid, uri)
	cli := mqtt.NewClient(opts)

	token := cli.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}

	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return cli
}

func createClientOptions(clid string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.Broker)

	// opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	// opts.SetUsername(url.User.Username())
	// opts.SetPassword(password)
	opts.SetClientID(clid)
	return opts
}

func mqttReader(uri *url.URL, topic string) {
	cli := connect("sub", uri)
	cli.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		tempstr := string(msg.Payload())
		fmt.Printf("* [%s] %s\n", msg.Topic(), tempstr)

		// send the weather to a websocket if we have one
		tlv := NewTLV(tlvTypeTempf, len(tempstr)+2, tempstr)
		fmt.Println("got our tlv")
		webQ <- tlv
	})
}

func mqtt_loop(broker string, wg *sync.WaitGroup) {
	defer wg.Done()

	uri, err := url.Parse(broker)
	if err != nil {
		log.Fatal(err)
	}

	//topic := uri.Path[1:len(uri.Path)]
	topic := uri.Path
	if topic == "" {
		topic = "test"
	}

	go mqttReader(uri, topic)

	u, err := uri.Parse(config.Broker)
	if err != nil {
		log.Fatal(err)
	}

	client := connect("pub", u)
	timer := time.NewTicker(1 * time.Second)
	for t := range timer.C {
		client.Publish("/health/office", 0, false, t.String())
	}
}
