# 进展记录 PROGRESS

> **AI 每次有实质性进展都要更新本文件**（见 [AGENTS.md](../AGENTS.md) 第 2 节）。
> 倒序排列：最新的在最上面。

## 当前状态

- 分支模型已建立：`master` 为上游纯净镜像，`dev` 为开发分支。
- **项目目标：从中国版适配为日本版**（原项目为中国区特斯拉）。
- 已完成中国特定耦合点全量调研，改造清单见 [JP-ADAPTATION.md](JP-ADAPTATION.md)。
- **已定选型**：地图=Google Maps；兼容=env 开关双区(CN/JP)；语言=vue-i18n 日/英。AI 模型待定。
- 尚未开始代码改动；下一步按 JP-ADAPTATION.md P0 起步（建议先做后端 region 配置开关）。

## 进展日志

### 2026-06-26（feat/coords-wgs84 分支）
- **P0 坐标系区域化**（JP-ADAPTATION #6, #7）：
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
