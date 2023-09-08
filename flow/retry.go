// Copyright 2023 Lack (xingyys@gmail.com).
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
