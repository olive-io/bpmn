//go:build !linux

package clock

import (
	"context"
	"time"
)

const hostForwardDriftTolerance = 3 * time.Second

func changeMonitor(ctx context.Context, changes chan time.Time) (err error) {
	go func(ctx context.Context) {
		for {
			t := time.Now()
			select {
			case <-ctx.Done():
				return
			case t1 := <-time.After(time.Second * 1):
				if t1.Before(t) {
					// backward drift
					changes <- t1
				} else if t1.Sub(t).Nanoseconds() > hostForwardDriftTolerance.Nanoseconds() {
					// forward drift
					changes <- t1
				}
			}
		}
	}(ctx)
	return
}
