import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { groundTaskApi } from "../api/GroundTask";
import type { GroundTask } from "../types/GroundTask";

export const fetchTasks = createAsyncThunk("tasks/fetchAll", async () => groundTaskApi.list());

interface State {
  rows: GroundTask[];
  loading: boolean;
}

const slice = createSlice({
  name: "tasks",
  initialState: { rows: [], loading: false } as State,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchTasks.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchTasks.fulfilled, (state, action) => {
        state.rows = action.payload;
        state.loading = false;
      })
      .addCase(fetchTasks.rejected, (state) => {
        state.loading = false;
      });
  }
});

export default slice.reducer;
