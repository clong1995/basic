package mqtt

import (
	"basic/color"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
)

type (
	Client struct {
		Broker   string
		Port     int
		Id       string
		Username string
		Password string
	}

	SubscribeCallback func(topic string, Message []byte)
	mqttClient        struct {
	}
)

var (
	MqttClient      *mqttClient
	client          MQTT.Client
	subscribeMap    = map[string]SubscribeCallback{}
	globalSubscribe SubscribeCallback
)

//Publish 发送消息
func (c mqttClient) Publish(topic, message string) {
	token := client.Publish(topic, 1, false, message)
	token.Wait()
}

//subscribe 监听
func (c mqttClient) subscribe(topic string, callback SubscribeCallback) {
	if token := client.Subscribe(topic, 1, func(client MQTT.Client, msg MQTT.Message) {
		callback(msg.Topic(), msg.Payload())
	}); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	}
}

func AddSubscribe(topic string, callback SubscribeCallback) {
	subscribeMap[topic] = callback
}

func AddGlobalSubscribe(callback SubscribeCallback) {
	globalSubscribe = callback
}

/*var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}*/

func (s Client) Run() error {
	opts := MQTT.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", s.Broker, s.Port))
	opts.SetClientID(s.Id)
	opts.SetUsername(s.Username)
	opts.SetPassword(s.Password)
	opts.SetCleanSession(false)

	/*log.Println(s.Id)
	log.Println(fmt.Sprintf("mqtt://%s:%d", s.Broker, s.Port))
	log.Println(s.Username)
	log.Println(s.Password)*/

	//全局消息处理
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
		if globalSubscribe != nil {
			globalSubscribe(msg.Topic(), msg.Payload())
		}
	})

	opts.OnConnect = func(client MQTT.Client) {
		color.Success(fmt.Sprintf("[mqtt client] %s connected %s:%d success", s.Id, s.Broker, s.Port))
	}
	opts.OnConnectionLost = func(client MQTT.Client, err error) {
		color.Fail(fmt.Sprintf("[mqtt client] connect lost. err: %v", err))
	}

	client = MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		color.Fail(fmt.Sprintf("[mqtt client] connect err: %v", token.Error()))
		return token.Error()
	}
	MqttClient = new(mqttClient)

	//注册监听
	for s2, callback := range subscribeMap {
		go MqttClient.subscribe(s2, callback)
	}
	return nil
}
