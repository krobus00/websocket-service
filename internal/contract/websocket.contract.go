package contract

import (
	"net"

	"github.com/gobwas/ws"
)

type WebsocketService interface {
	HandleIncomingMessage(conn net.Conn, messageData []byte, operationCode ws.OpCode) error
}
