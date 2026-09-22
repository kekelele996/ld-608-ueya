import { useEffect } from "react";
import { BrowserRouter, Navigate, Route, Routes, useNavigate } from "react-router-dom";
import { Provider } from "react-redux";
import { ConfigProvider, Spin } from "antd";
import zhCN from "antd/locale/zh_CN";
import { store, useAppDispatch, useAppSelector } from "./stores/store";
import { hydrate } from "./stores/authSlice";
import { authApi } from "./api/auth";
import { tokenStore } from "./api/client";
import { AppLayout } from "./components/AppLayout";
import { RequireAuth } from "./router/RequireAuth";
import { LoginPage } from "./pages/LoginPage";
import { DashboardPage } from "./pages/DashboardPage";
import { TurnaroundsPage } from "./pages/TurnaroundsPage";
import { TurnaroundDetailPage } from "./pages/TurnaroundDetailPage";
import { TasksPage } from "./pages/TasksPage";
import { ResourcesPage } from "./pages/ResourcesPage";
import { DelaysPage } from "./pages/DelaysPage";

function Bootstrap({ children }: { children: React.ReactNode }) {
  const dispatch = useAppDispatch();
  const user = useAppSelector((s) => s.auth.user);
  const navigate = useNavigate();

  useEffect(() => {
    // 刷新页面后用 token 恢复用户信息；失效则清理并回登录页。
    if (tokenStore.get() && !user) {
      authApi
        .me()
        .then((info) => dispatch(hydrate(info)))
        .catch(() => {
          tokenStore.clear();
          navigate("/login");
        });
    }
  }, [dispatch, user, navigate]);

  if (tokenStore.get() && !user) {
    return (
      <div className="bootstrap-loading">
        <Spin tip="加载保障数据..." />
      </div>
    );
  }
  return <>{children}</>;
}

export default function App() {
  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        token: {
          colorPrimary: "#155eef",
          borderRadius: 8,
          fontSize: 14
        }
      }}
    >
      <Provider store={store}>
        <BrowserRouter>
          <Bootstrap>
            <Routes>
              <Route path="/login" element={<LoginPage />} />
              <Route
                element={
                  <RequireAuth>
                    <AppLayout />
                  </RequireAuth>
                }
              >
                <Route path="/dashboard" element={<DashboardPage />} />
                <Route path="/turnarounds" element={<TurnaroundsPage />} />
                <Route path="/turnarounds/:id" element={<TurnaroundDetailPage />} />
                <Route path="/tasks" element={<TasksPage />} />
                <Route path="/resources" element={<ResourcesPage />} />
                <Route path="/delays" element={<DelaysPage />} />
              </Route>
              <Route path="*" element={<Navigate to="/dashboard" replace />} />
            </Routes>
          </Bootstrap>
        </BrowserRouter>
      </Provider>
    </ConfigProvider>
  );
}
