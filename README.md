# 航空地勤周转保障平台

面向机场地勤团队的**航班过站保障、资源预约、整批冲突校验、延误重排与任务签收**系统。航班到站后按任务类型/班组/资源可用时段整批生成地勤任务和资源预约；任何冲突整批不保存并逐项说明原因；登记延误分钟后自动重排未签收任务、冻结已签收任务、把被挤出的预约置为待处理，页面刷新后看到同一份结果。

## 快速启动

```bash
cp .env.example .env && docker compose up -d
```

- 前端：<http://localhost:20108>
- 后端健康检查：<http://localhost:21108/health>

首次登录在页面右上角弹窗选择本地账号（已签发 JWT）：

| 账号 | 角色 | 可执行动作 |
|---|---|---|
| dispatcher / supervisor | 地勤调度 / 运行督导 | 到站登记、生成计划、登记延误、关闭归因 |
| resource | 资源管理员 | 资源状态/维护窗口、调整待处理预约 |
| team-bag / team-fuel / team-ramp | 保障班组 | 任务签收、完成、阻塞 |

种子数据（2026-09-22 10:00 基准）内置一条完整演示链路：`CA101`（已到站、未生成计划，首次生成会被整批驳回）、`MU202`（有已签收行李任务）、`CZ303`（已签收清洁任务）。

## 保障闭环演示步骤

1. 以 `dispatcher` 登录，打开 **航班过站 → CA101 → 生成保障任务与预约**：皮带车重叠/离线、加油车维护/占用，**整批不保存**，抽屉逐项列出每条冲突原因。
2. 切换 `resource`：把 `BL-02` 置为可用、`FT-02` 置为可用并登记 10:40–11:15 维护窗口、释放占用 `FT-02` 的 MU202 预约。
3. 回到 `dispatcher` 重新生成：6 个任务、6 条预约一次保存，航班进入 `IN_SERVICE`。
4. 切换 `team-bag` **签收行李任务**（计划时点随即冻结）。
5. `dispatcher` 登记延误 `+30` 分钟：5 个未签收任务整体平移，行李任务保持 10:02 原时点；清洁预约撞上 CZ303 已签收窗口、加油预约撞上 FT-02 维护窗口，两条预约进入 **PENDING 待处理**。
6. 看板 / 资源调度页都能看到同一份待处理列表；`resource` 调整时段：调整到仍冲突的窗口会返回逐项原因并保持待处理，调整到干净窗口则恢复已确认（并同步未签收任务时间）。
7. 刷新任意页面，看板统计、时间轴、日历、待处理列表均来自同一后端状态。

后端有两条自动化测试完整覆盖该链路（`go test ./...`）：
`src/services/closedloop_scenario_test.go`（服务+事务层）与 `src/routes/api_closedloop_test.go`（真实 HTTP/JWT/RBAC 层）。

## 访问地址或 CLI 示例

```bash
# 健康检查
curl http://localhost:21108/health

# 登录拿 JWT
TOKEN=$(curl -s -X POST http://localhost:21108/api/auth/login \
  -H 'Content-Type: application/json' -d '{"username":"dispatcher"}' | jq -r .token)

# 生成计划（冲突时返回 409 + details[]，整批未写入）
curl -X POST http://localhost:21108/api/turnarounds/1/plan \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{}'

# 登记延误并重排
curl -X POST http://localhost:21108/api/turnarounds/1/delays \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"delay_type":"LATE_ARRIVAL","minutes":30,"root_cause":"前站晚到","responsibility_team":"AIRLINE"}'

# 调整待处理预约（仍冲突时 409 + details[]，保持 PENDING）
curl -X POST http://localhost:21108/api/resource-bookings/9/adjust \
  -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  -d '{"start_time":"2026-09-22T11:10:00+08:00","end_time":"2026-09-22T11:25:00+08:00"}'
```

## 本地开发方式

- 前端：`cd frontend && npm install && npm run dev`（Vite 端口 20108，开发期可在 `vite.config.ts` 增加 `/api` 代理）
- 后端：`cd backend && go run ./main.go`，接口统一挂在 `/api`，默认读 `DB_HOST=127.0.0.1 DB_PORT=3306` 等环境变量
- 后端测试：`cd backend && CGO_ENABLED=1 go test ./...`（用内存 SQLite 跑全链路，无需 MySQL）

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | React 18 + TypeScript + Vite + Ant Design + Redux Toolkit + React Router |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0（迁移由 GORM AutoMigrate 执行，`database/init.sql` 为同构 DDL） |
| 认证 | JWT（HS256）+ RBAC（四类角色） |
| 部署 | Docker Compose（db 命名卷 + healthcheck 依赖链） |

