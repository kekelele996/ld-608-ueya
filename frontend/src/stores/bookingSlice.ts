import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { resourceBookingApi } from "../api/ResourceBooking";
import type { ResourceBooking } from "../types/ResourceBooking";

export const fetchBookings = createAsyncThunk("bookings/fetchAll", async () =>
  resourceBookingApi.list()
);
export const fetchPendingBookings = createAsyncThunk("bookings/fetchPending", async () =>
  resourceBookingApi.listPending()
);

const slice = createSlice({
  name: "bookings",
  initialState: {
    rows: [] as ResourceBooking[],
    pending: [] as ResourceBooking[],
    loading: false
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchBookings.fulfilled, (state, action) => {
        state.rows = action.payload;
      })
      .addCase(fetchPendingBookings.fulfilled, (state, action) => {
        state.pending = action.payload;
      });
  }
});

export default slice.reducer;
