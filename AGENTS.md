# AGENTS.md — AI 协作规范

> 本文件给所有在本仓库工作的 AI 助手（Claude Code / Cursor / 等）阅读。
> 人类维护者：Szyoo。开始任何任务前**先读本文件和 [docs/PROGRESS.md](docs/PROGRESS.md)**。

## 0. 这是什么仓库

特斯拉中国区车联网平台，**fork 自** [798491-collab/tesla](https://github.com/798491-collab/tesla)（AGPL-3.0）。
本 fork 在原项目基础上做二次开发，**目标是不影响原始项目**。

- 后端 `tesla-server/`：Go（Gin + GORM/MySQL + Redis），Tesla Fleet API OAuth、VCP 签名命令代理、车辆状态机、WebSocket 实时推送、遥测接收解析、行程/充电/AI 分析。
- 前端 `tesla-app/`：UniApp + Vue3（可编译 H5 / App / 小程序），Pinia 状态管理，three.js 车模，腾讯地图。

更详细的架构见 [README.md](README.md) 和 [项目说明.md](项目说明.md)。

## 1. 分支模型（重要，别搞错）

| 分支 | 用途 | 规则 |
|------|------|------|
| `master` | **原作者代码的纯净镜像** | 只用来同步 `upstream/master`，**禁止**在上面直接改 |
| `main` | upstream 的占位分支（几乎只有 LICENSE） | 忽略，不用 |
| `dev` | **我们的主开发分支** | 所有改动都在这里或从这里切出的分支 |

remote 约定：
- `origin` = `Szyoo/tesla`（你的 fork，**可以 push**）
- `upstream` = `798491-collab/tesla`（原作者，**只 pull，永不 push**）

**绝对规则：永远不要 push 到 `upstream`。** 这样原始项目天然不受影响。

### 同步原作者更新的标准流程
```bash
git fetch upstream
git checkout master
git merge --ff-only upstream/master   # master 必须保持纯净，只快进
git push origin master                # 可选：更新你 fork 的 master
git checkout dev
git merge master                      # 把上游更新并入 dev，在这里解决冲突
```

## 2. AI 必须遵守的工作流

1. **动手前**：读本文件 + `docs/PROGRESS.md`，确认当前在 `dev`（或其子分支），不在 `master`。
2. **干活时**：遵循下方代码规范；改动尽量小而聚焦。
3. **完成一个阶段后（强制）**：更新 `docs/PROGRESS.md`——
   - 在「进展日志」追加一条（日期 + 做了什么 + 影响的文件）。
   - 遇到的问题 / 坑 / 待办，记入「问题与待办」区。
   - 不要等用户提醒，每次有实质性进展就更新。
4. **提交**：commit message 用中文或英文均可，但要清晰；遵循已有风格（见 `git log`，多用 `feat:` / `fix:` 前缀）。除非用户明确要求，**不要自动 push**。

## 3. 代码规范

- **不要碰密钥**：`.env`、`manifest.json`、`*.pem` 已被 gitignore，配置改 `.env.example` / `manifest.json.example`。
- **匹配现有风格**：写代码前先看周围文件的命名、缩进、注释密度，保持一致。
- 后端 Go：包结构在 `tesla-server/internal/<domain>/`，新功能按领域分包。
- 前端：页面在 `tesla-app/pages/`，组件在 `components/`，状态在 `store/`（Pinia），API 封装在 `api/`。
- 保留原作者版权声明（AGPL-3.0 要求）。

## 4. 构建与运行

后端：
```bash
cd tesla-server
cp .env.example .env   # 填配置
go mod tidy
go build -o tesla-server cmd/main.go && ./tesla-server
```

前端：
```bash
cd tesla-app
cp .env.example .env && cp manifest.json.example manifest.json
npm install
npm run dev:h5
```

## 5. AGPL-3.0 合规提醒

- 修改后的代码必须开源；部署为网络服务也必须公开源码。
- 必须保留原作者版权声明，不得用原作者商标推广衍生产品。
