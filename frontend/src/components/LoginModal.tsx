import React, { useState } from "react";
import { Modal, Card, Button, Form, Select, message, Typography } from "antd";
import { login as loginApi } from "../api/auth";
import { setSession } from "../stores/authSlice";
import { useAppDispatch } from "../stores/hooks";
import { ROLE_USERNAMES, ROLE_TEXT } from "../types/auth";

// LoginModal implements the local JWT login + role switch. No passwords in
// the local seed; picking a user is the RBAC switch for demos/reviews.
export const LoginModal: React.FC = () => {
  const dispatch = useAppDispatch();
  const [open, setOpen] = useState(true);
  const [loading, setLoading] = useState(false);
  const [username, setUsername] = useState(ROLE_USERNAMES[0].username);

  const submit = async () => {
    setLoading(true);
    try {
      const res = await loginApi(username);
      dispatch(
        setSession({
          token: res.token,
          username: res.username,
          displayName: res.display_name,
          role: res.role,
          teamId: res.team_id,
        }),
      );
      message.success(`已以「${res.display_name}」身份登录`);
      setOpen(false);
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Modal
      open={open}
      closable={false}
      maskClosable={false}
      footer={null}
      title="选择登录角色（JWT + RBAC）"
      width={460}
    >
      <Card size="small">
        <Form layout="vertical">
          <Form.Item label="本地账号">
            <Select
              value={username}
              onChange={setUsername}
              options={ROLE_USERNAMES.map((u) => ({
                value: u.username,
                label: `${u.username}（${ROLE_TEXT[u.role]}${u.team ? " · " + u.team : ""}）`,
              }))}
            />
          </Form.Item>
          <Typography.Paragraph type="secondary" style={{ marginBottom: 12 }}>
            地勤调度可登记到站、生成保障计划、登记延误；资源管理员可调整待处理预约与维护窗口；
            保障班组可签收/完成任务；运行督导拥有全部权限。
          </Typography.Paragraph>
          <Button type="primary" block loading={loading} onClick={submit}>
            登录
          </Button>
        </Form>
      </Card>
    </Modal>
  );
};
