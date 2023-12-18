package httpserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"wbtech1421/iternal/cache"
	"wbtech1421/iternal/postgres"
)

var mu sync.Mutex

func getOrderHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	orderID := r.URL.Query().Get("id")
	if orderID == "" {
		http.Error(w, "Не указан идентификатор заказа", http.StatusBadRequest)
		return
	}

	order, found := cache.GetFromCache(orderID)
	if !found {
		// Если заказ не найден в кеше, попробет восстановить из бд
		order, err := postgres.GetOrder(orderID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Ошибка при получении заказа из БД: %v", err), http.StatusInternalServerError)
			return
		}

		if order == nil {
			http.Error(w, "Заказ не найден", http.StatusNotFound)
			return
		}

		// Сохранение в кеш
		cache.SaveToCache(*order)
	}

	// Отправка заказа в ответ
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(order)
	if err != nil {
		return
	}
}

// StartHTTPServer запуск HTTP-сервер
func StartHTTPServer() {
	// Инициализация HTTP-сервера
	http.HandleFunc("/order", getOrderHandler)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Printf("Ошибка при запуске HTTP-сервеа: %v\n", err)
		}
	}()
}
