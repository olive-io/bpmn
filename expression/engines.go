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

package expression

import (
	"context"
	"sync"
)

var enginesLock sync.RWMutex
var enginesMap = make(map[string]func(ctx context.Context) IEngine)

func RegisterEngine(url string, engine func(ctx context.Context) IEngine) {
	enginesLock.Lock()
	defer enginesLock.Unlock()
	enginesMap[url] = engine
}

func GetEngine(ctx context.Context, url string) (engine IEngine) {
	enginesLock.RLock()
	defer enginesLock.RUnlock()
	if engineConstructor, ok := enginesMap[url]; ok {
		engine = engineConstructor(ctx)
	} else {
		engine = enginesMap["http://www.w3.org/1999/XPath"](ctx)
	}
	return
}
