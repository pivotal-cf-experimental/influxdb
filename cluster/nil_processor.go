package cluster

import (
	"fmt"

	"gopkg.in/pivotal-cf-experimental/influxdb.v0/engine"
	"gopkg.in/pivotal-cf-experimental/influxdb.v0/protocol"
)

type NilProcessor struct{}

func (np NilProcessor) Name() string {
	return "NilProcessor"
}

func (np NilProcessor) Yield(s *protocol.Series) (bool, error) {
	return false, fmt.Errorf("Shouldn't get any data")
}

func (np NilProcessor) Close() error {
	return nil
}

func (np NilProcessor) Next() engine.Processor {
	return nil
}
