import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import type { GroundResource, ResourceBooking } from "../types/entities";
import * as api from "../api/GroundResource";
import * as bookingApi from "../api/ResourceBooking";

interface ResourceState {
  resources: GroundResource[];
  bookings: ResourceBooking[];
  pending: ResourceBooking[];
  loading: boolean;
}

const initialState: ResourceState = {
  resources: [],
  bookings: [],
  pending: [],
  loading: false,
};

export const fetchResources = createAsyncThunk("resources/fetch", async () =>
  api.listResources({}),
);

export const fetchAllBookings = createAsyncThunk("resources/bookings", async () =>
  bookingApi.listBookings({}),
);

export const fetchPendingBookings = createAsyncThunk("resources/pending", async () => {
  const res = await bookingApi.listBookings({ status: "PENDING" });
  return res.items;
});

export const fetchCalendar = createAsyncThunk(
  "resources/calendar",
  async (range: { from?: string; to?: string }) => api.getCalendar(range.from, range.to),
);

const slice = createSlice({
  name: "resources",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchResources.fulfilled, (state, action) => {
        state.resources = action.payload.items;
      })
      .addCase(fetchAllBookings.fulfilled, (state, action) => {
        state.bookings = action.payload.items;
      })
      .addCase(fetchPendingBookings.fulfilled, (state, action) => {
        state.pending = action.payload;
      })
      .addCase(fetchCalendar.fulfilled, (state, action) => {
        state.resources = action.payload.resources;
        state.bookings = action.payload.items;
      });
  },
});

export default slice.reducer;
