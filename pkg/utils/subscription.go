package utils

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	"github.com/dstryaaa/l0/pkg/models"
	"github.com/nats-io/stan.go"
)

func StanChanelSubscription(sc stan.Conn, db *sql.DB, sub string, cache *sync.Map) (stan.Subscription, error) {
	subscription, err := sc.Subscribe(sub, func(msg *stan.Msg) {
		var order models.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println("Failed to Unmarshal user:", err)
			return
		}

		cache.Store(order.OrderUID, order)
		err = SaveOrder(db, order)
		if err != nil {
			log.Println("Failed to save user to DB:", err)
			return
		}
		log.Printf("User Created: %v", order)
	}, stan.DurableName("user-created-durable"))
	return subscription, err
}
