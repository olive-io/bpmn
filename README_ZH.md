# github.com/olive-io/bpmn

**纯 Go 实现的轻量级 bpmn2.0 工作流引擎**

[![Go Reference](https://pkg.go.dev/badge/github.com/olive-io/bpmn.svg)](https://pkg.go.dev/github.com/olive-io/bpmn)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache-blue.svg)](LICENSE.md)

---

## 简介

`github.com/olive-io/bpmn` 是一个以 Go 语言实现、轻量级的 **BPMN 2.0** 工作流引擎，旨在简化在 Go 应用中嵌入业务流程的建模与执行。它支持包括任务类型（User Task、Service Task、Script Task）、事件（Start / End / Catch）、各种网关（Exclusive, Inclusive, Parallel, Event-based）、子流程、流程控制、以及可自定义属性等核心 BPMN 元素。

---

## 特性亮点

- **原生 BPMN 2.0 支持**：直接以标准 BPMN 元素构建流程，符合建模规范。
- **轻量与简洁**：核心依赖少，适合直接嵌入业务系统，无需额外服务部署。
- **支持多任务类型**：包括用户任务、脚本任务、服务任务、自定义任务等。
- **丰富的流程控制**：支持子流程、并行、排他和事件驱动决策等结构。
- **可定制扩展**：支持流程属性扩展，满足业务场景多样性。
- **全面测试覆盖**：示例和模块均附有单元测试，保证执行逻辑可靠。

---

## 开始使用

### 安装

```bash
go get -u github.com/olive-io/bpmn/schema
go get -u github.com/olive-io/bpmn/v2
```

### 快速开始

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/olive-io/bpmn/schema"
	"github.com/olive-io/bpmn/v2"
	"github.com/olive-io/bpmn/v2/pkg/tracing"
)

func main() {

	var err error
	data, err := os.ReadFile("task.bpmn")
	if err != nil {
		log.Fatalf("Can't read bpmn: %v", err)
	}
	definitions, err := schema.Parse(data)
	if err != nil {
		log.Fatalf("XML unmarshalling error: %v", err)
	}

	engine := bpmn.NewEngine()
	options := []bpmn.Option{
		bpmn.WithVariables(map[string]any{
			"c": map[string]string{"name": "cc"},
		}),
		bpmn.WithDataObjects(map[string]any{
			"a": struct{}{},
		}),
	}
	ctx := context.Background()
	ins, err := engine.NewProcess(&definitions, options...)
	if err != nil {
		log.Fatalf("failed to instantiate the process: %s", err)
		return
	}
	traces := ins.Tracer().Subscribe()
	defer ins.Tracer().Unsubscribe(traces)
	err = ins.StartAll(ctx)
	if err != nil {
		log.Fatalf("failed to run the instance: %s", err)
	}
	go func() {
		for {
			var trace tracing.ITrace
			select {
			case trace = <-traces:
			}

			trace = tracing.Unwrap(trace)
			switch trace := trace.(type) {
			case bpmn.FlowTrace:
			case bpmn.TaskTrace:
				trace.Do(bpmn.DoWithResults(
					map[string]any{
						"c": map[string]string{"name": "cc1"},
						"a": 2,
					}),
				)
			case bpmn.ErrorTrace:
				log.Fatalf("%#v", trace)
				return
			case bpmn.CeaseFlowTrace:
				return
			default:
				log.Printf("%#v", trace)
			}
		}
	}()
	ins.WaitUntilComplete(ctx)

	pros := ins.Locator().CloneVariables()
	log.Printf("%#v", pros)
}
```

### 更多实例
- [单任务](https://github.com/olive-io/bpmn/tree/main/examples/basic): 最简单的任务
- [用户任务](https://github.com/olive-io/bpmn/tree/main/examples/user_task): 执行用户任务
- [网关](https://github.com/olive-io/bpmn/tree/main/examples/gateway): 执行带网关实例
- [排他网关](https://github.com/olive-io/bpmn/tree/main/examples/gateway_expr): 带条件的排他网关
- [自定义参数](https://github.com/olive-io/bpmn/tree/main/examples/properties): 支持自定义任务参数
- [子进程](https://github.com/olive-io/bpmn/tree/main/examples/subprocess): 执行子进程

## 许可协议

本项目采用 Apache-2.0 许可证，欢迎商用与二开