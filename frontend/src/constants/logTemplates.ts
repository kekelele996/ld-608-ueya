// Log action templates mirror backend constants/logTemplates.go. Every
// entity exposes >= 4 templates; audit pages render the formatted results.
export const LOG_TEMPLATES = {
  FlightTurnaround: ["航班过站 %s 已登记", "航班过站 %s 状态变更", "航班过站 %s 已放行离港", "航班过站导出"],
  GroundTask: ["航班 %s 整批生成保障计划", "任务 %d 已签收", "任务 %d 已完成", "任务 %d 标记阻塞", "航班 %s 保障计划校验失败整批未保存"],
  GroundResource: ["保障资源 %s 登记入库", "保障资源 %s 状态变更", "保障资源 %s 维护窗口登记", "预约 %d 时段调整"],
  ResourceBooking: ["资源预约 %d 生成", "预约 %d 进入待处理", "资源预约 %d 已释放", "资源预约导出"],
  DelayEvent: ["航班 %s 登记延误 +%d 分钟", "延误事件 %d 已关闭归因", "延误事件更新", "延误事件导出"],
} as const;
