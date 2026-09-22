import dayjs from "dayjs";
import type { ResourceBooking } from "../types/ResourceBooking";

export function createDefaultResourceBooking(
  overrides: Partial<ResourceBooking> = {}
): ResourceBooking {
  return {
    id: 0,
    resource_id: 0,
    turnaround_id: 0,
    task_id: null,
    start_time: dayjs().format(),
    end_time: dayjs().add(30, "minute").format(),
    booking_status: "PENDING",
    conflict_reason: "",
    ...overrides
  };
}

// 调整时段弹窗的表单默认值：以被挤出预约当前窗口为基准。
export function createAdjustBookingForm(booking: ResourceBooking) {
  return {
    resource_id: booking.resource_id,
    range: [dayjs(booking.start_time), dayjs(booking.end_time)] as [dayjs.Dayjs, dayjs.Dayjs]
  };
}
