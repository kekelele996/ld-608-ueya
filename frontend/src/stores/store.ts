import { configureStore } from "@reduxjs/toolkit";
import authReducer from "./authSlice";
import turnaroundReducer from "./turnaroundSlice";
import taskReducer from "./taskSlice";
import resourceReducer from "./resourceSlice";
import delayReducer from "./delaySlice";

export const store = configureStore({
  reducer: {
    auth: authReducer,
    turnarounds: turnaroundReducer,
    tasks: taskReducer,
    resources: resourceReducer,
    delays: delayReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
