package postgres

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"wbtech1421/iternal/data"
)

//const (
//	host     = "localhost"
//	port     = 5432
//	user     = "admin"
//	password = "admin"
//	dbname   = "wbtech_db"
//	ssl      = "disable"
//)

type Order2 struct {
	ID   string `gorm:"primaryKey"`
	Data []byte
}

func SaveOrder(order data.Order) error {
	dsn := "host=localhost user=admin password=admin dbname=wbtech_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	orderData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = db.Table("orders").Create(&Order2{ID: order.OrderUID, Data: orderData}).Error
	if err != nil {
		return err
	}

	return nil
}

func InitTable() error {
	dsn := "host=localhost user=admin password=admin dbname=wbtech_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Создание таблицы order, если ее нет
	err = db.Exec(`
		CREATE TABLE IF NOT EXISTS public.orders
		(
			id text NOT NULL,
			data bytea,
			PRIMARY KEY (id)
		);
	`).Error
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

//// Пример заказа
//order := data.Order{
//	OrderUID: "12521111415",
//}
//
//err = saveOrder(db, order)
//if err != nil {
//	log.Fatal(err)
//}
//
//ordersdd, err := getOrder(db, order.OrderUID)
//if err != nil {
//	log.Fatal(err)
//}
//fmt.Println(ordersdd)
//
//fmt.Println("Заказ успешно сохранен в базу данныхю.")
