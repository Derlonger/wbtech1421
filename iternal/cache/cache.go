package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
	"wbtech1421/iternal/data"
)

var myCache = cache.New(5*time.Minute, 10*time.Minute)

func SaveToCache(order data.Order) {
	myCache.Set(order.OrderUID, order, cache.DefaultExpiration)
}

func GetFromCache(orderUID string) (data.Order, bool) {
	cachedOrder, found := myCache.Get(orderUID)
	if found {
		return cachedOrder.(data.Order), true
	}
	return data.Order{}, false
}
