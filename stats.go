package world

import (
  "time"
  "github.com/rcrowley/go-metrics"
  "log"
  "os"
)


type stats struct {
	peepsAlive metrics.Gauge
	peepsDead  metrics.Gauge
	ages       metrics.Histogram
}

func newStats() *stats {
	r := metrics.NewRegistry()

	stats := &stats{
		peepsAlive: metrics.NewGauge(),
		peepsDead:  metrics.NewGauge(),
		ages:       metrics.NewHistogram(metrics.NewUniformSample(1028)),
	}

	r.Register("peeps_alive", stats.peepsAlive)
	r.Register("peeps_dead", stats.peepsDead)
	r.Register("ages", stats.ages)

	go metrics.Log(metrics.DefaultRegistry, time.Second*1, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

	//go influxdb.Influxdb(r, time.Second*1, &influxdb.Config{
	//	Host:     "127.0.0.1:8086",
	//	Database: "world",
	//	Username: "world",
	//	Password: "peepsrule",
	//})
	return stats

}
