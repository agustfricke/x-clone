package socket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type WebSocketMessage struct {
    ID     string `json:"id"`
    Title  string `json:"title"`
    Action string `json:"action"`
}

func FeedWebsocket(c *websocket.Conn) {
    defer func() {
      unregister <- c
      c.Close()
    }()

    register <- c

    for {
      messageType, message, err := c.ReadMessage()
      if err != nil {
        if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
          log.Println("read error:", err)
        }
        return
      }

      if messageType == websocket.TextMessage {

        var wsMessage WebSocketMessage
        err := json.Unmarshal(message, &wsMessage)
        if err != nil {
          log.Println("Error al decodificar el mensaje JSON:", err)
          continue
          }

        id := wsMessage.ID
        title := wsMessage.Title
        action := wsMessage.Action

        if action == "delete" {
          newMessage := WebSocketMessage{
            ID:     id,
            Title:  title,
            Action: "normal",
          }
          messageToBeDeleted, _ := json.Marshal(newMessage)
          deleteMessage <- string(messageToBeDeleted)
          fmt.Println("El mensaje se eliminó")
        } else {
          // Reenviar el mensaje JSON original
          broadcast <- string(message)
        }
      }
    }
}
