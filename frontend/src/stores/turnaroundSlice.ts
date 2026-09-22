import { createAsyncThunk, createSlice } from "@reduxjs/toolkit";
import { turnaroundApi } from "../api/FlightTurnaround";
import type { FlightTurnaround } from "../types/FlightTurnaround";
import type { TurnaroundDetail } from "../types/api";

export const fetchTurnarounds = createAsyncThunk("turnarounds/fetchAll", async () => {
  return turnaroundApi.list();
});

export const fetchTurnaroundDetail = createAsyncThunk(
  "turnarounds/fetchDetail",
  async (id: number) => turnaroundApi.detail(id)
);

interface State {
  rows: FlightTurnaround[];
  detail: TurnaroundDetail | null;
  loading: boolean;
}

const initialState: State = { rows: [], detail: null, loading: false };

const slice = createSlice({
  name: "turnarounds",
  initialState,
  reducers: {
    clearDetail(state) {
      state.detail = null;
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchTurnarounds.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchTurnarounds.fulfilled, (state, action) => {
        state.rows = action.payload;
        state.loading = false;
      })
      .addCase(fetchTurnarounds.rejected, (state) => {
        state.loading = false;
      })
      .addCase(fetchTurnaroundDetail.fulfilled, (state, action) => {
        state.detail = action.payload;
      });
  }
});

export const { clearDetail } = slice.actions;
export default slice.reducer;
