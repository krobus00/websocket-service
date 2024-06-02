package websocket

import (
	"context"
	"errors"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/krobus00/websocket-service/internal/config"
)

func (s *Service) HandleIncomingMessage(conn net.Conn, messageData []byte, operationCode ws.OpCode) error {
	ctx, cancel := context.WithTimeout(context.Background(), config.WebsocketDeadline())
	defer cancel()

	processDone := make(chan error)
	go func() {
		defer close(processDone)
		err := wsutil.WriteServerMessage(conn, operationCode, messageData)

		processDone <- err
	}()

	select {
	case <-ctx.Done():
		_ = conn.Close()
		return errors.New("timeout")
	case err := <-processDone:
		return err
	}
}
