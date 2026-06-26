# 进展记录 PROGRESS

> **AI 每次有实质性进展都要更新本文件**（见 [AGENTS.md](../AGENTS.md) 第 2 节）。
> 倒序排列：最新的在最上面。

## 当前状态

- 分支模型已建立：`master` 为上游纯净镜像，`dev` 为开发分支。
- 尚未开始功能改动；等待第一个开发任务。

## 进展日志

### 2026-06-26
- 初始化 fork 工作环境：
  - 确认真实代码在 `master` 分支（本地 `main` 仅含 LICENSE）。
  - 从 `origin/master` 切出 `dev` 分支作为主开发分支，并解除其对 master 的 upstream 跟踪（避免误 push 到 master）。
  - 新增 `AGENTS.md`（AI 协作规范）和本进展记录文件。

## 问题与待办

### 待办
- [ ] 确定第一个要改的功能 / 模块。
- [ ] 首次将 `dev` push 到 `origin/dev`（设置远程跟踪）。

### 已知问题 / 坑
- `dev` 分支当前**没有**设置 upstream 远程跟踪——首次 push 用 `git push -u origin dev`，**不要**误推到 master。
- `master` 分支严禁直接修改；同步上游见 AGENTS.md 第 1 节流程。

## 决策记录（为什么这么做）
- **为什么 master 保持纯净**：方便随时 `git merge --ff-only upstream/master` 拉取原作者更新，且让自己的改动与上游清晰分离、降低冲突。
- **为什么改动放 dev 而非 master**：保持上游镜像干净；AGENTS.md/PROGRESS.md 等本 fork 专属文件不污染镜像。
