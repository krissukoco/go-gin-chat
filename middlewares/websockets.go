package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	wsUpgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Skip checking origin
			return true
		},
	}
)

func WebsocketMiddleware(c *gin.Context) {
	// Upgrade request to websocket
	ws, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	defer ws.Close()
	c.Set("ws", ws)
	c.Next()
}
