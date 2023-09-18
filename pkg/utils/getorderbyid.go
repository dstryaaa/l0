package utils

import (
	"database/sql"

	"github.com/dstryaaa/l0/pkg/models"
)

func GetOrderByOrderID(orderUID string, db *sql.DB) (models.Order, error) {
	var order models.Order
	err := db.QueryRow(`
	SELECT
		o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id, o.delivery_service,
		o.shardkey, o.sm_id, o.date_created, o.oof_shard,
		d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
		p.transaction, p.request_id, p.currency, p.provider,
		p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
		FROM orders o
		LEFT JOIN delivery d ON o.order_uid = d.order_uid
		LEFT JOIN payment p ON o.order_uid = p.order_uid
		WHERE o.order_uid = $1
	`, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated,
		&order.OofShard, &order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region, &order.Delivery.Email,
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency, &order.Payment.Provider,
		&order.Payment.Amount, &order.Payment.PaymentDt, &order.Payment.Bank, &order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal, &order.Payment.CustomFee,
	)
	if err != nil {
		return models.Order{}, err
	}
	rows, err := db.Query(`
	SELECT
		chrt_id, track_number, price, rid, name, sale, size,
		total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`, orderUID)
	if err != nil {
		return models.Order{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return models.Order{}, err
		}

		order.Items = append(order.Items, item)
	}

	err = rows.Err()
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}
