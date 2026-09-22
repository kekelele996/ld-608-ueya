import dayjs from "dayjs";

// Empty booking adjust window form; pages seed start/end from booking rows.
export interface BookingAdjustForm {
  start_time: string;
  end_time: string;
}

export const createBookingAdjustForm = (startTime?: string, durationMinutes = 25): BookingAdjustForm => {
  const base = startTime ? dayjs(startTime) : dayjs();
  return {
    start_time: base.toISOString(),
    end_time: base.add(durationMinutes, "minute").toISOString(),
  };
};
