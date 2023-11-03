package socket

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/contrib/websocket"
)


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

        messageStr := string(message)
        log.Println("websocket message received:", messageStr)

        parts := strings.Split(messageStr, " || ")
        data := make(map[string]string)

        for _, part := range parts {
          keyValue := strings.Split(part, ": ")
          if len(keyValue) == 2 {
            key := strings.TrimSpace(keyValue[0])
            value := strings.TrimSpace(keyValue[1])
            data[key] = value
          }
        }

        id := data["id"]
        title := data["title"]
        action := data["action"]

        if (action == "delete") {
          messageToBeDeleted := fmt.Sprintf("id: %s || title: %s || action: normal", id, title)
          deleteMessage <- messageToBeDeleted
          fmt.Println("El mensaje se elimino")
        } else {
          broadcast <- messageStr
        }
      }
    }
}
