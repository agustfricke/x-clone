package socket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type Message struct {
	Repo *MessageStore 
}

func RunHub() {
  for {
    select {
    case connection := <-register:
      clients[connection] = &client{}
      log.Println("connection registered")

      MessageStorage.Lock()
      messages, err := MessageStorage.Client.LRange(MessageStorage.Client.Context(), "messages", 0, -1).Result()
      if err == nil {
        for _, message := range messages {
          connection.WriteMessage(websocket.TextMessage, []byte(message))
        }
      }
      MessageStorage.Unlock()

    case message := <-broadcast:
      log.Println("message received:", message)

      var messageJson WebSocketMessage
      json.Unmarshal([]byte(message), &messageJson)
      fmt.Println("Heyy", messageJson)

      MessageStorage.Lock()
      MessageStorage.Client.LPush(MessageStorage.Client.Context(), "messages", message)
      MessageStorage.Unlock()

      for connection, c := range clients {
        go func(connection *websocket.Conn, c *client) {
          c.mu.Lock()
          defer c.mu.Unlock()
          if c.isClosing {
            return
          }
          if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
            c.isClosing = true
            log.Println("write error:", err)

            connection.WriteMessage(websocket.CloseMessage, []byte{})
            connection.Close()
            unregister <- connection
          }
        }(connection, c)
      }

    case message := <-deleteMessage:
      log.Println("message to delete:", message)

      MessageStorage.Lock()
      if err := MessageStorage.Client.LRem(MessageStorage.Client.Context(), "messages", 0, message).Err(); err != nil {
        log.Println("error deleting message:", err)
      }
      MessageStorage.Unlock()

      for connection, c := range clients {
        go func(connection *websocket.Conn, c *client) {
          c.mu.Lock()
          defer c.mu.Unlock()
          if c.isClosing {
            return
          }
          if err := connection.WriteMessage(websocket.TextMessage, []byte("Deleted message: "+message)); err != nil {
            c.isClosing = true
            log.Println("error al eliminar write error:", err)

            connection.WriteMessage(websocket.CloseMessage, []byte{})
            connection.Close()
            unregister <- connection
          }
        }(connection, c)
      }
    }
  }
}

