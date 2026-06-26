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

### 2026-06-26
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
