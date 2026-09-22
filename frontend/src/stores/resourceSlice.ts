import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { groundResourceApi } from "../api/GroundResource";
import type { GroundResource } from "../types/GroundResource";

export const fetchResources = createAsyncThunk("resources/fetchAll", async () =>
  groundResourceApi.list()
);

const slice = createSlice({
  name: "resources",
  initialState: { rows: [] as GroundResource[], loading: false },
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchResources.fulfilled, (state, action) => {
        state.rows = action.payload;
      })
      .addCase(fetchResources.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchResources.rejected, (state) => {
        state.loading = false;
      });
  }
});

export default slice.reducer;
