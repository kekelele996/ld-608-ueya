import { configureStore } from "@reduxjs/toolkit";
import { TypedUseSelectorHook, useDispatch, useSelector } from "react-redux";
import turnaroundReducer from "./turnaroundSlice";
import taskReducer from "./taskSlice";
import resourceReducer from "./resourceSlice";
import bookingReducer from "./bookingSlice";
import delayReducer from "./delaySlice";
import authReducer from "./authSlice";

export const store = configureStore({
  reducer: {
    auth: authReducer,
    turnarounds: turnaroundReducer,
    tasks: taskReducer,
    resources: resourceReducer,
    bookings: bookingReducer,
    delays: delayReducer
  }
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export const useAppDispatch: () => AppDispatch = useDispatch;
export const useAppSelector: TypedUseSelectorHook<RootState> = useSelector;
