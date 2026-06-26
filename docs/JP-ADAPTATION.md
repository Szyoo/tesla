# 日本版适配清单 JP Adaptation

> 本项目 fork 自中国区版本，目标是适配**日本版**特斯拉。
> 本文件记录所有「中国特定 → 日本」的改造点。AI 完成每一项后在表格勾选，并在 [PROGRESS.md](PROGRESS.md) 记一笔。

## ⚠️ 关键认知：Tesla Fleet API 没有日本专属区域

Tesla Fleet API 只有 **3 个区域**：北美/亚太（NA/APAC）、欧洲（EU）、中国（CN）。
**日本属于亚太（APAC），走北美区 `na` 端点**，不存在 `*.tesla.jp` 域名。

| 配置 | 中国（现状） | 日本 / APAC（正确值） |
|------|------|------|
| Auth/Token 域名 | `auth.tesla.cn` | `auth.tesla.com` |
| Fleet API | `fleet-api.prd.cn.vn.cloud.tesla.cn` | `fleet-api.prd.na.vn.cloud.tesla.com` |
| Audience | 同 Fleet API（cn） | 同 Fleet API（na） |
| 虚拟钥匙配对 | `tesla.cn/_ak/` | `tesla.com/_ak/` |
| 开发者门户 | `developer.tesla.cn` | `developer.tesla.com` |

> 部署前请以 Tesla 官方文档复核一次区域映射（developer.tesla.com/docs/fleet-api）。
> EU 是 `...prd.eu...`，仅作对照，日本不用。

## 改造项（按优先级）

### P0 — 不改无法连接 / 位置错误

| # | 类别 | 文件 | 中国值 → 日本值 | 状态 |
|---|------|------|------|------|
| 1 | OAuth Token URL | `tesla-server/config/config.go:90` | `auth.tesla.cn/...token` → `auth.tesla.com/...token` | ☐ |
| 2 | OAuth Auth URL | `tesla-server/config/config.go:91` | `auth.tesla.cn/...authorize` → `auth.tesla.com/...authorize` | ☐ |
| 3 | Fleet API URL | `tesla-server/config/config.go:92` | `fleet-api...cn...tesla.cn` → `fleet-api.prd.na.vn.cloud.tesla.com` | ☐ |
| 4 | Audience | `tesla-server/config/config.go:93` | 同上 na URL | ☐ |
| 5 | .env 示例 + 注释 | `tesla-server/.env.example:18-23` | `(China)` 默认值 + `developer.tesla.cn` → na 端点 + `developer.tesla.com` | ☐ |
| 6 | 坐标系：停用 GCJ-02 偏移 | `tesla-server/internal/telemetry/receiver.go`、`internal/fleet/client.go` | 移除 `geo.WGS84ToGCJ02()` 调用，日本直接用 WGS-84 | ☐ |
| 7 | 坐标转换函数 | `tesla-server/internal/geo/geocode.go:169-186` | 保留代码但不调用（或加开关），日本无需偏移 | ☐ |
| 8 | 地图反向地理编码 | `tesla-server/internal/geo/geocode.go:51` | 腾讯 `apis.map.qq.com/.../geocoder` → Google/其他 | ☐ |
| 9 | 地图距离查询 | `tesla-server/internal/geo/geocode.go:111` | 腾讯 `apis.map.qq.com/.../distance` → Google/其他 | ☐ |
| 10 | 前端地图 Key | `tesla-app/.env.example:3-5`、`tesla-server/.env.example:40` | `VITE_TENCENT_MAP_KEY` → `VITE_GOOGLE_MAPS_KEY`（或所选服务） | ☐ |
| 11 | 前端地图页面 | `pages/vehicle/location.vue`、`charging/map.vue`、`trip/route.vue`、`dashcam/map.vue`、`dashboard/dashboard.vue`、`dashboard/instrument.vue`、`components/cockpit/MapPanel.vue` | 替换腾讯地图 SDK 调用为新地图服务 | ☐ |
| 12 | 原生导航插件 | `tesla-app/nativeplugins/tencent-navi/` | 腾讯导航 SDK → 移除或换 Google/MapBox | ☐ |
| 13 | manifest 地图配置 | `tesla-app/manifest.json.example:60,73-75,117-130` | TencentMapSDK metadata → 新地图服务 | ☐ |

