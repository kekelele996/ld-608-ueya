import { DatePicker, Form, Modal, Select, Alert } from "antd";
import dayjs from "dayjs";
import { useEffect } from "react";
import type { ResourceBooking } from "../../types/ResourceBooking";
import type { GroundResource } from "../../types/GroundResource";
import { RESOURCE_STATUS_TEXT } from "../../constants/ResourceStatus";
import { toApiTime } from "../../utils/formatters";

// AdjustBookingModal：资源调度页 / 过站详情共用的"调整时段"弹窗。
// 保存后无论成功还是冲突都刷新同一结果（CONFIRMED / PENDING+原因）。
export function AdjustBookingModal({
  open,
  booking,
  resources,
  requiredType,
  onClose,
  onSubmit
}: {
  open: boolean;
  booking: ResourceBooking | null;
  resources: GroundResource[];
  // 若预约关联了任务，可传入任务资源类型，只展示同类资源；资源页全量调整时不传。
  requiredType?: string;
  onClose: () => void;
  onSubmit: (values: { resource_id: number; start_time: string; end_time: string }) => Promise<void>;
}) {
  const [form] = Form.useForm();

  useEffect(() => {
    if (booking) {
      form.setFieldsValue({
        resource_id: booking.resource_id,
        range: [dayjs(booking.start_time), dayjs(booking.end_time)]
      });
    }
  }, [booking, form]);

  return (
    <Modal
      title={`调整预约时段 #${booking?.id ?? ""}`}
      open={open}
      onCancel={onClose}
      onOk={async () => {
        const values = await form.validateFields();
        await onSubmit({
          resource_id: values.resource_id,
          start_time: toApiTime(values.range[0]),
          end_time: toApiTime(values.range[1])
        });
      }}
      destroyOnClose
    >
      {booking?.conflict_reason && (
        <Alert type="warning" showIcon className="modal-alert" message={booking.conflict_reason} />
      )}
      <Form form={form} layout="vertical">
        <Form.Item name="resource_id" label="资源" rules={[{ required: true }]}>
          <Select
            options={resources
              .filter((r) => !requiredType || r.resource_type === requiredType)
              .map((r) => ({
                value: r.id,
                disabled: r.availability_status !== "AVAILABLE",
                label: `${r.resource_code}（${RESOURCE_STATUS_TEXT[r.availability_status] ?? r.availability_status}）`
              }))}
          />
        </Form.Item>
        <Form.Item name="range" label="预约时段" rules={[{ required: true }]}>
          <DatePicker.RangePicker showTime={{ format: "HH:mm" }} format="MM-DD HH:mm" />
        </Form.Item>
      </Form>
    </Modal>
  );
}
