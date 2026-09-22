# 航空地勤周转保障平台 · ground-turn

面向机场地勤团队的**航班过站保障、资源调度、异常延误重排与任务签收**全栈系统。
航班到站后按任务类型/班组/资源可用时段**整批**生成地勤任务与资源预约；任何冲突（资源时段重叠、资源维护/离线、任务截止冲突）都会导致**整批不保存并逐项说明原因**；登记延误分钟数后，未签收任务与预约整体重排，**已签收任务保持原时点**，被挤出的预约进入**待处理**状态，页面可查看冲突、调整时段、签收并刷新看到同一份结果。

---

## 快速启动（推荐，Docker Compose 一键）

```bash
cp .env.example .env && docker compose up -d
```

启动后：

- 前端：<http://localhost:20108>
- 后端健康检查：<http://localhost:21108/health>

首次启动会自动建库建表（`database/init.sql` + GORM AutoMigrate 幂等对账），并写入本地演示种子数据（航班、资源、班组、4 个演示账号）。

### 演示账号（RBAC）

| 角色 | 账号 | 密码 | 主要权限 |
|---|---|---|---|
| 地勤调度 DISPATCHER | `dispatcher` | `dispatch123` | 过站登记、生成计划、延误登记、任务签收、放行 |
| 班组 CREW | `crew` | `crew123` | 任务签收/完成/阻塞、延误登记 |
| 资源管理员 RESOURCE | `resource` | `resource123` | 资源维护/离线、预约调整/释放 |
| 运行督导 SUPERVISOR | `supervisor` | `super123` | 全部只读 + 关闭延误归因 |

---

## 一分钟走通业务闭环

1. 用 `dispatcher` 登录 → 「航班过站」→ 打开航班 **CA1831（203 机位）**。
2. 点击 **生成保障任务与资源预约**：按 6 类任务（行李/清洁/配餐/清水污水/加油/推出）、对应班组与资源时段整批生成 6 项任务 + 6 条预约。
3. 打开 **CZ3108**（过站时间仅 40 分钟）再点生成：弹出 409 冲突面板，**整批未保存**，逐项列出「任务截止冲突」的任务、计划结束与截止时间。
4. 回到 CA1831，把「行李装卸」任务**签收**（签收后计划时点被冻结）。
5. 点击 **延误登记并重排**，输入 40 分钟：
   - 已签收的行李任务**保持 08:00 原时点**（重排结果里单列说明）；
   - 其余 5 项未完成任务整体顺延 40 分钟；
   - 配餐预约平移后撞上 CATER-01 既有占用 → 预约被**挤出为待处理**并带冲突原因。
6. 在详情页或「资源调度 → 待处理预约」点击 **调整时段**：
   - 改到空闲窗口 → 状态回到「已确认」；
   - 改到占用窗口/维护资源 → 仍为「待处理」并刷新出新的冲突原因。
7. 所有任务完成后 **放行航班**。刷新页面/重新登录，看到的始终是同一份持久化结果。

---

## 访问地址与 CLI 示例

```bash
# 登录获取 JWT
curl -s -X POST http://localhost:21108/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"dispatcher","password":"dispatch123"}'

# 整批生成（成功）
curl -s -X POST http://localhost:21108/api/turnarounds/1/generate-plan \
  -H "Authorization: Bearer <token>"

# 冲突时返回 HTTP 409，信封 data.items 逐项给出冲突原因
# 延误登记（触发重排/挤出）
curl -s -X POST http://localhost:21108/api/turnarounds/1/delays \
  -H "Authorization: Bearer <token>" -H 'Content-Type: application/json' \
  -d '{"delay_type":"ATC","minutes":40,"root_cause":"出港流控"}'
```

统一响应信封：成功 `{ "ok": true, "data": ... }`；失败 `{ "ok": false, "error": { "code", "message" }, "data"?: ... }`
（整批冲突时 `data` 携带逐项原因，HTTP 409）。

---

## 本地开发方式

- 后端（Go 1.22，可用 SQLite 免装 MySQL）：

  ```bash
  cd backend
  go run .                      # 默认连 MySQL（127.0.0.1:3306）
  # 或用 SQLite 本地调试：
  SQLITE_PATH=/tmp/ground-turn.db PORT=3000 go run .
  go test ./...                 # 闭环集成测试（整批拒绝/签收冻结/延误挤出/改期）
  ```

- 前端（Vite 已把 `/api` 代理到 `http://127.0.0.1:21108`）：

  ```bash
  cd frontend
  npm install
  npm run dev                   # http://localhost:20108
  ```

  前端代码统一请求 `/api`，未硬编码任何 `localhost` 主机名；生产由 Nginx 反代到 `http://backend:3000/`。

---

## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | React 18 + TypeScript + Vite + Ant Design 5 + Redux Toolkit + react-router 6 + dayjs |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0（本地/测试可用 SQLite，纯 Go 驱动 glebarez/sqlite） |
| 认证鉴权 | JWT（golang-jwt/v5，bcrypt 口令）+ RBAC 四角色矩阵 |
| 部署 | Docker Compose（db / backend / frontend，命名卷 + healthcheck） |

---

## 项目目录结构

```text
.
├── docker-compose.yml          # name: ground-turn，三服务编排
├── .env / .env.example         # 端口/账号/JWT/种子开关
├── database/init.sql           # MySQL 8 完整 DDL（7 张表）
├── backend/
│   ├── main.go
│   └── src/
│       ├── config/             # 环境配置
│       ├── database/           # 连接/迁移/种子
│       ├── constants/          # 枚举、错误码、错误消息、日志模板、RBAC、任务蓝图、班组
│       ├── models/             # GORM 模型（5 实体 + AuditLog/User）
│       ├── repositories/       # 数据访问层
│       ├── services/           # 事务/冲突检测/整批生成/延误重排（含闭环测试）
│       ├── controllers/        # 控制器（各自包装 AppError）
│       ├── routes/             # 路由装配 + 中间件挂载
│       ├── middlewares/        # auth/rbac/auditLog/errorHandler/rateLimit
│       ├── constructors/       # 请求/模型工厂
│       ├── types/              # DTO、错误信封、聚合视图
│       └── utils/              # JWT、时间格式化、日志渲染
└── frontend/
    ├── nginx.conf              # /api/ 反代 + SPA try_files
    └── src/
        ├── api/                # 按实体分文件 + 统一请求封装
        ├── stores/             # Redux Toolkit 按实体分片 + auth
        ├── types/              # 共享类型
        ├── constants/          # 枚举/错误码/日志模板/状态文案/RBAC 镜像
        ├── constructors/       # 默认对象与表单工厂
        ├── components/common/  # StatusBadge/TurnaroundTimeline/ResourceCalendar/
        │                       #   DelayTag/ConflictBadge/StatCard/TeamTag/TimelineList/AdjustBookingModal
        ├── hooks/              # useTurnaroundProgress/useResourceConflict/usePagination
        ├── pages/              # Dashboard/Turnarounds/TurnaroundDetail/Tasks/Resources/Delays/Login
        ├── router/             # 路由表 + RequireAuth 守卫
        ├── utils/formatters.ts # 日期/数字/风险混合格式化
        └── mocks/              # 本地种子快照（仅只读兜底）
```

---

## 环境变量说明

| 变量 | 默认值 | 说明 |
|---|---|---|
| `COMPOSE_PROJECT_NAME` | `ground-turn` | Compose 项目名与容器名前缀 |
| `FRONTEND_PORT` | `20108` | 前端宿主机端口（容器内 80） |
| `BACKEND_PORT` | `21108` | 后端宿主机端口（容器内 3000） |
| `DB_PORT` | `33060` | MySQL 宿主机端口（容器内 3306） |
| `DB_NAME / DB_USER / DB_PASSWORD` | `app_db / app_user / app_password` | 数据库凭据 |
| `JWT_SECRET` | `local-dev-secret` | JWT 签名密钥（生产请修改） |
| `SEED_ON_START` | `true` | 启动空库时写入演示种子 |
| `SQLITE_PATH` | 空 | 仅本地开发：设置后使用 SQLite 文件而非 MySQL |

---

## Docker 部署说明

- 根 Compose 无 `version:` 字段，顶层 `name: ground-turn`；容器名均为 `${COMPOSE_PROJECT_NAME:-ground-turn}-{db,backend,frontend}`，在任意目录名（含中文）下均可启动。
- 端口映射：`${FRONTEND_PORT:-20108}:80`、`${BACKEND_PORT:-21108}:3000`。
- 数据库使用**命名卷** `db_data`，不绑定挂载到宿主机路径，规避中文路径权限问题。
- `db` 配置 `mysqladmin ping` healthcheck；`backend` 通过 `depends_on: condition: service_healthy` 等待数据库，并在应用层再做约 60 秒连接重试；`frontend` 等待后端健康。
- 后端镜像为 CGO_ENABLED=0 的静态二进制，最终镜像基于 alpine，自带 `/health` 与容器 healthcheck。
- 常见问题：
  - 端口占用 → 修改 `.env` 中 `FRONTEND_PORT/BACKEND_PORT/DB_PORT` 后 `docker compose up -d`。
  - 重置演示数据 → `docker compose down -v && docker compose up -d`（删除命名卷）。
  - 验证编排 → `docker compose config --quiet`。

---

## 枚举/常量出现位置清单

> 三组枚举均在**前后端重复定义**，新增枚举值会强制触达：常量、类型、构造器、日志模板、错误消息、筛选器、列表/详情展示组件与种子数据。

### GroundTaskType（CLEANING / CATERING / BAGGAGE / REFUEL / WATER_SERVICE / PUSHBACK）

