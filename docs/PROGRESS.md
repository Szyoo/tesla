# 进展记录 PROGRESS

> **AI 每次有实质性进展都要更新本文件**（见 [AGENTS.md](../AGENTS.md) 第 2 节）。
> 倒序排列：最新的在最上面。

## 架构决策（2026-06-26）：闭源 SwiftUI App + 复用 AGPL 后端

- 发布平台 = **iOS App**。前端用**新建的原生 SwiftUI App（闭源）**替代 UniApp，调后端 API。
- 后端 = **继续复用现有 Go 项目（AGPL）+ 日本化**，不重写。理由：VCP/遥测/状态机是作者踩坑沉淀，重写=几个月重复劳动 + 净室法律风险。
- **闭源诉求如何满足**：用户的私有业务/UI 逻辑放在 SwiftUI App（独立程序，不受 AGPL 约束）；如后端也有私有逻辑，抽成独立闭源服务经 API 调用。AGPL 只约束「被修改的那个程序」本身，不传染独立的 API 客户端。
- 后端日本化（feat/jp-localization 集群）已完成所有不依赖前端的项；UniApp 前端的 i18n/地图（#18-22, #8-13）**作废**，由 SwiftUI App 重做。

## 当前状态

- **后端日本化集群（feat/jp-localization）完成并合并回 dev**：端点区域化、坐标 WGS-84、VIN 电池容量、时区、货币注释、AI Provider 可切换、Accept-Language/AI 报告语言。均编译验证通过。
- 下一阶段：新建闭源 **SwiftUI iOS App** 作为前端。
- 分支模型已建立：`master` 为上游纯净镜像，`dev` 为开发分支。
- **项目目标：从中国版适配为日本版**（原项目为中国区特斯拉）。
- 已完成中国特定耦合点全量调研，改造清单见 [JP-ADAPTATION.md](JP-ADAPTATION.md)。
- **已定选型**：地图=Google Maps；兼容=env 开关双区(CN/JP)；语言=vue-i18n 日/英。AI 模型待定。
- 尚未开始代码改动；下一步按 JP-ADAPTATION.md P0 起步（建议先做后端 region 配置开关）。

## 进展日志

### 2026-06-26（feat/jp-localization 集群分支）
- **AI Provider 可切换**（JP-ADAPTATION #24，设计见 docs/AI-PROVIDERS.md）：
  - 重构 `internal/ai/client.go` 为 `Provider` 接口 + `Chat()` 按 `AI_PROVIDER` 分发；`Chat()` 签名不变，handler 4 处调用零改动。
  - `provider_openai.go`：OpenAI 兼容实现（重构自原逻辑，覆盖智谱/OpenAI/DeepSeek 等）。
  - `provider_anthropic.go`：Claude 官方 anthropic-sdk-go v1.52.0（`client.Messages.New`），响应映射进通用 ChatResponse。
  - config 新增 `AI_PROVIDER`/`ANTHROPIC_API_KEY`/`ANTHROPIC_MODEL`（默认 openai-compat，模型 claude-haiku-4-5）。
  - 苹果端侧 AI（仅 iOS）已做架构设计、预留路径，**未实现**，等发布 iOS App 再做。
  - `go build` + `go vet` + `go mod tidy` 通过。
- **货币单位后端注释**（JP-ADAPTATION #19）：9 处「元」→「円」（models/tesla.go、charging/trip tracker、routes.go），均为注释/字段说明，无逻辑改动。go build 通过。
- **VIN 电池容量映射**（JP-ADAPTATION #15）：
  - 新增 `internal/battery/battery.go`，把原先在 charging/trip 两包**完全重复**的 `getBatteryCapacity` 合并为 `battery.CapacityByVIN`。
  - 认知更正：VIN 前缀按生产地/车型估算，**非中国专属**（LRW=上海产，含出口日本车型），日本沿用即可；仅修正误导性注释。
  - 顺手修复边界 bug：原 `len(vin)<4` 才读 `vin[0:3]`，改为 `<3`。
  - 已 `go build ./...` + `go vet` 通过。
