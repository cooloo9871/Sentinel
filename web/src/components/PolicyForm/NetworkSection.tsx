import { Button, Form, Input, InputNumber, Select, Space, Typography } from 'antd'
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons'

const { Text } = Typography

interface Props {
  name: string
}

export function NetworkSection({ name }: Props) {
  return (
    <Form.List name={name}>
      {(fields, { add, remove }) => (
        <>
          <Text strong>Network Rules</Text>
          {fields.map(({ key, name: fieldName, ...rest }) => (
            <Space key={key} align="baseline" style={{ display: 'flex', marginBottom: 4 }}>
              <Form.Item {...rest} name={[fieldName, 'protocol']}>
                <Select placeholder="TCP/UDP" style={{ width: 90 }}>
                  <Select.Option value="TCP">TCP</Select.Option>
                  <Select.Option value="UDP">UDP</Select.Option>
                </Select>
              </Form.Item>
              <Form.Item
                {...rest}
                name={[fieldName, 'cidr']}
                rules={[{ required: true, message: 'Enter CIDR' }]}
              >
                <Input placeholder="10.0.0.0/8" style={{ width: 160 }} />
              </Form.Item>
              <Form.Item {...rest} name={[fieldName, 'port']}>
                <InputNumber placeholder="Port" min={1} max={65535} style={{ width: 90 }} />
              </Form.Item>
              <MinusCircleOutlined onClick={() => remove(fieldName)} />
            </Space>
          ))}
          <Button type="dashed" onClick={() => add()} icon={<PlusOutlined />} size="small">
            Add Network Rule
          </Button>
        </>
      )}
    </Form.List>
  )
}
