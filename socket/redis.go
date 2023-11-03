package socket

import (
	"sync"

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
