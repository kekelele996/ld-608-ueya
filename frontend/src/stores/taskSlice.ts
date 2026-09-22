import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import type { GroundTask } from "../types/entities";
import * as api from "../api/GroundTask";

interface TaskState {
  rows: GroundTask[];
  total: number;
  loading: boolean;
}

const initialState: TaskState = { rows: [], total: 0, loading: false };

export const fetchTasks = createAsyncThunk(
  "tasks/fetch",
  async (params?: { status?: string }) => api.listTasks({ ...params, page: 1, page_size: 200 }),
);

export const acceptTask = createAsyncThunk("tasks/accept", async (id: number) => api.acceptTask(id));
export const completeTask = createAsyncThunk("tasks/complete", async (id: number) => api.completeTask(id));
export const blockTask = createAsyncThunk(
  "tasks/block",
  async (payload: { id: number; note: string }) => api.blockTask(payload.id, payload.note),
);

const slice = createSlice({
  name: "tasks",
  initialState,
  reducers: {},
  extraReducers: (builder) => {
    builder
      .addCase(fetchTasks.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchTasks.fulfilled, (state, action) => {
        state.rows = action.payload.items;
        state.total = action.payload.total;
        state.loading = false;
      })
      .addCase(fetchTasks.rejected, (state) => {
        state.loading = false;
      });
  },
});

export default slice.reducer;
