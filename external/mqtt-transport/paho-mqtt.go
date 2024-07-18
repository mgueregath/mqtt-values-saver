package mqtt_transport

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type PahoMqtt struct {
	clientId   string
	url        string
	port       *string
	username   *string
	password   *string
	opts       mqtt.ClientOptions
	pahoClient mqtt.Client
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {

	//fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Println("Connect lost: %v", err)
}

func NewPahoMqtt(clientId string, url string, port *string, username *string, password *string) *PahoMqtt {
	pahoMqtt := PahoMqtt{clientId: clientId, url: url, port: port, username: username, password: password}

	opts := mqtt.NewClientOptions()
	if port != nil {
		opts.AddBroker(fmt.Sprintf("tcp://%s:%d", url, port))
	} else {
		opts.AddBroker(url)
	}
	opts.SetClientID(clientId)
	if (username != nil) && (password != nil) {
		opts.SetUsername(*username)
		opts.SetPassword(*password)
	}
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnectionLost = connectLostHandler
	opts.SetMaxReconnectInterval(1 * time.Second)

	pahoMqtt.opts = *opts

	s, _ := json.MarshalIndent(opts, "", "\t")
	fmt.Printf(string(s))

	return &pahoMqtt
}

type callback func(topic string, message []byte)

func (client *PahoMqtt) Connect(topics []string, fn callback) error {
	client.opts.SetOnConnectHandler(func(c mqtt.Client) {
		fmt.Println("Connected")

		client.subscribe(topics, fn)
	})
	client.pahoClient = mqtt.NewClient(&client.opts)
	if token := client.pahoClient.Connect(); token.Wait() && token.Error() != nil {
		println(token.Error())
	}
	return nil
}

func (client *PahoMqtt) subscribe(topics []string, fn callback) {
	topicsMap := make(map[string]byte, 0)
	for i := 0; i < len(topics); i += 2 {
		topicsMap[topics[i]] = 1
	}
	client.pahoClient.SubscribeMultiple(topicsMap, func(c mqtt.Client, message mqtt.Message) {
		fn(message.Topic(), message.Payload())
	})

}

/*

 */
