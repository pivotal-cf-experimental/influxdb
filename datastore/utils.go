package datastore

import (
	"code.google.com/p/goprotobuf/proto"
	"code.google.com/p/log4go"

	"gopkg.in/pivotal-cf-experimental/influxdb.v0/engine"
	"gopkg.in/pivotal-cf-experimental/influxdb.v0/protocol"
)

func yieldToProcessor(s *protocol.Series, p engine.Processor, aliases []string) (bool, error) {
	for _, alias := range aliases {
		series := &protocol.Series{
			Name:   proto.String(alias),
			Fields: s.Fields,
			Points: s.Points,
		}
		log4go.Debug("Yielding to %s %s", p.Name(), series)
		if ok, err := p.Yield(series); !ok || err != nil {
			return ok, err
		}
	}
	return true, nil
}
