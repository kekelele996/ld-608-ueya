export const GroundTaskType = ["CLEANING","CATERING","BAGGAGE","REFUEL","WATER_SERVICE","PUSHBACK"] as const;
export type GroundTaskType = (typeof GroundTaskType)[number];
export const GroundTaskTypeText: Record<GroundTaskType, string> = Object.fromEntries(GroundTaskType.map((value) => [value, value.replace(/_/g, " ")])) as Record<GroundTaskType, string>;
