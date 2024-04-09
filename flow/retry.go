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

package flow

// Retry describes the retry handler for Flow
type Retry struct {
	limit    int32
	attempts int32
}

func (r *Retry) IsContinue() bool {
	if r.limit == -1 {
		return true
	}
	return r.limit > r.attempts
}

func (r *Retry) Reset(retries int32) {
	r.limit = retries
}

// Step Attempts += 1
func (r *Retry) Step() {
	r.attempts += 1
}
