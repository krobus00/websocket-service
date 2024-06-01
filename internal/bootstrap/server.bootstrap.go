package bootstrap

import (
	"log"
	"net"
	"time"

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

	ln, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Fatal(err)
	}

	u := ws.Upgrader{}

	for {
		// zero allocation upgrade
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		_, err = u.Upgrade(conn)
		if err != nil {
			log.Printf("upgrade error: %s", err)
		}
		if err := epollerService.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
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
			conn.SetDeadline(time.Now().Add(timeout))

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
