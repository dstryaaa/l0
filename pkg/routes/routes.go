package routes

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"

	"github.com/dstryaaa/l0/pkg/controllers"
	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router, sc stan.Conn, db *sql.DB) {

	r.HandleFunc("/order/show", controllers.ServeHTML)
	r.HandleFunc("/order/styles.css", controllers.ServeCSS)
	r.HandleFunc("/order/show/submit", controllers.SubmitForm)
	r.HandleFunc("/order/show/{order_uid}", controllers.ShowOrder(db, sc)).Methods("GET")
	r.HandleFunc("/order/create", controllers.CreateOrderHandler(db, sc))

}
