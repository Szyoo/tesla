# AI 提供方架构设计（可切换 + 苹果端侧）

> 目标：AI 数据分析（行程/充电/车辆）的大语言模型**可自由切换**，并为 **iOS 苹果端侧 AI** 预留路径。
> 现状：后端 `internal/ai/client.go` 写死走 OpenAI 兼容接口、默认智谱 GLM（`open.bigmodel.cn`）。

## 一、两层架构（关键认知）

AI 推理可以在**两个完全不同的位置**发生，它们物理隔离、不能混在一个开关里：

```
┌─ 云端 AI（后端调用）──────────────┐     ┌─ 端侧 AI（仅 iOS App）────────────┐
│ Go 后端把数据拼 prompt，HTTP 调用  │     │ iPhone 上的 Swift 原生代码调用      │
│ 云端 LLM，存库、推前端              │     │ Apple FoundationModels(iOS 26+)，  │
│ → 全平台通用                       │     │ 数据不出设备                        │
│ Zhipu / OpenAI / Claude / DeepSeek │     │ → 仅 iOS App 且设备支持时可用       │
└────────────────────────────────────┘     └────────────────────────────────────┘
        后端碰不到端侧模型；后端在云上，调不到用户手机里的 Apple 模型。
```

## 二、后端：Provider 可切换（本期实现）

用 `AI_PROVIDER` 环境变量选择实现，统一接口，`Chat()` 对上层签名不变（`handler.go` 4 处调用零改动）。

```
AI_PROVIDER=openai-compat   # 默认。覆盖 Zhipu / OpenAI / DeepSeek / Moonshot 等 OpenAI 兼容服务
AI_PROVIDER=anthropic       # Claude（官方 Anthropic API，非兼容 shim）
```

### 接口

```go
// internal/ai/provider.go
type Provider interface {
    Chat(systemPrompt, userPrompt string) (*ChatResponse, error)
}
```

- `ChatResponse` 复用现有结构（含 `Choices[0].Message.Content`、`Usage.PromptTokens/CompletionTokens`），保证 `handler.go` 契约不变。
- `Chat()` 顶层函数按 `cfg.AI.Provider` 选 Provider 实例后转发。

### 两个云端 Provider

| Provider | 实现 | 适用 |
|---|---|---|
| `openaiCompat` | 现有 `net/http` + `/chat/completions`（重构自当前 client.go） | Zhipu GLM、OpenAI GPT、DeepSeek、Moonshot 等 |
| `anthropic` | **官方 anthropic-sdk-go**（`client.Messages.New`），把响应映射进 `ChatResponse` | Claude（Opus/Sonnet/Haiku） |

> 为何 Claude 不走 OpenAI 兼容 shim：官方 SDK 是 Anthropic 推荐做法，能正确处理 thinking、stop_reason、token 统计等差异。模型 ID 用准确字符串（如 `claude-opus-4-8` / `claude-haiku-4-5`，便宜的 Haiku 适合做数据分析摘要、控成本）。

### 配置（.env）

```
AI_PROVIDER=openai-compat
# openai-compat:
AI_BASE_URL=https://open.bigmodel.cn/api/paas/v4   # 或 OpenAI/DeepSeek 等
AI_MODEL=glm-4-flash
AI_API_KEY=...
# anthropic（AI_PROVIDER=anthropic 时）:
ANTHROPIC_API_KEY=...
ANTHROPIC_MODEL=claude-haiku-4-5
```

## 三、苹果端侧 AI（仅 iOS，未实现，预留设计）

**约束**：只能在打包成 iOS App、且设备支持 Apple Intelligence 时，由原生 Swift 通过 `FoundationModels` 框架调用。H5/Android/小程序用不了；Go 后端调不到。

### 设计路径（待发布 iOS App 时实现）

```
前端判断：iOS App + 设备支持端侧 AI ?
  ├─ 是 → 走「端侧分析」：UniApp 调原生 Swift 插件(FoundationModels) → 本地出分析文字
  │         → 可选：把结果 POST 回后端存库（复用现有 AIAnalysis 表/接口）
  └─ 否 → 走「云端分析」：调后端现有 AI 接口（Provider 见上）
```

- 端侧需要一个 **UniApp 原生插件**（Swift），类似现有 tencent-navi 的形态。
- 前端用一个统一的「AI 分析服务」封装，内部决定端侧 vs 云端，对页面透明。
- 后端只需新增/复用一个「保存外部分析结果」的接口（端侧产出回传存档）。

## 四、实现顺序

1. **本期**：后端 Provider 抽象 + openai-compat（重构现有）+ anthropic（官方 SDK）。纯后端、可编译验证。
2. **发布 iOS App 时**：端侧 Swift 插件 + 前端 AI 分析服务封装 + 后端结果回传接口。

## 五、已确认并实现（2026-06-26）

- ✅ 本期实现全部：`Provider` 接口 + `openaiCompat`（重构自原 client.go）+ `anthropic`（官方 anthropic-sdk-go v1.52.0）。
- ✅ 默认 `AI_PROVIDER=openai-compat`，保持现有配置（改动最小）。
- ✅ `Chat()` 签名不变，handler 4 处调用零改动；`go build` + `go vet` 通过。
- 文件：`internal/ai/client.go`（接口+分发）、`provider_openai.go`、`provider_anthropic.go`。
- 端侧 AI（第三节）仍未实现，等发布 iOS App 时再做。
