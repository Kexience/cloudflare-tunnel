# AGENT.md - cloudflared-tunnel-go 后端规范

## 环境要求
- 依赖服务：SQLite

## 核心命令
| 命令 | 说明 |
|------|------|
| `go run ./cmd/web` | 启动 web 服务 |
| `go test ./...` | 运行测试 |
| `go run -mod=mod entgo.io/ent/cmd/ent generate ./ent/schema` | 修改 schema 后必须执行 |
| `go build -o dist/web ./cmd/web` | 编译到 dist/ |

## 架构分层
```
cmd/web                 # 入口
internal/               # config, entry, infra, module
ent/schema/             # 唯一数据源（修改后需 ent generate）
ent/*.go                # 生成代码（禁止手动修改）
```

| 层级 | 职责 | 错误处理 |
|------|------|----------|
| ui/api/ctrl | HTTP 绑定、参数校验 | pkg/core 封装响应 |
| svc | 业务逻辑 | pkg/core/errno |
| repo | 数据访问 | 查不到返回 nil, nil |

## 关键约束
- ORM：必须 Ent，禁止手写 SQL
- HTTP：仅 GET/POST，统一返回 200 + JSON body
- DI：Uber FX
- ID 生成：`xid.New().String()`
- Schema 字段需 Comment + `entsql.Annotation`

## 测试规范
- 修改 `**/svc/**.go` 或 `pkg/**/**.go` 必须更新 `*_test.go`
- 测试包名与被测文件同包（如 `svc`）

## 语言与协作
- 所有文本使用中文（日志、错误、注释、commit message）
- 构建产物输出到 `dist/`，禁止提交
- 未经确认不得 git commit/push
