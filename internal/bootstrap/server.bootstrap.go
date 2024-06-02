package bootstrap

import (
	"fmt"
	"net"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/krobus00/websocket-service/internal/config"
	"github.com/krobus00/websocket-service/internal/contract"
	epollerSvc "github.com/krobus00/websocket-service/internal/service/epoller"
	websocketSvc "github.com/krobus00/websocket-service/internal/service/websocket"
	"github.com/sirupsen/logrus"
)

func StartServer() {
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
			logrus.Error(err)
			return
		}

		conn.SetDeadline(time.Now().Add(config.WebsocketDeadline()))
		if err := epollerService.Add(conn); err != nil {
			logrus.Error(err)
			_ = conn.Close()
		}
	})

	err = router.Run(fmt.Sprintf("0.0.0.0:%s", config.ServerPort()))
	if err != nil {
		logrus.Fatal(err)
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
					logrus.Error(err)
				} else {
					logrus.Error(err)
				}
				if err := epollerService.Remove(conn); err != nil {
					logrus.Error(err)
				}
				_ = conn.Close()
				continue
			}

			err = websocketService.HandleIncomingMessage(conn, messageData, opCode)
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
