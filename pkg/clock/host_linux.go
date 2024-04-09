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
