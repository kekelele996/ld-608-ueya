export const DELAY_TYPES = [
  "WEATHER",
  "ATC",
  "CATERING",
  "BAGGAGE",
  "REFUEL",
  "MAINTENANCE",
  "OTHER"
] as const;
export type DelayTypeValue = (typeof DELAY_TYPES)[number];

export const DELAY_TYPE_TEXT: Record<DelayTypeValue, string> = {
  WEATHER: "天气",
  ATC: "空管流控",
  CATERING: "配餐延误",
  BAGGAGE: "行李延误",
  REFUEL: "加油延误",
  MAINTENANCE: "机务维护",
  OTHER: "其他"
};
