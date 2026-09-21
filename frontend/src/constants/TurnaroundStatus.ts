export const TurnaroundStatus = ["ARRIVING","ON_STAND","IN_SERVICE","READY","DEPARTED","DELAYED"] as const;
export type TurnaroundStatus = (typeof TurnaroundStatus)[number];
export const TurnaroundStatusText: Record<TurnaroundStatus, string> = Object.fromEntries(TurnaroundStatus.map((value) => [value, value.replace(/_/g, " ")])) as Record<TurnaroundStatus, string>;
