import { useEffect } from 'react'
import { Form, Input, Select, Row, Col, Divider, Button, Space } from 'antd'
import { ProcessSection } from './ProcessSection'
import { FileSection } from './FileSection'
import { NetworkSection } from './NetworkSection'
import { YamlPreview } from './YamlPreview'
import { formToYaml } from '../../utils/formToYaml'
import type { PolicyFormInput } from '../../api/types'

interface Props {
  initialValues?: Partial<PolicyFormInput>
  namespaces: string[]
  action: string
  onSubmit: (values: PolicyFormInput) => Promise<void>
  loading?: boolean
}

export function PolicyForm({ initialValues, namespaces, action, onSubmit, loading }: Props) {
  const [form] = Form.useForm<PolicyFormInput>()
  const formValues = Form.useWatch([], form)

  const yamlPreview = (() => {
    try {
      if (!formValues?.name) return ''
      return formToYaml(formValues as PolicyFormInput, action)
    } catch {
      return ''
    }
  })()

  useEffect(() => {
    if (initialValues) form.setFieldsValue(initialValues)
  }, [initialValues, form])

  return (
    <Row gutter={24}>
      <Col span={14}>
        <Form form={form} layout="vertical" onFinish={onSubmit}>
          <Form.Item name="name" label="Policy Name" rules={[{ required: true }]}>
            <Input placeholder="my-policy" />
          </Form.Item>
          <Form.Item name="namespace" label="Namespace (leave empty for cluster-wide)">
            <Select
              allowClear
              placeholder="cluster-wide"
              options={namespaces.map((n) => ({ value: n, label: n }))}
            />
          </Form.Item>
          <Divider />
          <ProcessSection name="process" />
          <Divider />
          <FileSection name="file" />
          <Divider />
          <NetworkSection name="network" />
          <Divider />
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" loading={loading}>
                Apply Policy
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Col>
      <Col span={10}>
        <div style={{ position: 'sticky', top: 24 }}>
          <strong>YAML Preview</strong>
          <YamlPreview yamlText={yamlPreview} />
        </div>
      </Col>
    </Row>
  )
}
