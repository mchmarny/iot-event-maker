package hardware

import (
	"context"
	"sync"
	"time"

	"github.com/mchmarny/eventmaker/pkg/event"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/cpu"
)

// NewCPUMetricProvider creates nee CPUMetricProvider
func NewCPUMetricProvider() *CPUMetricProvider {
	return &CPUMetricProvider{}
}

// CPUMetricProvider is a host CPU metric provider
type CPUMetricProvider struct{}

// Provide provides os process events
func (p *CPUMetricProvider) Provide(ctx context.Context, wg *sync.WaitGroup, src string, d time.Duration, h func(e *event.SimpleEvent)) error {
	defer wg.Done()
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			e, err := getCPUMetric(src, d)
			if err != nil {
				return err
			}
			h(e)
		case <-ctx.Done():
			ticker.Stop()
			return nil
		}
	}
}

func getCPUMetric(src string, d time.Duration) (e *event.SimpleEvent, err error) {
	if src == "" {
		return nil, errors.New("nil source ID")
	}
	mp, err := cpu.Percent(d, false)
	if err != nil {
		return nil, errors.Wrap(err, "error creating cpu utilization metric")
	}
	e = &event.SimpleEvent{
		ID:    uuid.NewV4().String(),
		SrcID: src,
		Time:  time.Now().UTC().Unix(),
		Label: "CPU Utilization",
		Unit:  "percent",
		Value: mp[0],
	}
	return
}