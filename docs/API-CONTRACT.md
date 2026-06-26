# 后端 API 契约（供 SwiftUI App 对接）

> 本文件汇总 Go 后端对前端暴露的 REST + WebSocket 接口，供独立的闭源 SwiftUI App 对接。
> 来源：`tesla-server/routes/routes.go`。基址默认 `http://localhost:8080`（部署后换正式 HTTPS 域名）。

## 认证

- 注册/登录返回 JWT；后续请求带 `Authorization: Bearer <token>`。
- `POST /api/register`、`POST /api/login`、`POST /api/refresh-token`（有 RateLimitAuth 限流）。
- 受保护接口走 `JWTAuth` 中间件。

## 公钥托管（VCP 配对用）

- `GET /.well-known/appspecific/com.tesla.3p.public-key.pem` —— Tesla 虚拟钥匙配对要求，由后端在正式域名根托管。

## 用户

- `POST /api/logout`
- `GET /api/user/info`
- `POST /api/user/change_password`
- `POST /api/user/update`

## Tesla 账号 / 车辆绑定

- `GET /api/tesla/auth` —— 获取 OAuth 授权 URL
- `GET /api/tesla/callback` —— OAuth 回调
- `GET /api/tesla/auth_data`
- `POST /api/tesla/partner/register`、`GET /api/tesla/partner/check-public-key`、`GET /api/tesla/partner/check-hosting`
- `POST /api/tesla/bind`、`DELETE /api/tesla/unbind/:vin`
- `GET /api/tesla/vehicles` —— 用户车辆列表
- `GET /api/tesla/vehicle/:vin/detail`
- `GET /api/tesla/vehicle/:vin/fleet-status`、`/fleet-telemetry-errors`
- `GET /api/tesla/vehicle/:vin/pairing-url` —— 虚拟钥匙配对 URL
- `POST /api/tesla/refresh-vehicle-info`

## 车辆状态

- `GET /api/vehicle/:vin/state` —— 当前状态
- `POST /api/vehicle/:vin/refresh` —— 主动刷新
- `POST /api/vehicle/:vin/wake` —— 唤醒
- `GET /api/vehicle/:vin/data` —— 详细数据
- `GET /api/vehicle/:vin/tracks` —— 轨迹

## 行程

- `GET /api/trip/:vin/logs`、`/stats`、`/monthly-list`、`/monthly-stats`
- `GET /api/trip/:vin/points/:tripId` —— 单次行程轨迹点

## 充电

- `GET /api/charging/:vin/logs`、`/stats`、`/monthly-list`、`/monthly-stats`
- `POST /api/charging/log/:id/price` —— 补充电价（慢充 price_per_kwh / 快充 total_cost，单位円）

## AI 分析

- `POST /api/ai/trip/:vin/:refId`、`/ai/charging/:vin/:refId`、`/ai/vehicle/:vin` —— 触发分析（异步）
- `GET /api/ai/trip/:vin/:refId`、`/ai/charging/:vin/:refId`、`/ai/vehicle/:vin` —— 取结果
- `GET /api/ai/history/:vin`、`/ai/list/:vin`、`/ai/latest/:vin/:type`
- 完成后经 WebSocket 推送 `analysis_complete` 事件（含 type/ref_id/status）。

## VCP 车辆控制（需虚拟钥匙配对）

车门：`door_lock`、`door_unlock` · 空调：`auto_conditioning_start`、`auto_conditioning_stop`、`set_temps`、`remote_seat_heater`、`remote_steering_wheel_heater` · 提示：`honk_horn`、`flash_lights` · 行李厢：`actuate_trunk`、`actuate_frunk` · 哨兵：`set_sentry_mode` · 充电：`charge_start`、`charge_stop`、`set_charge_limit`、`charge_port_door_open`、`charge_port_door_close` · 媒体：`media_toggle_playback`、`media_next_track`、`media_prev_track`、`media_next_fav`、`media_prev_fav`、`media_volume_up`、`media_volume_down`、`adjust_volume`

全部为 `POST /api/vcp/<command>`。`GET /api/vcp/commands` 列出可用命令。

⚠️ VCP 命令需正式域名 + 公钥配对，localhost 不可用（仅读数据可调试）。

## WebSocket 实时推送

- `GET /api/ws`、`GET /api/ws/vin/:vin` —— 实时车辆状态、遥测、analysis_complete 等事件。
- App 应优先用 WebSocket 获取实时数据（省 Fleet API 调用，控成本），少用轮询。

## 遥测

- `POST /api/telemetry/config` —— 配置遥测字段。
