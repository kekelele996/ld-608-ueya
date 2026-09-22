import { Button, Card, Form, Input, Tag } from "antd";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { authApi } from "../api/auth";
import { useAppDispatch } from "../stores/store";
import { loginSucceeded } from "../stores/authSlice";
import { SEED_ACCOUNTS } from "../constants/UserRole";
import { ApiError } from "../api/client";

export function LoginPage() {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);

  const onFinish = async (values: { username: string; password: string }) => {
    setLoading(true);
    try {
      const result = await authApi.login(values.username, values.password);
      dispatch(loginSucceeded(result));
      navigate("/dashboard");
    } catch (err) {
      // ApiError 已在 request 中弹出统一 message；此处只需停止 loading。
      if (err instanceof ApiError) return;
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login-page">
      <Card className="login-card" title={<div className="login-title">航空地勤周转保障平台</div>}>
        <Form layout="vertical" onFinish={onFinish} initialValues={{ username: "dispatcher", password: "dispatch123" }}>
          <Form.Item name="username" label="账号" rules={[{ required: true, message: "请输入账号" }]}>
            <Input placeholder="dispatcher / crew / resource / supervisor" autoComplete="username" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, message: "请输入密码" }]}>
            <Input.Password autoComplete="current-password" />
          </Form.Item>
          <Button type="primary" htmlType="submit" block loading={loading}>
            登录
          </Button>
        </Form>
        <div className="login-accounts">
          <div className="muted">演示账号（点击自动填充）：</div>
          <div className="account-tags">
            {SEED_ACCOUNTS.map((acc) => (
              <Tag
                key={acc.username}
                className="account-tag"
                onClick={() =>
                  onFinish({ username: acc.username, password: acc.password })
                }
              >
                {acc.hint}：{acc.username} / {acc.password}
              </Tag>
            ))}
          </div>
        </div>
      </Card>
    </div>
  );
}
