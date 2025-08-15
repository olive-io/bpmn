//go:build !linux

/*
Copyright 2023 The bpmn Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
