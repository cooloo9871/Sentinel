import { Button, Form, Input, Space, Typography } from 'antd'
import { PlusOutlined, MinusCircleOutlined } from '@ant-design/icons'

const { Text } = Typography

interface Props {
  name: string
}

export function ProcessSection({ name }: Props) {
  return (
    <Form.List name={name}>
      {(fields, { add, remove }) => (
        <>
          <Text strong>Process Rules</Text>
          {fields.map(({ key, name: fieldName, ...rest }) => (
            <Space key={key} align="baseline" style={{ display: 'flex', marginBottom: 4 }}>
              <Form.Item
                {...rest}
                name={[fieldName, 'binaries', 0]}
                rules={[{ required: true, message: 'Enter binary path' }]}
              >
                <Input placeholder="/bin/bash" style={{ width: 300 }} />
              </Form.Item>
              <MinusCircleOutlined onClick={() => remove(fieldName)} />
            </Space>
          ))}
          <Button type="dashed" onClick={() => add()} icon={<PlusOutlined />} size="small">
            Add Process Rule
          </Button>
        </>
      )}
    </Form.List>
  )
}
