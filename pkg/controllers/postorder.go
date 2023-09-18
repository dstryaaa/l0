package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"

	"github.com/dstryaaa/l0/pkg/models"
)

func CreateOrderHandler(db *sql.DB, sc stan.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			var order models.Order
			json.NewDecoder(r.Body).Decode(&order)

			data, _ := json.Marshal(order)

			err := sc.Publish("order-create", data)
			if err != nil {
				http.Error(w, "Failed to publish user.", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		}
	}
}
