package bootstrap

import (
	"log"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/krobus00/websocket-service/internal/contract"
	epollerSvc "github.com/krobus00/websocket-service/internal/service/epoller"
	websocketSvc "github.com/krobus00/websocket-service/internal/service/websocket"
)

const (
	timeout = 5 * time.Second
)

func StartServer() {
	var err error
	epollerService, err := epollerSvc.New()
	if err != nil {
		panic(err)
	}

	websocketService := websocketSvc.
		New().
		WithEpollerService(epollerService)

	go Start(epollerService, websocketService)

	router := gin.Default()

	router.GET("/ws", func(c *gin.Context) {
		conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
		if err != nil {
			return
		}

		conn.SetDeadline(time.Now().Add(timeout))
		if err := epollerService.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	})

	if err := router.Run("0.0.0.0:8000"); err != nil {
		log.Fatal(err)
	}

}

func Start(epollerService contract.EpollerService, websocketService contract.WebsocketService) {
	for {
		connections, err := epollerService.Wait()
		if err != nil {
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}

			messageData, opCode, err := wsutil.ReadClientData(conn)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					log.Printf("read timeout for connection %v", conn)
				} else {
					log.Printf("read error for connection %v: %v", conn, err)
				}
				if err := epollerService.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
				conn.Close()
				continue
			}

			err = websocketService.HandleIncomingMessage(conn, messageData, opCode)
			if err != nil {
				log.Printf("failed to handle in coming message %v", err)
			}
		}
	}
}
