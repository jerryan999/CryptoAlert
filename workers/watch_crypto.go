package workers

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/jerryan999/CryptoAlert/utils"
	"github.com/labstack/gommon/log"
)

var QUEUE_KEY_WATCH_CRYPTO = "tasks.send_email"
var QUEUE_ADDR_WATCH_CRYPTO = "localhost:6379"

// redis
var rdb *redis.Client

func WatchCryptoWorker() {
	ctx := context.TODO()

	// we only care about crypto here
	crypto := "bitcoin"

	conn, _, err := websocket.DefaultDialer.Dial("wss://ws.coincap.io/prices?assets="+crypto, nil)
	if err != nil {
		log.Fatal(err)
	}

	// read from the websocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		// parse the json message
		var data map[string]string
		json.Unmarshal(message, &data)
		log.Info("Got data from websocket:", data)

		// get the price
		currentPrice := data[crypto]

		// Get alert key
		key_price_gt := utils.GetAlertQueueKey(crypto, true)
		key_price_lt := utils.GetAlertQueueKey(crypto, false)

		// Get alert ids which needs to be sent
		res := rdb.ZRangeByScore(ctx, key_price_gt, &redis.ZRangeBy{
			Min: "0",
			Max: currentPrice,
		}).Val()

		// remove the alert from redis
		// convert []string to []interface{} first
		members := make([]interface{}, len(res))
		for i, v := range res {
			members[i] = v
		}
		rdb.ZRem(ctx, key_price_gt, members...)

		// add to send email queue
		for _, v := range res {
			PushRequest(map[string]interface{}{"id": v, "price": currentPrice, "crypto": crypto, "direction": 0})
		}

		res = rdb.ZRangeByScore(ctx, key_price_lt, &redis.ZRangeBy{
			Min: currentPrice,
			Max: "inf",
		}).Val()

		// remove the alert from redis
		// convert []string to []interface{} first
		members = make([]interface{}, len(res))
		for i, v := range res {
			members[i] = v
		}
		rdb.ZRem(ctx, key_price_lt, members...)

		// remove the alert from redis
		for _, v := range res {
			PushRequest(map[string]interface{}{"id": v, "price": currentPrice, "crypto": crypto, "direction": 1})
		}
	}
}

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: QUEUE_ADDR_WATCH_CRYPTO,
	})
}
