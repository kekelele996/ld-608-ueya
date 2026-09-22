// 日志模板前端镜像：与后端 constants/logTemplates.go 对应，
// 写操作触发的 message/notification 均引用此处，字段变更需同步后端模板。
export const LOG_TEMPLATES = {
  FlightTurnaround: {
    create: "航班过站登记：{flightNo} 机位 {standNo}",
    arrive: "航班 {flightNo} 到站确认，状态进入 ON_STAND",
    release: "航班 {flightNo} 放行",
    generate: "保障计划生成：航班 {flightNo} 共 {count} 项任务"
  },
  GroundTask: {
    dispatch: "任务批量派发：航班 {flightNo} 生成 {count} 项任务",
    sign: "任务签收：任务 #{taskId} 由班组签收",
    finish: "任务完成：任务 #{taskId}",
    block: "任务阻塞：任务 #{taskId}，原因 {note}"
  },
  ResourceBooking: {
    create: "资源预约生成：{resource} {start} ~ {end}",
    adjust: "预约时段调整：预约 #{bookingId} -> {status}",
    bump: "预约挤出：预约 #{bookingId} 进入待处理",
    release: "预约释放：预约 #{bookingId}",
    confirm: "预约确认：预约 #{bookingId}"
  },
  DelayEvent: {
    register: "延误登记：航班 {flightNo} 延误 {minutes} 分钟",
    replan: "延误重排：航班 {flightNo} 调整任务 {moved} 项，挤出 {bumped} 项",
    resolve: "延误关闭：事件 #{eventId}"
  },
  GroundResource: {
    maintenance: "资源维护登记：{resource} 状态 -> {status}"
  }
} as const;

export function renderLog(template: string, vars: Record<string, string | number>): string {
  return template.replace(/\{(\w+)\}/g, (_, key: string) => String(vars[key] ?? ""));
}
