import { Layout, Menu, Button, Tag, Space } from "antd";
import {
  DashboardOutlined,
  PlaneOutlineParams,
  UnorderedListOutlined,
  AppstoreOutlined,
  FieldTimeOutlined,
  LogoutOutlined
} from "./menuIcons";
import { useLocation, useNavigate, Outlet } from "react-router-dom";
import { ROUTES } from "../router/routes";
import { useAppDispatch, useAppSelector } from "../stores/store";
import { logout } from "../stores/authSlice";
import { USER_ROLE_TEXT } from "../constants/UserRole";

const ICONS: Record<string, React.ReactNode> = {
  dashboard: <DashboardOutlined />,
  plane: <PlaneOutlineParams />,
  tasks: <UnorderedListOutlined />,
  resource: <AppstoreOutlined />,
  delay: <FieldTimeOutlined />
};

const { Header, Sider, Content } = Layout;

export function AppLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const dispatch = useAppDispatch();
  const user = useAppSelector((s) => s.auth.user);

  const selectedKey =
    ROUTES.find((r) => location.pathname.startsWith(r.path) && r.path !== "/dashboard")?.path ??
    (location.pathname.startsWith("/turnarounds/") ? "/turnarounds" : location.pathname);

  return (
    <Layout className="app-layout">
      <Sider breakpoint="lg" collapsedWidth="0" width={220} theme="dark">
        <div className="brand">
          <span className="brand-mark">GT</span>
          <span>航空地勤周转保障</span>
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[selectedKey]}
          items={ROUTES.map((route) => ({
            key: route.path,
            icon: ICONS[route.icon],
            label: route.name,
            onClick: () => navigate(route.path)
          }))}
        />
      </Sider>
      <Layout>
        <Header className="app-header">
          <div className="header-title">
            {ROUTES.find((r) => selectedKey === r.path)?.name ?? "过站运行看板"}
          </div>
          <Space>
            {user && (
              <>
                <Tag color="blue">{USER_ROLE_TEXT[user.role]}</Tag>
                <span className="header-user">{user.name}</span>
              </>
            )}
            <Button
              size="small"
              icon={<LogoutOutlined />}
              onClick={() => {
                dispatch(logout());
                navigate("/login");
              }}
            >
              退出
            </Button>
          </Space>
        </Header>
        <Content className="app-content">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
}
