import { useEffect, useState } from 'react'
import { Switch, Space, Typography, Tag, Tooltip, Alert } from 'antd'
import { modeApi } from '../api/client'
import type { Mode } from '../api/types'

const { Text } = Typography

interface Props {
  onModeChange?: (mode: Mode) => void
}

export function ModeToggle({ onModeChange }: Props) {
  const [mode, setMode] = useState<Mode>('Monitoring')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    modeApi.get().then(setMode)
  }, [])

  const handleChange = async (checked: boolean) => {
    const next = checked ? 'Protect' : 'Monitoring'
    setLoading(true)
    try {
      await modeApi.set(next)
      setMode(next)
      onModeChange?.(next)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Space direction="vertical" size="small">
      <Space>
        <Text>Monitoring</Text>
        <Tooltip title={mode === 'Mixed' ? 'Policies have mixed actions. Switch to apply uniformly.' : ''}>
          <Switch
            checked={mode === 'Protect'}
            onChange={handleChange}
            loading={loading}
            checkedChildren="Protect"
            unCheckedChildren="Monitor"
          />
        </Tooltip>
        <Text>Protect</Text>
        {mode === 'Protect' && <Tag color="red">PROTECT</Tag>}
        {mode === 'Monitoring' && <Tag color="green">MONITORING</Tag>}
      </Space>
      {mode === 'Mixed' && (
        <Alert
          type="warning"
          message="Policies have mixed actions. Switch mode to apply uniformly."
          showIcon
        />
      )}
    </Space>
  )
}
