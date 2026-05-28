import { useState } from 'react'
import MonacoEditor from '@monaco-editor/react'
import { Button, Alert, Space } from 'antd'
import yaml from 'js-yaml'

interface Props {
  initialValue?: string
  onSubmit: (yamlText: string) => Promise<void>
  loading?: boolean
}

export function YamlEditor({ initialValue = '', onSubmit, loading }: Props) {
  const [value, setValue] = useState(initialValue)
  const [error, setError] = useState('')

  const handleSubmit = async () => {
    setError('')
    try {
      yaml.load(value)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Invalid YAML')
      return
    }
    try {
      await onSubmit(value)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to apply YAML')
    }
  }

  return (
    <Space direction="vertical" style={{ width: '100%' }}>
      {error && <Alert type="error" message={error} showIcon />}
      <MonacoEditor
        height="500px"
        language="yaml"
        theme="vs-dark"
        value={value}
        onChange={(v) => setValue(v ?? '')}
        options={{ minimap: { enabled: false }, fontSize: 13 }}
      />
      <Button type="primary" onClick={handleSubmit} loading={loading}>
        Apply YAML
      </Button>
    </Space>
  )
}
