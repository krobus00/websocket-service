package websocket

import (
	"context"
	"errors"
	"net"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func (s *Service) HandleIncomingMessage(conn net.Conn, messageData []byte, operationCode ws.OpCode) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	processDone := make(chan error)
	go func() {
		defer close(processDone)
		err := wsutil.WriteServerMessage(conn, operationCode, messageData)

		processDone <- err
	}()

	select {
	case <-ctx.Done():
		conn.Close()
		return errors.New("timeout")
	case err := <-processDone:
		return err
	}
}
