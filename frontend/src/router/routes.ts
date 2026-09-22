export interface AppRoute {
  path: string;
  name: string;
  icon: string;
}

export const ROUTES: AppRoute[] = [
  { path: "/dashboard", name: "过站运行看板", icon: "dashboard" },
  { path: "/turnarounds", name: "航班过站", icon: "plane" },
  { path: "/tasks", name: "地勤任务", icon: "tasks" },
  { path: "/resources", name: "资源调度", icon: "resource" },
  { path: "/delays", name: "延误归因", icon: "delay" }
];
