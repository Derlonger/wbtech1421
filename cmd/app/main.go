package main

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
	"wbtech1421/iternal/cache"

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
		_, found := cache.GetFromCache(msg.OrderUID)
		if found {
			fmt.Println("Заказ найден в кеше")
		} else {
			// Если не найден, сохраняем в бд и в кеш
			err := postgres.SaveOrder(msg)
			if err != nil {
				return
			}

			// Сохранение в кеш
			cache.SaveToCache(msg)

			fmt.Printf("Заказ сохранен")
		}
	})

	// Поддержание работы приложения
	select {}
}

func createOrderMessage() data.Order {
	delivery := data.Delivery{
		Name:   "Test Testov",
		Phone:  "+9720000000",
		Zip:    "2639809",
		City:   "Kiryat Mozkin",
		Adress: "Ploshad Mira 15",
		Region: "Kraiot",
		Email:  "test@gmail.com",
	}

	payment := data.Payment{
		Transaction:  "b563feb7b2b84b6test",
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       1817,
		PaymentDt:    1637907727,
		Bank:         "alpha",
		DeliveryCost: 1500,
		GoodsTotal:   317,
		CustomFee:    0,
	}

	items := []data.Items{
		{
			ChrtID:      9934930,
			TrackNumber: "WBILMTESTTRACK",
			Price:       453,
			RID:         "ab4219087a764ae0btest",
			Name:        "Mascaras",
			Sale:        30,
			Size:        "0",
			TotalPrice:  317,
			NmID:        2389212,
			Brand:       "Vivienne Sabo",
			Status:      202,
		},
	}

	dateCreated, _ := time.Parse(time.RFC3339, "2021-11-26T06:22:19Z")

	order := data.Order{
		OrderUID:          generateUID(),
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              "99",
		DateCreated:       dateCreated,
		OofShard:          "1",
	}

	return order
}

func generateUID() string {
	uid := uuid.New().String()
	return uid
}
