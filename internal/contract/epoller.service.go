package contract

import "net"

type EpollerService interface {
	Add(conn net.Conn) error
	Remove(conn net.Conn) error
	Wait() ([]net.Conn, error)
}
