import { Button, Form, Input, Select, Space, Typography } from 'antd'
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons'

const { Text } = Typography

interface Props {
  name: string
}

export function FileSection({ name }: Props) {
  return (
    <Form.List name={name}>
      {(fields, { add, remove }) => (
        <>
          <Text strong>File Rules</Text>
          {fields.map(({ key, name: fieldName, ...rest }) => (
            <Space key={key} align="baseline" style={{ display: 'flex', marginBottom: 4 }}>
              <Form.Item
                {...rest}
                name={[fieldName, 'paths', 0]}
                rules={[{ required: true, message: 'Enter file path' }]}
              >
                <Input placeholder="/etc/passwd" style={{ width: 250 }} />
              </Form.Item>
              <Form.Item
                {...rest}
                name={[fieldName, 'operation']}
                rules={[{ required: true, message: 'Select operation' }]}
              >
                <Select placeholder="Operation" style={{ width: 110 }}>
                  <Select.Option value="read">read</Select.Option>
                  <Select.Option value="write">write</Select.Option>
                  <Select.Option value="open">open</Select.Option>
                </Select>
              </Form.Item>
              <MinusCircleOutlined onClick={() => remove(fieldName)} />
            </Space>
          ))}
          <Button type="dashed" onClick={() => add()} icon={<PlusOutlined />} size="small">
            Add File Rule
          </Button>
        </>
      )}
    </Form.List>
  )
}
