/*
Copyright 2026 The bpmn Authors

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

package id

import (
	"encoding/json"
	"strconv"
	"sync/atomic"
	"time"
)

type fallbackId struct {
	value string
}

func (id fallbackId) Bytes() []byte {
	return []byte(id.value)
}

func (id fallbackId) String() string {
	return id.value
}

type fallbackGenerator struct {
	prefix  string
	counter uint64
}

func NewFallbackGenerator() IGenerator {
	return &fallbackGenerator{
		prefix: strconv.FormatInt(time.Now().UnixNano(), 36),
	}
}

func (g *fallbackGenerator) Snapshot() ([]byte, error) {
	state := map[string]any{
		"prefix":  g.prefix,
		"counter": atomic.LoadUint64(&g.counter),
	}
	return json.Marshal(state)
}

func (g *fallbackGenerator) New() Id {
	n := atomic.AddUint64(&g.counter, 1)
	return fallbackId{value: "fallback-" + g.prefix + "-" + strconv.FormatUint(n, 10)}
}
