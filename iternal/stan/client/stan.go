package client

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"log"
	"wbtech1421/iternal/data"
)

const (
	clusterID = "test-cluster"
	clientID  = "example-client"
	subject   = "TEST"
)

var sc stan.Conn

func init() {
	var err error
	sc, err = stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Не удалось подключиться к серверу: %v", err)
	}
}

func SubscribeToOrders(callback func(data.Order)) {
	subscription := "example-subscription"

	_, err := sc.Subscribe(subject, func(m *stan.Msg) {
		var receivedMsg data.Order
		err := json.Unmarshal(m.Data, &receivedMsg)
		if err != nil {
			log.Printf("Не удалось размаршаллировать сообщение: %v", err)
			return
		}

		fmt.Printf("Получено сообщение: %+v\n", receivedMsg)
		callback(receivedMsg)
	}, stan.SetManualAckMode(), stan.DurableName(subscription))

	if err != nil {
		log.Fatalf("Не удалось подписаться на сообщения: %v", err)
	}

	fmt.Println("Ожидание сообщений...")
}

func PublishOrderMessage(msg data.Order) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("не удалось маршаллировать сообщение: %v", err)
	}

	err = sc.Publish(subject, msgBytes)
	if err != nil {
		return fmt.Errorf("не удалось отправить сообщение: %v", err)
	}

	fmt.Println("Сообщение отправлено")

	return nil
}
