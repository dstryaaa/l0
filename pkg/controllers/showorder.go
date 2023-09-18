package controllers

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/dstryaaa/l0/pkg/models"
	"github.com/dstryaaa/l0/pkg/utils"
	"github.com/dstryaaa/l0/pkg/utils2"
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
)

func serveOrder(w http.ResponseWriter, r *http.Request, order models.Order) {
	tmp2, err := template.ParseFiles("./static/order.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmp2.Execute(w, order)
}

func ShowOrder(db *sql.DB, sc stan.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		orderUID := vars["order_uid"]

		order, err := utils.GetOrderByOrderID(orderUID, db)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("Error while fetching order"))
			return
		}

		data, _ := json.Marshal(order)

		err = sc.Publish("order-post", data)
		if err != nil {
			http.Error(w, "Failed to publish user.", http.StatusInternalServerError)
			return
		}
		log.Println(order.OrderUID)
		utils2.OrderMutex.Lock()
		receivedOrder := utils2.ReceivedOrder
		utils2.OrderMutex.Unlock()

		w.Header().Set("Content-Type", "text/html")
		serveOrder(w, r, receivedOrder)
	}
}
