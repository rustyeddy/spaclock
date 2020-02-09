package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"sync"

	"github.com/eclipse/paho.mqtt.golang"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func mqtt_loop(broker string, topic string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)

	//opts := mqtt.NewClientOptions().AddBroker("tcp://iot.eclipse.org:1883").SetClientID("gotrivial")
	if broker == "" {
		broker = "tcp://10.24.10.10:1883"
	}

	fmt.Println("MQTT start loop")
	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID("spaclock")

	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	fmt.Println("MQTT new client")
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("MQTT subscribe ")
	if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	fmt.Println("look for tempfs")
	for i := 0; i < 5; i++ {
		text := fmt.Sprintf("tempf is msg #%d!", i)
		token := c.Publish(topic, 0, false, text)
		token.Wait()
	}

	fmt.Println("sleep some and unsub")
	time.Sleep(6 * time.Second)

	fmt.Println("sleep and unsub")
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	c.Disconnect(250)
	time.Sleep(1 * time.Second)
}

