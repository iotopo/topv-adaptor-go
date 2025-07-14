package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"log/slog"
	"math/rand"
	"time"
)

var nc *nats.Conn
var randGen = rand.New(rand.NewSource(time.Now().UnixNano()))

func pushRealtimeValue(item *ValueItem) {
	payload, err := json.Marshal(item)
	if err != nil {
		slog.Error("Error marshaling JSON: %v", err)
		return
	}
	err = nc.Publish(fmt.Sprintf("rtdb.iotopo.%s", item.Tag), payload)
	if err != nil {
		slog.Warn("Error publishing to nats: %v", err)
	}
}

// 实时数据推送
func realPush(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			ts := time.Now()
			for i := 0; i < 3; i++ {
				for j := 0; j < 10; j++ {
					// 生成1-100之间的随机数
					value := randGen.Float64()*99 + 1
					pushRealtimeValue(&ValueItem{
						Value: value,
						Tag:   fmt.Sprintf("group%d.dev%d.a", i+1, j+1),
						Time:  ts,
					})
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func init() {
	opts := []nats.Option{
		nats.MaxReconnects(-1),
		nats.RetryOnFailedConnect(true),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			slog.Info("NATS reconnecting...")
		}),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			slog.Info("NATS disconnected:", err)
		}),
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			slog.Info("NATS error: %v", err)
		}),
		nats.ConnectHandler(func(conn *nats.Conn) {
			slog.Info("nats connected")
		}),
		nats.Name("topv-adaptor"),
	}

	var err error
	nc, err = nats.Connect("nats://127.0.0.1:4222", opts...)
	if err != nil {
		log.Fatal(err)
	}
}
