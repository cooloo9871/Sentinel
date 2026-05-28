import { useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Layout, Button, Table, Space, Popconfirm, Typography, message } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { policyApi, authApi } from '../api/client'
import { ModeToggle } from '../components/ModeToggle'
import type { PolicyRecord } from '../api/types'

const { Header, Content } = Layout
const { Title } = Typography

export function PolicyListPage() {
  const navigate = useNavigate()
  const [policies, setPolicies] = useState<PolicyRecord[]>([])
  const [loading, setLoading] = useState(true)

  const fetchPolicies = async () => {
    setLoading(true)
    try {
      const data = await policyApi.list()
      setPolicies(data)
    } catch {
      message.error('Failed to load policies')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchPolicies() }, [])

  const handleDelete = async (name: string, namespace?: string) => {
    try {
      await policyApi.delete(name, namespace)
      message.success('Policy deleted')
      fetchPolicies()
    } catch {
      message.error('Failed to delete policy')
    }
  }

  const handleLogout = async () => {
    await authApi.logout()
    navigate('/login')
  }

  const columns = [
    { title: 'Name', dataIndex: 'name', key: 'name' },
    { title: 'Scope', dataIndex: 'scope', key: 'scope' },
    { title: 'Namespace', dataIndex: 'namespace', key: 'namespace', render: (v: string) => v || '—' },
    { title: 'Created', dataIndex: 'createdAt', key: 'createdAt' },
    {
      title: 'Actions',
      key: 'actions',
      render: (_: unknown, record: PolicyRecord) => (
        <Space>
          <Button size="small" onClick={() => navigate(`/policies/${record.name}?namespace=${record.namespace ?? ''}`)}>
            Edit
          </Button>
          <Popconfirm
            title="Delete this policy?"
            onConfirm={() => handleDelete(record.name, record.namespace)}
            okText="Delete"
            okButtonProps={{ danger: true }}
          >
            <Button size="small" danger>Delete</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Title level={4} style={{ color: '#fff', margin: 0 }}>Sentinel</Title>
        <Space>
          <ModeToggle />
          <Button onClick={handleLogout} style={{ marginLeft: 16 }}>Logout</Button>
        </Space>
      </Header>
      <Content style={{ padding: 24 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
          <Title level={5} style={{ margin: 0 }}>TracingPolicies</Title>
          <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate('/policies/new')}>
            New Policy
          </Button>
        </div>
        <Table
          dataSource={policies}
          columns={columns}
          rowKey={(r) => `${r.scope}-${r.namespace ?? ''}-${r.name}`}
          loading={loading}
        />
      </Content>
    </Layout>
  )
}