## 项目目录结构

```text
frontend/src/
├── api/                  # client.ts 统一 /api 请求与错误信封；按实体拆分
├── stores/               # Redux Toolkit slices（auth/turnaround/task/resource/delay）
├── types/                # 实体与枚举类型
├── constants/            # 枚举、错误码、错误消息、日志模板、冲突原因
├── constructors/         # 表单/空对象构造器
├── components/           # AppShell、LoginModal
│   └── common/           # StatusBadge/TurnaroundTimeline/ResourceCalendar/ConflictBadge/
│                         # ConflictPanel/PendingBookingList/DelayTag/StatCard/TeamTag/TimelineList
├── hooks/                # useTurnaroundProgress/useResourceConflict/usePagination/useRole
├── pages/                # Dashboard/Turnarounds/Tasks/Resources/Delays
├── router/               # routes.tsx + RequireAuth 路由守卫 + LoginGate
├── utils/formatters.ts   # 日期/数字/风险等混合格式化
└── mocks/                # 仅保留 localOnly 标记，禁止第三方 API

backend/src/
├── routes/               # 路由装配 + JWT/RBAC 挂载（含全链路 HTTP 测试）
├── controllers/          # 按实体拆分，单独包装错误
├── services/             # PlanService/DelayService/TaskService/BookingService 等业务核心
├── models/               # GORM 实体
├── repositories/         # 数据访问
├── middlewares/          # auth/rbac/auditLog/errorHandler/rateLimit
├── constants/            # 枚举、错误码、错误消息、日志模板、冲突原因、角色
├── constructors/         # 请求 DTO 与响应视图构造器
├── utils/                # 时间、JWT、错误信封
├── types/                # ConflictItem/PlanResult/RescheduleResult/View DTO
└── config/               # 配置、MySQL 连接/迁移、本地种子数据
```

## 核心业务规则

- **整批生成（PlanService.GeneratePlan）**：按任务类型偏移表生成计划时段，逐任务在对应资源类型内校验：资源离线、资源维护（状态或维护窗口重叠）、同资源时段与既有预约重叠、同批互撞、任务截止晚于离港、任务早于到站。任何一项不通过，事务回滚，任务和预约**都不保存**，409 响应的 `details[]` 逐项给出原因。
- **延误重排（DelayService.RegisterDelay）**：同一事务内登记延误事件；仅平移 `PLANNED/BLOCKED` 任务与其预约，`ACCEPTED/COMPLETED`（已签收/已完成）任务保持原计划时点；平移后撞上维护/离线/其它航班冻结预约的预约转为 `PENDING` 并写入机器可读 `conflict_reason`；航班到/离港时间与累计延误同步更新。
- **时段调整（BookingService.Adjust）**：再次执行维护/离线/重叠/离港截止校验，失败则保持 `PENDING` 并逐项返回原因；成功则恢复 `CONFIRMED`，并同步未签收关联任务的计划窗口。
- **签收冻结（TaskService.Accept）**：签收写入 `signed_at`，之后任何延误重排都不再移动该任务及其预约。

## 环境变量说明

- `COMPOSE_PROJECT_NAME`：Compose 项目名，默认 `ground-turn`
- `FRONTEND_PORT`：前端宿主机端口，默认 `20108`
- `BACKEND_PORT`：后端宿主机端口，默认 `21108`（容器内固定 3000）
- `DB_PORT`：MySQL 宿主机端口，默认 `33060`
- `DB_NAME / DB_USER / DB_PASSWORD`：数据库名与凭据
- `JWT_SECRET`：JWT 签名密钥
- `SEED_ON_START`：启动时写入本地演示种子（默认 true，已有用户时跳过）
- `RATE_LIMIT_QPS`：每 IP 令牌桶限流（默认 20）

新增配置需要同步：`.env.example`、`.env`、`docker-compose.yml`、`backend/src/config/config.go`。

## Docker 部署说明

