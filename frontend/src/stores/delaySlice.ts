import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import type { DelayEvent } from "../types/entities";
import * as api from "../api/DelayEvent";

interface DelayState {
  rows: DelayEvent[];
  total: number;
  loading: boolean;
}

const initialState: DelayState = { rows: [], total: 0, loading: false };

export const fetchDelays = createAsyncThunk("delays/fetch", async () =>
  api.listDelays({ page: 1, page_size: 100 }),
);

const slice = createSlice({
  name: "delays",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchDelays.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchDelays.fulfilled, (state, action) => {
        state.rows = action.payload.items;
        state.total = action.payload.total;
        state.loading = false;
      })
      .addCase(fetchDelays.rejected, (state) => {
        state.loading = false;
      });
  },
});

export default slice.reducer;
