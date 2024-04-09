//go:build !linux

/*
   Copyright 2023 The bpmn Authors

   This program is offered under a commercial and under the AGPL license.
   For AGPL licensing, see below.

   AGPL licensing:
   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
