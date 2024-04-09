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
