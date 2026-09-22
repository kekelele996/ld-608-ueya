export const BOOKING_STATUSES = ["CONFIRMED", "PENDING", "RELEASED"] as const;
export type BookingStatusValue = (typeof BOOKING_STATUSES)[number];

export const BOOKING_STATUS_TEXT: Record<BookingStatusValue, string> = {
  CONFIRMED: "已确认",
  PENDING: "待处理",
  RELEASED: "已释放"
};
