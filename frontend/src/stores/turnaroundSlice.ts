import { createSlice, createAsyncThunk } from "@reduxjs/toolkit";
import type { FlightTurnaround, TurnaroundDetail } from "../types/entities";
import * as api from "../api/FlightTurnaround";

interface TurnaroundState {
  rows: FlightTurnaround[];
  total: number;
  loading: boolean;
  detail: TurnaroundDetail | null;
  detailLoading: boolean;
}

const initialState: TurnaroundState = {
  rows: [],
  total: 0,
  loading: false,
  detail: null,
  detailLoading: false,
};

export const fetchTurnarounds = createAsyncThunk("turnarounds/fetch", async () => {
  return api.listTurnarounds({ page: 1, page_size: 100 });
});

export const fetchTurnaroundDetail = createAsyncThunk(
  "turnarounds/fetchDetail",
  async (id: number) => api.getTurnaround(id),
);

const slice = createSlice({
  name: "turnarounds",
  initialState,
  reducers: {
    clearDetail(state) {
      state.detail = null;
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchTurnarounds.pending, (state) => {
        state.loading = true;
      })
      .addCase(fetchTurnarounds.fulfilled, (state, action) => {
        state.rows = action.payload.items;
        state.total = action.payload.total;
        state.loading = false;
      })
      .addCase(fetchTurnarounds.rejected, (state) => {
        state.loading = false;
      })
      .addCase(fetchTurnaroundDetail.pending, (state) => {
        state.detailLoading = true;
      })
      .addCase(fetchTurnaroundDetail.fulfilled, (state, action) => {
        state.detail = action.payload;
        state.detailLoading = false;
      })
      .addCase(fetchTurnaroundDetail.rejected, (state) => {
        state.detailLoading = false;
      });
  },
});

export const { clearDetail } = slice.actions;
export default slice.reducer;
