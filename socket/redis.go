package socket

import (
	"sync"
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/contrib/websocket"
)

type client struct {
  isClosing bool
  mu        sync.Mutex
}

var clients = make(map[*websocket.Conn]*client)
var register = make(chan *websocket.Conn)
var unregister = make(chan *websocket.Conn)

var broadcast = make(chan string)
var deleteMessage = make(chan string)

type MessageStore struct {
  sync.Mutex
  Client *redis.Client
}

var MessageStorage *MessageStore

func InitializeRedis() (*redis.Client, error) {
  client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379", 
    DB:   0,               
  })
  _, err := client.Ping(client.Context()).Result()
  return client, err
}

func (r *MessageStore) Insert(ctx context.Context, message WebSocketMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to encode order: %w", err)
	}

	key := message.ID

	txn := r.Client.TxPipeline()

	res := txn.SetNX(ctx, key, string(data), 0)
	if err := res.Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to set: %w", err)
	}

	if err := txn.SAdd(ctx, "messages", key).Err(); err != nil {
		txn.Discard()
		return fmt.Errorf("failed to add to orders set: %w", err)
	}

	if _, err := txn.Exec(ctx); err != nil {
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}

