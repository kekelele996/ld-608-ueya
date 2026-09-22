import React from "react";
import { Layout, Menu, Button, Space, Tag } from "antd";
import { useNavigate, useLocation } from "react-router-dom";
import { useAppDispatch, useAppSelector } from "../stores/hooks";
import { logout } from "../stores/authSlice";
import { ROLE_TEXT } from "../types/auth";

const { Header, Sider, Content } = Layout;

const NAV = [
  { key: "/dashboard", label: "过站运行看板" },
  { key: "/turnarounds", label: "航班过站" },
  { key: "/tasks", label: "地勤任务" },
  { key: "/resources", label: "资源调度" },
  { key: "/delays", label: "延误归因" },
];

export const AppShell: React.FC<React.PropsWithChildren> = ({ children }) => {
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useAppDispatch();
  const { displayName, role, username } = useAppSelector((s) => s.auth);

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Header
        style={{
          display: "flex",
          alignItems: "center",
          gap: 16,
          background: "#001529",
          paddingInline: 24,
        }}
      >
        <div style={{ color: "#fff", fontSize: 18, fontWeight: 600, whiteSpace: "nowrap" }}>
          航空地勤周转保障平台
        </div>
        <Tag color="blue">ground-turn</Tag>
        <div style={{ flex: 1 }} />
        <Space>
          {username && (
            <span style={{ color: "rgba(255,255,255,0.85)" }}>
              {displayName} · {role ? ROLE_TEXT[role as keyof typeof ROLE_TEXT] : ""}
            </span>
          )}
          <Button size="small" onClick={() => dispatch(logout())}>
            切换角色
          </Button>
        </Space>
      </Header>
      <Layout>
        <Sider width={180} theme="light">
          <Menu
            mode="inline"
            selectedKeys={[NAV.find((n) => location.pathname.startsWith(n.key))?.key ?? "/dashboard"]}
            items={NAV.map((n) => ({ key: n.key, label: n.label }))}
            onSelect={({ key }) => navigate(key)}
            style={{ borderInlineEnd: "none", paddingTop: 8 }}
          />
        </Sider>
        <Content style={{ padding: 20, background: "#f0f2f5" }}>{children}</Content>
      </Layout>
    </Layout>
  );
};