- 后端：`backend/src/constants/GroundTaskType.go`（枚举+中文文案）、`constants/taskBlueprints.go`（时间偏移/时长/截止）、`constants/teams.go`（类型→班组/资源池）、`constructors/GroundTask.go`、`services/TurnaroundService.go`（生成/截止校验）、`database/seed.go`、`utils/logger.go`（日志模板入参）。
- 前端：`frontend/src/constants/GroundTaskType.ts`、`types/GroundTask.ts`、`constructors/GroundTaskConstructor.ts`、`constants/logTemplates.ts`、`constants/errorMessages.ts`、列表筛选器（`pages/TasksPage.tsx`、`pages/TurnaroundsPage.tsx`）、展示组件（`StatusBadge`/`TurnaroundTimeline`/`ResourceCalendar`/`TeamTag`）、`mocks/seedData.ts`。

### TurnaroundStatus（ARRIVING / ON_STAND / IN_SERVICE / READY / DEPARTED / DELAYED）

- 后端：`constants/TurnaroundStatus.go`、`models/FlightTurnaround.go`、`services/TurnaroundService.go`（到站/生成/放行/延误置 DELAYED）、`services/ReportService.go`（看板聚合）、`database/seed.go`。
- 前端：`constants/TurnaroundStatus.ts`、`types/FlightTurnaround.ts`、`constructors/FlightTurnaroundConstructor.ts`、筛选器（`TurnaroundsPage`）、展示（`StatusBadge`、`DashboardPage`、详情页 Descriptions）。

### ResourceStatus（AVAILABLE / BOOKED / MAINTENANCE / OFFLINE）

- 后端：`constants/ResourceStatus.go`、`models/GroundResource.go`、`services/conflict.go`（维护/离线/可用时段判定）、`services/ResourceService.go`（状态维护）、`database/seed.go`。
- 前端：`constants/ResourceStatus.ts`、`types/GroundResource.ts`、`constructors/GroundResourceConstructor.ts`、资源台账筛选与状态列（`ResourcesPage`）、`ResourceCalendar`、`AdjustBookingModal`（资源下拉文案）、`StatusBadge`。

### 其他贯穿全栈的枚举

- `GroundTaskStatus`（PLANNED/SIGNED/FINISHED/BLOCKED）、`BookingStatus`（CONFIRMED/PENDING/RELEASED）、`DelayType`（WEATHER/ATC/...）、`UserRole`（DISPATCHER/CREW/RESOURCE/SUPERVISOR）：同样各自在 `backend/src/constants/*.go` 与 `frontend/src/constants/*.ts` 成对出现，并被路由守卫、RBAC 矩阵、store 选择器、按钮显隐、日志模板和错误码引用。

---

## 为什么该项目会「牵一发动全身」

- 枚举前后端重复定义，且被常量、类型、构造器、日志模板、错误消息、筛选器、多个共享组件和种子数据共同引用；新增一个状态值必须同步十余个文件。
- 日志模板集中在 `constants/logTemplates.*`（后端 Go map、前端 `LOG_TEMPLATES`），所有写操作都渲染模板；字段改名要同步模板与所有调用处。
- 错误码与错误消息分离（`errorCodes.*` / `errorMessages.*`），service 与 controller **分别包装** `AppError`，全局中间件只兜底未捕获异常，不在单点吞错。
- 每个实体都有独立 constructor/factory，页面、store、service 不允许散写默认结构。
- `utils/formatters.ts` 故意混合日期、分钟、风险等级格式化，被多个页面与组件共同依赖。
- 配置分散经过 `.env.example`、`.env`、`docker-compose.yml`、`backend/src/config/config.go`、前端请求封装（`/api`）与 Nginx 反代，新增配置需多处同步。
- 核心事务逻辑（整批生成、延误重排）横跨 model → repository → service → controller → 五个页面与共享组件，一个规则改动会强制触及 3–5 个以上文件。

---

## 核心事务与一致性规则

- **整批生成**：在单个数据库事务内投影全部任务与预约；只要存在任一资源时段重叠、资源 MAINTENANCE/OFFLINE、可用时段不覆盖或任务晚于离港截止，事务回滚 → 任务与预约**一条都不落库**，响应逐项给出 `RESOURCE_TIME_OVERLAP / RESOURCE_MAINTENANCE / RESOURCE_OFFLINE / RESOURCE_WINDOW_UNAVAILABLE / TASK_DEADLINE_CONFLICT / NO_MATCHING_RESOURCE`。
- **延误重排**：单个事务内写入延误事件、平移航班时刻；`SIGNED/FINISHED` 任务及其预约冻结，`PLANNED/BLOCKED` 任务顺延并重算截止；平移后冲突的预约置 `PENDING` 并写 `conflict_reason`。
- **调整时段**：同样的冲突规则即时校验；改到合法窗口回到 `CONFIRMED`，否则保持 `PENDING` 并回写原因——页面刷新后结果一致。

## License

MIT
