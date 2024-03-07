// Copyright 2023 The bpmn Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clock

import (
	"context"
	"errors"
	"math"
	"time"

	"golang.org/x/sys/unix"
)

func changeMonitor(ctx context.Context, changes chan time.Time) (err error) {
	fd, err := unix.TimerfdCreate(unix.CLOCK_REALTIME, 0)
	if err != nil {
		return
	}

	err = unix.TimerfdSettime(fd, unix.TFD_TIMER_ABSTIME|unix.TFD_TIMER_CANCEL_ON_SET,
		&unix.ItimerSpec{
			Value: unix.Timespec{
				Sec:  math.MaxInt64,
				Nsec: 999_999_999,
			},
		}, nil)
	if err != nil {
		_ = unix.Close(fd)
		return
	}

	err = unix.SetNonblock(fd, true)
	if err != nil {
		_ = unix.Close(fd)
		return
	}

	go func(ctx context.Context) {
		defer func(fd int) {
			_ = unix.Close(fd)
		}(fd)

		buf := make([]byte, 8)
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second * 1):
				_, err := unix.Read(fd, buf)
				if err != nil {
					if errors.Is(err, unix.ECANCELED) {
						changes <- time.Now()
					}
				}
			}
		}
	}(ctx)
	return
}
