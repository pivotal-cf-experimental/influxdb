package udp

import (
	"bytes"
	"encoding/json"
	"net"

	log "code.google.com/p/log4go"
	"gopkg.in/pivotal-cf-experimental/influxdb.v0/api"
	"gopkg.in/pivotal-cf-experimental/influxdb.v0/cluster"
	. "gopkg.in/pivotal-cf-experimental/influxdb.v0/common"
	"gopkg.in/pivotal-cf-experimental/influxdb.v0/protocol"
)

type Server struct {
	listenAddress string
	database      string
	coordinator   api.Coordinator
	clusterConfig *cluster.ClusterConfiguration
	conn          *net.UDPConn
	user          *cluster.ClusterAdmin
	shutdown      chan bool
}

func NewServer(listenAddress string, database string, coord api.Coordinator, clusterConfig *cluster.ClusterConfiguration) *Server {
	self := &Server{}

	self.listenAddress = listenAddress
	self.database = database
	self.coordinator = coord
	self.shutdown = make(chan bool, 1)
	self.clusterConfig = clusterConfig

	return self
}

func (self *Server) getAuth() {
	// just use any (the first) of the list of admins.
	names := self.clusterConfig.GetClusterAdmins()
	if len(names) == 0 {
		panic("UDP Plugin: No cluster admins found - couldn't initialize.")
	}
	self.user = self.clusterConfig.GetClusterAdmin(names[0])
}

func (self *Server) ListenAndServe() {
	var err error

	self.getAuth()

	addr, err := net.ResolveUDPAddr("udp4", self.listenAddress)
	if err != nil {
		log.Error("UDPServer: ResolveUDPAddr: ", err)
		return
	}

	if self.listenAddress != "" {
		self.conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			log.Error("UDPServer: Listen: ", err)
			return
		}
	}
	defer self.conn.Close()
	self.HandleSocket(self.conn)
}

func (self *Server) HandleSocket(socket *net.UDPConn) {
	buffer := make([]byte, 2048)

	for {
		n, _, err := socket.ReadFromUDP(buffer)
		if err != nil || n == 0 {
			log.Error("UDP ReadFromUDP error: %s", err)
			continue
		}

		serializedSeries := []*SerializedSeries{}
		decoder := json.NewDecoder(bytes.NewBuffer(buffer[0:n]))
		decoder.UseNumber()
		err = decoder.Decode(&serializedSeries)
		if err != nil {
			log.Error("UDP json error: %s", err)
			continue
		}

		for _, s := range serializedSeries {
			if len(s.Points) == 0 {
				continue
			}

			series, err := ConvertToDataStoreSeries(s, SecondPrecision)
			if err != nil {
				log.Error("UDP cannot convert received data: %s", err)
				continue
			}

			serie := []*protocol.Series{series}
			err = self.coordinator.WriteSeriesData(self.user, self.database, serie)
			if err != nil {
				log.Error("UDP cannot write data: %s", err)
				continue
			}
		}

	}

}
