import React from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { Provider } from "react-redux";
import { App as AntApp, ConfigProvider } from "antd";
import zhCN from "antd/locale/zh_CN";
import { store } from "./stores/store";
import { router } from "./router/routes";
import "antd/dist/reset.css";
import "./styles.css";

createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ConfigProvider locale={zhCN}>
      <Provider store={store}>
        <AntApp>
          <RouterProvider router={router} />
        </AntApp>
      </Provider>
    </ConfigProvider>
  </React.StrictMode>,
);
