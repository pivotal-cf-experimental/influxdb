package coordinator

import (
	"net"

	"gopkg.in/pivotal-cf-experimental/influxdb.v0/protocol"
)

type Handler interface {
	HandleRequest(*protocol.Request, net.Conn) error
}
