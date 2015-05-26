package batcher

import (
	"sync/atomic"
	"time"

	"github.com/influxdb/influxdb"
)

// Batcher accepts Points and will emit a batch of those points when either
// a) the batch reaches a certain size, or b) a certain time passes.
type Batcher struct {
	size     int
	duration time.Duration

	stats BatcherStats
}

// NewBatcher returns a new Batcher.
func NewBatcher(sz int, d time.Duration) *Batcher {
	return &Batcher{size: sz, duration: d}
}

// BatcherStats are the statistics each batcher tracks.
type BatcherStats struct {
	BatchTotal   uint64 // Total count of batches transmitted.
	PointTotal   uint64 // Total count of points processed.
	SizeTotal    uint64 // Nubmer of batches that reached size threshold.
	TimeoutTotal uint64 // Nubmer of timeouts that occurred.
}

// Start starts the batching process. It should be called from a goroutine.
func (b *Batcher) Start(in <-chan influxdb.Point, out chan<- []influxdb.Point) {
	timer := time.NewTimer(0)
	var batch []influxdb.Point
	var timerCh <-chan time.Time

	for {
		select {
		case p := <-in:
			atomic.AddUint64(&b.stats.PointTotal, 1)
			if batch == nil {
				batch = make([]influxdb.Point, 0, b.size)
				timer.Reset(b.duration)
				timerCh = timer.C
			}

			batch = append(batch, p)
			if len(batch) == b.size {
				atomic.AddUint64(&b.stats.SizeTotal, 1)
				out <- batch
				atomic.AddUint64(&b.stats.BatchTotal, 1)
				batch = nil
				timerCh = nil
			}

		case <-timerCh:
			atomic.AddUint64(&b.stats.TimeoutTotal, 1)
			out <- batch
			atomic.AddUint64(&b.stats.BatchTotal, 1)
			batch = nil
		}
	}
}

// Stats returns a BatcherStats object for the Batcher. While the each statistic should be
// closely correlated with each other statistic, it is not guaranteed.
func (b *Batcher) Stats() *BatcherStats {
	stats := BatcherStats{}
	stats.BatchTotal = atomic.LoadUint64(&b.stats.BatchTotal)
	stats.PointTotal = atomic.LoadUint64(&b.stats.PointTotal)
	stats.SizeTotal = atomic.LoadUint64(&b.stats.SizeTotal)
	stats.TimeoutTotal = atomic.LoadUint64(&b.stats.TimeoutTotal)
	return &stats
}
