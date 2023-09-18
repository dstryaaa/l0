package utils2

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"

	"github.com/dstryaaa/l0/pkg/models"
	"github.com/nats-io/stan.go"
)

var ReceivedOrder models.Order
var OrderMutex sync.Mutex

func StanChanelSubscription2(sc stan.Conn, db *sql.DB, sub string) (stan.Subscription, error) {
	subscription, err := sc.Subscribe(sub, func(msg *stan.Msg) {

		var order models.Order
		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Println("Failed to Unmarshal order:", err)
			return
		}
		OrderMutex.Lock()
		ReceivedOrder = order
		OrderMutex.Unlock()
	}, stan.DurableName("user-created-durable"))

	return subscription, err
}
