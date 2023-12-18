package postgres

import (
	"encoding/json"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"wbtech1421/iternal/data"
)

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

func GetOrder(orderID string) (*data.Order, error) {
	dsn := "host=localhost user=admin password=admin dbname=wbtech_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var order2 Order2
	if err := db.Table("orders").Where("id = ?", orderID).First(&order2).Error; err != nil {
		return nil, err
	}

	var order data.Order
	err = json.Unmarshal(order2.Data, &order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}
