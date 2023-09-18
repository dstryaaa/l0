package utils

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dstryaaa/l0/pkg/models"
)

func SaveOrder(db *sql.DB, order models.Order) error {
	err := insertOrder(db, order)
	if err != nil {
		return err
	}

	order.Delivery.OrderUID = order.OrderUID
	err = saveDelivery(db, order.Delivery)
	if err != nil {
		return err
	}

	order.Payment.OrderUID = order.OrderUID
	err = savePayment(db, order.Payment)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		item.OrderUID = order.OrderUID
		err = saveItem(db, item)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertOrder(db *sql.DB, order models.Order) error {
	parsedDate, err := time.Parse(time.RFC3339, order.DateCreated)
	if err != nil {
		return fmt.Errorf("cannot parse date: %w", err)
	}
	_, err = db.Exec(`INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, parsedDate, order.OofShard)
	return err
}

func saveDelivery(db *sql.DB, delivery models.Delivery) error {
	_, err := db.Exec(`INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		delivery.OrderUID, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email)
	return err
}

func savePayment(db *sql.DB, payment models.Payment) error {
	_, err := db.Exec(`INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		payment.OrderUID, payment.Transaction, payment.RequestID, payment.Currency, payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee)
	return err
}

func saveItem(db *sql.DB, item models.Item) error {
	_, err := db.Exec(`INSERT INTO items (chrt_id, order_uid, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		item.ChrtID, item.OrderUID, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
	return err
}
