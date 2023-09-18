package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dstryaaa/l0/pkg/routes"
	"github.com/dstryaaa/l0/pkg/utils"
	"github.com/dstryaaa/l0/pkg/utils2"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var MyCache sync.Map

func main() {
	var err error
	db, err := utils.ConnectToPostgres()
	if err != nil {
		log.Fatal(err)
	}

	err = utils.LoadOrdersToCache(db, &MyCache)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}

	// MyCache.Range(func(key, value interface{}) bool {
	// 	fmt.Printf("Key: %s\nValue: %s\n", key, utils.ToPrettyJSON(value))
	// 	return true
	// })

	sc, err := utils.ConnectToStan()
	if err != nil {
		log.Fatal(err)
	}

	postSubscription, err := utils.StanChanelSubscription(sc, db, "order-create", &MyCache)
	if err != nil {
		log.Fatal(err)
	}

	getSubscription, err := utils2.StanChanelSubscription2(sc, db, "order-post")
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	routes.UserRoutes(r, sc, db)
	http.Handle("/", r)
	log.Println("Starting server on :8080")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", r))
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	log.Println("Shutting down...")

	sc.Close()
	postSubscription.Close()
	getSubscription.Close()
}
