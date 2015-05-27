package wal

import (
	"gopkg.in/pivotal-cf-experimental/influxdb.v0/common"
	"gopkg.in/pivotal-cf-experimental/influxdb.v0/protocol"
)

type replayRequest struct {
	requestNumber uint32
	request       *protocol.Request
	shardId       uint32
	startOffset   int64
	endOffset     int64
	err           error
}

func newErrorReplayRequest(err error) *replayRequest {
	return &replayRequest{
		err: common.NewErrorWithStacktrace(err, "Replay error"),
	}
}
