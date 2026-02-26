[English](https://github.com/olive-io/bpmn/tree/main/README.md) | 中文

**纯 Go 实现的轻量级 BPMN 2.0 工作流引擎**

[![Go Reference](https://pkg.go.dev/badge/github.com/olive-io/bpmn.svg)](https://pkg.go.dev/github.com/olive-io/bpmn)
[![License: Apache-2.0](https://img.shields.io/badge/license-Apache-blue.svg)](LICENSE.md)
[![Build Status](https://github.com/olive-io/bpmn/actions/workflows/main.yml/badge.svg?branch=main)](https://github.com/olive-io/bpmn/actions/workflows/main.yml?query=branch%3Amain)
[![Last Commit](https://img.shields.io/github/last-commit/olive-io/bpmn)](https://github.com/olive-io/bpmn/commits/main)

---

## 这是什么？

`github.com/olive-io/bpmn` 是一个可嵌入 Go 应用的 BPMN 2.0 运行时。

适合用于：

- 在 Go 服务内执行 BPMN 流程定义。
- 通过用户任务/服务任务/脚本任务接入业务逻辑。
- 通过 trace 事件流做运行时观测和排障。

仓库包含两个模块：

- 运行时模块：`github.com/olive-io/bpmn/v2`
- BPMN schema 模块：`github.com/olive-io/bpmn/schema`

---

## 安装

```bash
go get -u github.com/olive-io/bpmn/schema
go get -u github.com/olive-io/bpmn/v2
```

---

## 快速开始

下面示例会加载 BPMN 文件、启动流程、处理任务 trace，并等待流程结束。

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
	data, err := os.ReadFile("task.bpmn")
	if err != nil {
		log.Fatalf("读取 BPMN 文件失败: %v", err)
	}

	definitions, err := schema.Parse(data)
	if err != nil {
		log.Fatalf("解析 BPMN XML 失败: %v", err)
	}

	engine := bpmn.NewEngine()
	ctx := context.Background()

	proc, err := engine.NewProcess(definitions,
		bpmn.WithVariables(map[string]any{"customer": "alice"}),
	)
	if err != nil {
		log.Fatalf("创建流程失败: %v", err)
	}

	traces := proc.Tracer().Subscribe()
	defer proc.Tracer().Unsubscribe(traces)

	if err := proc.StartAll(ctx); err != nil {
		log.Fatalf("启动流程失败: %v", err)
	}

	go func() {
		for trace := range traces {
			trace = tracing.Unwrap(trace)
			switch t := trace.(type) {
			case bpmn.TaskTrace:
				// 回写任务执行结果。
				t.Do(bpmn.DoWithResults(map[string]any{"approved": true}))
			case bpmn.ErrorTrace:
				log.Printf("流程错误: %v", t.Error)
			}
		}
	}()

	if ok := proc.WaitUntilComplete(ctx); !ok {
		log.Printf("流程被取消")
	}

	log.Printf("流程变量: %#v", proc.Locator().CloneVariables())
}
```

---

## 核心概念

- `Engine`：从 BPMN 定义创建流程实例。
- `Process`：执行单个 executable process。
- `ProcessSet`：执行同一个 definitions 中的多个 executable process。
- `Tracer`：运行时事件流，通过订阅/退订获取 trace。
- `TaskTrace`：任务等待外部输入时的回调入口，可返回结果、数据对象或错误。

---

## 什么时候用 NewProcess / NewProcessSet

当 definitions 中只有一个 executable process 时，使用 `NewProcess`。

当 definitions 中存在多个 executable process 时，使用 `NewProcessSet`。

```go
proc, err := engine.NewProcess(definitions)
// 如果有多个 executable process，会返回错误。

set, err := engine.NewProcessSet(&definitions)
// 同时运行 definitions 中所有 executable process。
```

---

## 可观测性与运行时错误

建议订阅 tracer 事件来观察流程状态：

- `FlowTrace`：顺序流流转。
- `TaskTrace`：任务等待外部处理。
- `ErrorTrace`：运行时发生错误。
- `CeaseFlowTrace` / `CeaseProcessSetTrace`：流程或流程集合执行结束。

ID 生成器降级行为：

- 运行时会优先初始化基于 SNO 的默认 ID 生成器。
- 若初始化失败，会自动降级到本地 fallback 生成器。
- 同时输出 warning trace，便于在监控平台中定位。

---

## BPMN 能力覆盖

当前支持常见服务流程需要的核心 BPMN 元素，包括：

- 任务：User、Service、Script、Manual、Business Rule、Call Activity。
- 事件：Start、End、Intermediate Catch/Throw、Timer 相关流转。
- 网关：Exclusive、Inclusive、Parallel、Event-based。
- 子流程和顺序流。

更多具体用法请参考 examples 与测试用例。

---

## 开发与测试命令

在仓库根目录执行：

```bash
# 运行运行时全部测试
go test -v ./...

# 运行 schema 模块测试
go test -v ./schema

# 静态检查
go vet ./...

# 按名称运行单个测试
go test -v ./... -run 'TestUserTask'

# 在指定包运行单测并禁用缓存
go test -v ./pkg/data -run '^TestContainer$' -count=1
```

---

## 示例

- [quickstart](https://github.com/olive-io/bpmn/tree/main/examples/readme_quickstart)
- [basic](https://github.com/olive-io/bpmn/tree/main/examples/basic)
- [user_task](https://github.com/olive-io/bpmn/tree/main/examples/user_task)
- [gateway](https://github.com/olive-io/bpmn/tree/main/examples/gateway)
- [gateway_expr](https://github.com/olive-io/bpmn/tree/main/examples/gateway_expr)
- [properties](https://github.com/olive-io/bpmn/tree/main/examples/properties)
- [catch_event](https://github.com/olive-io/bpmn/tree/main/examples/catch_event)
- [collaboration](https://github.com/olive-io/bpmn/tree/main/examples/collaboration)
- [subprocess](https://github.com/olive-io/bpmn/tree/main/examples/subprocess)
- [multiprocess](https://github.com/olive-io/bpmn/tree/main/examples/multiprocess)

---

## 参与贡献

欢迎提交 PR。贡献与 agent 协作规范见 `AGENTS.md`。

---

## 许可证

Apache-2.0