- 根 `docker-compose.yml` 顶层 `name: ground-turn`，无 `version:` 字段。
- 容器名均为 `${COMPOSE_PROJECT_NAME:-ground-turn}-{db,backend,frontend}`。
- 数据库使用命名卷 `db_data`，不绑定挂载中文路径；`database/init.sql` 只读挂载到初始化目录。
- db 配置 healthcheck，backend `depends_on: condition: service_healthy`，frontend 再依赖 backend 健康。
- Nginx（`frontend/nginx.conf`）将 `/api/` 反代到 `http://backend:3000/api/`，SPA 使用 `try_files $uri $uri/ /index.html;`；前端代码只请求 `/api`，无硬编码 localhost。
- 常见问题：端口占用改 `.env` 后 `docker compose up -d`；重置数据执行 `docker compose down -v`。

## 枚举/常量出现位置清单

### GroundTaskType（CLEANING / CATERING / BAGGAGE / REFUEL / WATER_SERVICE / PUSHBACK）

- 后端：`constants/GroundTaskType.go`（枚举、中文文案、默认班组、偏移时长、资源类型映射）、`models/GroundTask.go`（字段）、`config/seed.go`、`services/PlanService.go`（构造与校验）、`services/GroundTask.go`、`constructors/GroundTask.go`（视图文案）、`constants/logTemplates.go`、`constants/errorMessages.go`、repositories、controllers、HTTP 测试
- 前端：`types/GroundTaskType.ts`、`constants/GroundTaskType.ts`（筛选项/颜色）、`constants/statusText.ts`、`components/common/StatusBadge.tsx`、`components/common/TeamTag.tsx`、`components/common/TurnaroundTimeline.tsx`、`pages/TasksPage.tsx`、`pages/TurnaroundsPage.tsx`、`constructors/GroundTaskConstructor.ts`、`stores/taskSlice.ts`

### TurnaroundStatus（ARRIVING / ON_STAND / IN_SERVICE / READY / DEPARTED / DELAYED）

- 后端：`constants/TurnaroundStatus.go`、`models/FlightTurnaround.go`、`config/seed.go`、`services/PlanService.go`、`services/DelayService.go`、`constructors/FlightTurnaround.go`、`constants/logTemplates.go`、controllers
- 前端：`types/TurnaroundStatus.ts`、`constants/TurnaroundStatus.ts`（颜色/筛选器）、`components/common/StatusBadge.tsx`、`components/common/DelayTag.tsx`、`pages/DashboardPage.tsx`、`pages/TurnaroundsPage.tsx`、`stores/turnaroundSlice.ts`

### ResourceStatus（AVAILABLE / BOOKED / MAINTENANCE / OFFLINE）

- 后端：`constants/ResourceStatus.go`（同文件含 BookingStatus CONFIRMED/PENDING/RELEASED 与任务状态）、`models/GroundResource.go`、`models/ResourceBooking.go`、`config/seed.go`、`services/PlanService.go`、`services/DelayService.go`、`services/BookingAdjustService.go`、`services/GroundResource.go`、`constructors/GroundResource.go`、`constructors/ResourceBooking.go`、`constants/conflictReasons.go`、`constants/logTemplates.go`
- 前端：`types/ResourceStatus.ts`、`constants/ResourceStatus.ts`（筛选器/颜色）、`constants/statusText.ts`（预约/任务状态颜色）、`components/common/StatusBadge.tsx`、`components/common/ConflictBadge.tsx`、`components/common/ResourceCalendar.tsx`、`components/common/PendingBookingList.tsx`、`pages/ResourcesPage.tsx`、`stores/resourceSlice.ts`

## 为什么会牵一发动全身

实体字段、枚举、日志模板、错误消息、DTO 构造器、冲突原因被刻意拆分到 `constants / types / constructors / services / controllers / stores / components` 多层并互相直接引用：

- 新增一个任务类型要同时改后端偏移表/班组/资源类型映射、种子、前端枚举、颜色、筛选项、状态徽章与展示页；
- 新增一个冲突原因要同时改后端 `conflictReasons`、检测逻辑、错误消息与前端 `CONFLICT_REASON_TEXT/COLOR`、ConflictBadge、ConflictPanel；
- 任意写操作都经 service 记录 `audit_log`，字段变更要同步 `logTemplates` 与调用处。

## License

MIT
