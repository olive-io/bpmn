package expression

import (
	"context"
	"sync"
)

var enginesLock sync.RWMutex
var enginesMap = make(map[string]func(ctx context.Context) Engine)

func RegisterEngine(url string, engine func(ctx context.Context) Engine) {
	enginesLock.Lock()
	defer enginesLock.Unlock()
	enginesMap[url] = engine
}

func GetEngine(ctx context.Context, url string) (engine Engine) {
	enginesLock.RLock()
	defer enginesLock.RUnlock()
	if engineConstructor, ok := enginesMap[url]; ok {
		engine = engineConstructor(ctx)
	} else {
		engine = enginesMap["http://www.w3.org/1999/XPath"](ctx)
	}
	return
}
