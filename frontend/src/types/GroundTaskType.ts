export type GroundTaskType =
  | "CLEANING"
  | "CATERING"
  | "BAGGAGE"
  | "REFUEL"
  | "WATER_SERVICE"
  | "PUSHBACK";

export const GROUND_TASK_TYPES: GroundTaskType[] = [
  "CLEANING",
  "CATERING",
  "BAGGAGE",
  "REFUEL",
  "WATER_SERVICE",
  "PUSHBACK",
];

export const GROUND_TASK_TYPE_TEXT: Record<GroundTaskType, string> = {
  CLEANING: "客舱清洁",
  CATERING: "航食配送",
  BAGGAGE: "行李装卸",
  REFUEL: "航油加注",
  WATER_SERVICE: "清水污水",
  PUSHBACK: "飞机推出",
};

export const TASK_RESOURCE_TYPE: Record<GroundTaskType, string> = {
  CLEANING: "CLEANING_KIT",
  CATERING: "CATERING_TRUCK",
  BAGGAGE: "BELT_LOADER",
  REFUEL: "FUEL_TRUCK",
  WATER_SERVICE: "WATER_TRUCK",
  PUSHBACK: "PUSHBACK_TRACTOR",
};
