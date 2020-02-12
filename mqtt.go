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

	fmt.Printf("client %+v\n", cli)

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
	opts.AddBroker("tcp://10.24.10.10:1883")

	//opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	// opts.SetUsername(url.User.Username())
	// opts.SetPassword(password)
	opts.SetClientID(clid)
	return opts
}

func listen(uri *url.URL, topic string) {
	cli := connect("sub", uri)
	cli.Subscribe("/topic/tempf", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}

func mqtt_loop(broker string, wg *sync.WaitGroup) {
	defer wg.Done()

	uri, err := url.Parse(broker)
	if err != nil {
		log.Fatal(err)
	}

	topic := uri.Path[1:len(uri.Path)]
	if topic == "" {
		topic = "test"
	}

	go listen(uri, topic)

	u, err := uri.Parse("tcp://10.24.10.10:1883")
	if err != nil {
		log.Fatal(err)
	}

	client := connect("pub", u)
	timer := time.NewTicker(1 * time.Second)
	for t := range timer.C {
		client.Publish("/health/office", 0, false, t.String())
	}
}
