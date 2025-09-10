
**纯 Go 实现的轻量级业务流程建模与标记(BPMN) 2.0 工作流引擎**

[![Go Reference](https://pkg.go.dev/badge/github.com/olive-io/bpmn.svg)](https://pkg.go.dev/github.com/olive-io/bpmn)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache-blue.svg)](LICENSE.md)
[![Build Status](https://github.com/olive-io/bpmn/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/olive-io/bpmn/actions/workflows/main.yml?query=branch%3Amain)
[![Last Commit](https://img.shields.io/github/last-commit/olive-io/bpmn)](https://github.com/olive-io/bpmn/commits/main)

---

## 简介

`github.com/olive-io/bpmn` 是一个以 Go 语言实现的轻量级**业务流程建模与标记(BPMN) 2.0** 工作流引擎，专为在 Go 应用中建模和执行业务流程而设计。  
它支持核心 BPMN 2.0 元素，包括活动(用户任务、服务任务、脚本任务)、事件(开始事件、结束事件、中间捕获事件)、网关(排他网关、包容网关、并行网关、基于事件网关)、子流程、顺序流以及可扩展的流程属性。

---

## 特性亮点

- **标准 BPMN 2.0 合规性** – 使用标准化业务流程建模与标记元素构建流程模型，完全符合 OMG 规范。
- **轻量级与可嵌入** – 最小化依赖；易于嵌入业务系统，零外部服务需求。
- **完整活动支持** – 用户任务、服务任务、脚本任务、手动任务、业务规则任务以及自定义活动实现。
- **全面流程控制** – 子流程、并行网关、排他网关、包容网关以及基于事件网关，实现复杂决策逻辑。
- **事件驱动处理** – 开始事件、结束事件、中间捕获事件、定时器事件和边界事件。
- **可扩展流程属性** – 自定义属性和数据对象，满足多样化业务流程需求。
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
			select {
			case trace := <-traces:
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
			default:

			}
		}
	}()
	ins.WaitUntilComplete(ctx)

	pros := ins.Locator().CloneVariables()
	log.Printf("%#v", pros)
}

```

### 更多实例
- [基础任务](https://github.com/olive-io/bpmn/tree/main/examples/basic): 简单活动执行示例
- [用户任务](https://github.com/olive-io/bpmn/tree/main/examples/user_task): 实现用户任务活动
- [网关](https://github.com/olive-io/bpmn/tree/main/examples/gateway): 网关流程控制示例
- [网关表达式](https://github.com/olive-io/bpmn/tree/main/examples/gateway_expr): 包容网关与表达式计算
- [自定义属性](https://github.com/olive-io/bpmn/tree/main/examples/properties): 流程特定数据属性和自定义参数
- [捕获事件](https://github.com/olive-io/bpmn/tree/main/examples/catch_event): 中间捕获事件和边界事件示例
- [抛出事件](https://github.com/olive-io/bpmn/tree/main/examples/collaboration): 中间抛出事件和泳池示例
- [子流程](https://github.com/olive-io/bpmn/tree/main/examples/subprocess): 嵌入式子流程执行
- [多通道](https://github.com/olive-io/bpmn/tree/main/examples/multiprocess): 多通道流程的示例

## 许可协议

本项目采用 Apache-2.0 许可证，欢迎商用与二开