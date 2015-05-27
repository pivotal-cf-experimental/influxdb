package engine

import "gopkg.in/pivotal-cf-experimental/influxdb.v0/protocol"

type PointRange struct {
	startTime int64
	endTime   int64
}

func (self *PointRange) UpdateRange(point *protocol.Point) {
	timestamp := *point.GetTimestampInMicroseconds()
	if timestamp < self.startTime {
		self.startTime = timestamp
	}
	if timestamp > self.endTime {
		self.endTime = timestamp
	}
}
