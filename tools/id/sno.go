// Copyright 2023 The olive Authors
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

package id

import (
	"context"

	json "github.com/json-iterator/go"
	"github.com/muyo/sno"
	"github.com/olive-io/bpmn/tracing"
)

type Sno struct{}

func GetSno() *Sno {
	sno_ := &Sno{}
	return sno_
}

type SnoGenerator struct {
	*sno.Generator
	tracer tracing.ITracer
}

func (g *Sno) NewIdGenerator(ctx context.Context, tracer tracing.ITracer) (result IGenerator, err error) {
	return g.RestoreIdGenerator(ctx, []byte{}, tracer)
}

func (g *Sno) RestoreIdGenerator(ctx context.Context, bytes []byte, tracer tracing.ITracer) (result IGenerator, err error) {
	var snapshot *sno.GeneratorSnapshot
	if len(bytes) > 0 {
		snapshot = new(sno.GeneratorSnapshot)
		err = json.Unmarshal(bytes, snapshot)
		if err != nil {
			return
		}
	}
	sequenceOverflowNotificationChannel := make(chan *sno.SequenceOverflowNotification)
	go func(ctx context.Context) {
		for {
			select {
			case notification := <-sequenceOverflowNotificationChannel:
				tracer.Trace(tracing.WarningTrace{Warning: notification})
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	var generator *sno.Generator
	generator, err = sno.NewGenerator(snapshot, sequenceOverflowNotificationChannel)
	if err != nil {
		return
	}
	result = &SnoGenerator{Generator: generator, tracer: tracer}
	return
}

func (g *SnoGenerator) Snapshot() (result []byte, err error) {
	result, err = json.Marshal(g.Generator.Snapshot())
	return
}

type SnoId struct {
	sno.ID
}

func (g *SnoGenerator) New() Id {
	return &SnoId{ID: g.Generator.New(0)}
}

func (id *SnoId) String() string {
	return id.ID.String()
}

func (id *SnoId) Bytes() []byte {
	return id.ID.Bytes()
}

func init() {
	DefaultIdGeneratorBuilder = GetSno()
}