- **时区**（JP-ADAPTATION #17）：`cmd/batch_analyze/main.go` 默认 `Asia/Tokyo`，保留 `TZ` 环境变量覆盖。已编译通过。
- **地图服务（#8-13）暂缓**：调研发现前端用 UniApp 内置 `<map>` 组件（绑定中国地图服务），换 Google Maps 需抛弃该组件改用 Google SDK，且实现方式取决于发布平台（H5/App/小程序）。待用户确定发布平台后再做。
- **AI 模型（#24-25）暂缓**：等用户确定 LLM 选型（Claude/GPT/其他）。

- **P0 坐标系区域化**（JP-ADAPTATION #6, #7，集群分支前称 feat/coords-wgs84）：
  - `internal/geo/geocode.go`：新增 `LocalizeCoords(lat,lng)`——仅 `REGION=cn` 时做 GCJ-02 偏移，日本及其他区域返回原始 WGS-84。`WGS84ToGCJ02` 原函数保留供 cn 用。
  - 3 个调用点（`telemetry/receiver.go` ×2、`fleet/client.go` ×1）改为调 `LocalizeCoords`。
  - 修复真实隐患：旧 `outOfChina` 矩形（经度 72~137.83）会把**西日本**（九州/冲绳/四国等，经度 < 137.83）误判为中国并施加坐标偏移；区域开关彻底规避。
  - 已 `go build ./...` 编译通过。

### 2026-06-26（feat/region-switch 分支）
- **P0 后端 region 开关 + Tesla 端点区域化**（JP-ADAPTATION #1-5, #14）：
  - `config.go`：新增 `REGION` 环境变量（jp/na/eu/cn，默认 jp）与 `teslaRegions` 端点表，按区域提供 OAuth/Fleet/audience/配对页默认端点；显式 `TESLA_*_URL` 仍优先覆盖。`Config` 加 `Region` 字段，`TeslaConfig` 加 `PairingBase`。
  - `internal/tesla/oauth.go`：虚拟钥匙配对 URL 由写死的 `tesla.cn/_ak/` 改为按区域的 `cfg.Tesla.PairingBase`。
  - `config.go` 修复 `ensureHTTPS()`：localhost/127.0.0.1 豁免强制 https（否则 localhost 调试时 redirect_uri 会与 Tesla 门户登记值不符，OAuth 报 mismatch）。
  - `.env.example`：文档化 REGION，回调示例改为 localhost，加端点覆盖说明。
  - ⚠️ **未本地编译验证**：开发机未安装 Go，改动经人工审阅。需在有 Go 环境处 `go build ./...` 复验。

### 2026-06-26（dev 分支）
- 调研中国区耦合点并制定日本版适配计划：
  - 全仓扫描出所有「中国特定」代码/配置，分 P0/P1/P2 共 25 项，记入 `docs/JP-ADAPTATION.md`。
  - **纠正关键误区**：Tesla Fleet API 无日本专属域名，日本属 APAC，走北美区 `na` 端点（`auth.tesla.com` / `fleet-api.prd.na.vn.cloud.tesla.com`），**不是** `*.tesla.jp`。
- 初始化 fork 工作环境：
  - 确认真实代码在 `master` 分支（本地 `main` 仅含 LICENSE）。
  - 从 `origin/master` 切出 `dev` 分支作为主开发分支，并解除其对 master 的 upstream 跟踪（避免误 push 到 master）。
  - 新增 `AGENTS.md`（AI 协作规范）和本进展记录文件。

## 问题与待办

### 待办
- [ ] **用户确认日本版选型**：地图服务 / AI 模型 / 语言策略 / 是否保留中国版兼容（见 JP-ADAPTATION.md）。
- [ ] 按 JP-ADAPTATION.md 优先级从 P0 开始改造。
- [ ] 首次将 `dev` push 到 `origin/dev`（设置远程跟踪）。

### 已知问题 / 坑
- `dev` 分支当前**没有**设置 upstream 远程跟踪——首次 push 用 `git push -u origin dev`，**不要**误推到 master。
- `master` 分支严禁直接修改；同步上游见 AGENTS.md 第 1 节流程。

## 决策记录（为什么这么做）
- **为什么 master 保持纯净**：方便随时 `git merge --ff-only upstream/master` 拉取原作者更新，且让自己的改动与上游清晰分离、降低冲突。
- **为什么改动放 dev 而非 master**：保持上游镜像干净；AGENTS.md/PROGRESS.md 等本 fork 专属文件不污染镜像。
