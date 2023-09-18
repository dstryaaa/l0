package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dstryaaa/l0/pkg/models"
)

func LoadOrdersToCache(db *sql.DB, cache *sync.Map) error {
	rows, err := db.Query(`
	SELECT
		o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service,
		o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		p.transaction, p.request_id, p.currency, p.provider,
		p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		LEFT JOIN delivery d ON o.order_uid = d.order_uid
		LEFT JOIN payment p ON o.order_uid = p.order_uid
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated,
			&order.OofShard, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
			&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
			&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
			&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal, &order.Payment.CustomFee,
		)
		if err != nil {
			return err
		}

		// Загрузка элементов
		itemsRows, err := db.Query(`
		SELECT
			chrt_id, track_number, price, rid, name, sale, size,
			total_price, nm_id, brand, status
			FROM items
			WHERE order_uid = $1
		`, order.OrderUID)
		if err != nil {
			return err
		}

		defer itemsRows.Close()

		for itemsRows.Next() {
			var item models.Item
			err := itemsRows.Scan(
				&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
				&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
			)
			if err != nil {
				return err
			}

			order.Items = append(order.Items, item)
		}

		err = itemsRows.Err()
		if err != nil {
			return err
		}

		cache.Store(order.OrderUID, order)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func ToPrettyJSON(value interface{}) string {
	jsonBytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("<error: %v>", err)
	}
	return string(jsonBytes)
}
