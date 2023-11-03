package database

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/contrib/websocket"
)

type Client struct {
  isClosing bool
  mu        sync.Mutex
}

var Clients = make(map[*websocket.Conn]*Client)
var Register = make(chan *websocket.Conn)
var Unregister = make(chan *websocket.Conn)

var Broadcast = make(chan string)
var DeleteMessage = make(chan string)

type MessageStore struct {
  sync.Mutex
  Client *redis.Client
}

var MessageStorage *MessageStore

func initializeRedis() (*redis.Client, error) {
  client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379", 
    DB:   0,               
  })
  _, err := client.Ping(client.Context()).Result()
  return client, err
}
