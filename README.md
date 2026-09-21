# 航空地勤周转保障平台

面向机场地勤团队的航班过站保障、资源调度、异常延误和任务签收系统。

## 快速启动

```bash
cp .env.example .env && docker compose up -d
```

## 访问地址或 CLI 示例

前端：<http://localhost:20108>

后端健康检查：<http://localhost:21108/health>


## 本地开发方式

- 前端：`cd frontend && npm install && npm run dev`
- 后端：进入 `backend` 后按技术栈运行开发命令，接口统一挂在 `/api`。


## 技术栈

| 层 | 技术 |
|---|---|
| 前端 | React 18 + TypeScript + Vite + Ant Design + Redux Toolkit |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | MySQL 8.0 |
| 部署 | Docker Compose |

## 项目目录结构

```text
frontend/src/api, stores, types, constants, constructors, components/common, hooks, pages, router, utils, mocks
backend/src/routes, controllers, services, models, repositories, middlewares, constants, constructors, utils, types, config
```

## 环境变量说明

- `COMPOSE_PROJECT_NAME`: Compose 项目名，默认 `ground-turn`
- `FRONTEND_PORT`: 前端端口，默认 `20108`
- `BACKEND_PORT`: 后端端口，默认 `21108`
- `DB_PORT`: 数据库宿主机端口
- `DB_USER/DB_PASSWORD/DB_NAME`: 本地数据库凭据

## Docker 部署说明

- 根 Compose 文件不写 `version`，顶层 `name: ground-turn`。
- 容器名均使用 `${COMPOSE_PROJECT_NAME:-ground-turn}` 前缀。
- 数据库使用命名卷，避免绑定中文路径。
- 常见问题：端口占用时修改 `.env` 中端口后重启；需要重置数据时执行 `docker compose down -v`。

## 枚举/常量出现位置清单

- GroundTaskType: constants/GroundTaskType、types/GroundTaskType、constructors、logTemplates、errorMessages、筛选器、展示组件/控制器均有引用。
- TurnaroundStatus: constants/TurnaroundStatus、types/TurnaroundStatus、constructors、logTemplates、errorMessages、筛选器、展示组件/控制器均有引用。
- ResourceStatus: constants/ResourceStatus、types/ResourceStatus、constructors、logTemplates、errorMessages、筛选器、展示组件/控制器均有引用。

## 为什么会牵一发动全身

实体字段、枚举、日志模板、错误消息、构造器、筛选器和展示组件被刻意拆散到多个目录；修改一个状态值通常需要同步类型、构造器、服务、控制器、store、页面、README 与数据库种子。

## License

MIT
