package world

import "time"

import (
	"github.com/rcrowley/go-metrics"
	"github.com/rcrowley/go-metrics/influxdb"
)

type stats struct {
	peeps metrics.Gauge
}

func newStats() *stats {
	r := metrics.NewRegistry()

	stats := &stats{
		peeps: metrics.NewGauge(),
	}

	r.Register("peeps", stats.peeps)

	// go metrics.Log(metrics.DefaultRegistry, time.Second*1, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	go influxdb.Influxdb(r, time.Second*1, &influxdb.Config{
		Host:     "127.0.0.1:8086",
		Database: "world",
		Username: "world",
		Password: "peepsrule",
	})
	return stats

}
