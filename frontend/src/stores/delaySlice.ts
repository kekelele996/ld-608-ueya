import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { delayEventApi, type DelayImpact } from "../api/DelayEvent";
import type { DelayEvent } from "../types/DelayEvent";

export const fetchDelays = createAsyncThunk("delays/fetchAll", async () => delayEventApi.list());
export const fetchDelayImpacts = createAsyncThunk("delays/fetchImpacts", async () =>
  delayEventApi.impacts()
);

const slice = createSlice({
  name: "delays",
  initialState: {
    rows: [] as DelayEvent[],
    impacts: [] as DelayImpact[],
    loading: false
  },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchDelays.fulfilled, (state, action) => {
        state.rows = action.payload;
      })
      .addCase(fetchDelayImpacts.fulfilled, (state, action) => {
        state.impacts = action.payload;
      });
  }
});

export default slice.reducer;