### P1 — 重要（功能正确性 / 配对）

| # | 类别 | 文件 | 中国值 → 日本值 | 状态 |
|---|------|------|------|------|
| 14 | 虚拟钥匙配对 URL | `tesla-server/internal/tesla/oauth.go:1253` | `tesla.cn/_ak/` → `tesla.com/_ak/` | ☐ |
| 15 | VIN 前缀/电池容量 | `tesla-server/internal/charging/tracker.go:118-128`、`internal/trip/tracker.go` | `LRW`(中国国行) 映射 → 加入日本进口车 VIN 前缀（待查实际值，多为 5YJ/7SA/XP7 进口） | ☐ |
| 16 | Fleet 注释/endpoints | `tesla-server/internal/fleet/client.go` | 移除「中国区不支持 closures_state」等假设，按 APAC 实测调整 | ☐ |
| 17 | 时区 | `tesla-server/cmd/batch_analyze/main.go` | `Asia/Shanghai` → `Asia/Tokyo` | ☐ |

### P2 — 本地化（语言 / 货币 / 单位）

| # | 类别 | 文件 | 中国值 → 日本值 | 状态 |
|---|------|------|------|------|
| 18 | 货币单位（前端） | `pages/trip/trip.vue`、`charging/month.vue` 等 | `元` / `元/kWh` → `円` / `円/kWh`（JPY） | ☐ |
| 19 | 货币单位（后端注释/字段） | `models/tesla.go`、`charging/tracker.go`、`trip/tracker.go`、`routes/routes.go` | 注释「元」→「円」 | ☐ |
| 20 | App 名称/描述 | `tesla-app/manifest.json.example:2,4` | `Tesla中国区...` → 日文/英文 | ☐ |
| 21 | 页面标题 i18n | `tesla-app/pages.json`（多处 navigationBarTitleText） | 中文 → 日文（或接入 vue-i18n，项目已装 `vue-i18n`） | ☐ |
| 22 | 页面内中文文案 | 各 `pages/*.vue` | 中文 → 日文 | ☐ |
| 23 | Accept-Language | `tesla-server/internal/tesla/oauth.go:908` | `zh-CN,zh...` → `ja-JP,ja;q=0.9,en;q=0.8` | ☐ |
| 24 | AI 模型 | `tesla-server/config/config.go:143`、`.env.example:43-45` | Zhipu `glm-4-flash` / `open.bigmodel.cn` → 国际 LLM（Claude/GPT） | ☐ |
| 25 | AI 报告输出语言 | `tesla-server/internal/ai/handler.go` | 中文 key/单位 → 日文 | ☐ |

## 已确认的决策（2026-06-26）

- **地图服务 = Google Maps**。理由：日本特斯拉车机本身即用 Google 地图（特斯拉全球车机除中国外均用 Google 地图数据），符合车机生态与日本用户习惯。影响 #8–13。
- **双区兼容 = 用 env 开关切 CN/JP**。保留中国代码路径，端点/地图/坐标/货币由配置控制。新增配置项 `REGION=jp|cn`（或等价开关），坐标偏移、地图服务、Fleet 端点据此分流。代码改动最小、同步 upstream 冲突最小。影响全局架构。
- **语言 = 接 vue-i18n 做日/英多语言**。文案抽成语言包（`ja` / `en`），默认日文。影响 #18–22, #25。
- **AI 模型**：仍待定（Claude / GPT / 其他国际 LLM）。影响 #24–25。

## 由决策衍生的架构约定

- 后端新增 region 配置（如 `REGION` 环境变量），`config.go` 据此选择 Tesla 端点、是否做 GCJ-02 偏移、默认货币/时区。
- 前端用 vue-i18n 管理文案；地图封装一层抽象，便于 CN(腾讯)/JP(Google) 切换。
- 同步 upstream 时，双区开关让我们的改动以「新增分支逻辑」为主，尽量不删原中国代码 → 降低 merge 冲突。
