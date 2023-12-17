package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"

	"wbtech1421/iternal/data"
	"wbtech1421/iternal/postgres"
	"wbtech1421/iternal/stan/client"
)

func main() {
	//инизацилация таблице в бд
	fmt.Println("INIT TABLE RESULT")
	err := postgres.InitTable()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("INIT FROM DB")

	go func() {
		for {
			msg := createOrderMessage()
			err = client.PublishOrderMessage(msg)
			if err != nil {
				log.Fatalf("Ошибка при отправке заказа: %v", err)
			}
			time.Sleep(10 * time.Second) // Задержка перед отправкой сообщения
		}
	}()

	client.SubscribeToOrders(func(msg data.Order) {
		err := postgres.SaveOrder(msg)
		if err != nil {
			return
		}
		fmt.Printf("Заказ сохранен")
	})

	// Поддержание работы приложения
	select {}
}

func createOrderMessage() data.Order {
	msg := data.Order{
		OrderUID: generateUID(),
		// todo заполнить остальные поля
	}
	return msg

}

func generateUID() string {
	uid := uuid.New().String()
	return uid
}
