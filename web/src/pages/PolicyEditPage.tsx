import { useEffect, useState } from 'react'
import { useNavigate, useParams, useSearchParams } from 'react-router-dom'
import { Layout, Tabs, Typography, Breadcrumb, message } from 'antd'
import { policyApi, namespaceApi, modeApi } from '../api/client'
import { PolicyForm } from '../components/PolicyForm/PolicyForm'
import { YamlEditor } from '../components/YamlEditor'
import type { PolicyFormInput, Mode } from '../api/types'

const { Header, Content } = Layout
const { Title } = Typography

export function PolicyEditPage() {
  const { name } = useParams<{ name: string }>()
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const namespace = searchParams.get('namespace') || undefined
  const isNew = !name || name === 'new'

  const [namespaces, setNamespaces] = useState<string[]>([])
  const [mode, setMode] = useState<Mode>('Monitoring')
  const [initialYaml, setInitialYaml] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    namespaceApi.list().then(setNamespaces)
    modeApi.get().then(setMode)
    if (!isNew && name) {
      policyApi.get(name, namespace).then((r) => setInitialYaml(r.rawYaml))
    }
  }, [name, namespace, isNew])

  const action = mode === 'Protect' ? 'Sigkill' : 'Post'

  const handleFormSubmit = async (values: PolicyFormInput) => {
    setLoading(true)
    try {
      const payload = { source: 'form' as const, form: values, action }
      if (isNew) {
        await policyApi.create(payload)
      } else {
        await policyApi.update(name!, payload)
      }
      message.success('Policy applied')
      navigate('/')
    } catch {
      message.error('Failed to apply policy')
    } finally {
      setLoading(false)
    }
  }

  const handleYamlSubmit = async (rawYaml: string) => {
    setLoading(true)
    try {
      const payload = { source: 'yaml' as const, rawYaml }
      if (isNew) {
        await policyApi.create(payload)
      } else {
        await policyApi.update(name!, payload)
      }
      message.success('Policy applied')
      navigate('/')
    } catch {
      message.error('Failed to apply YAML')
    } finally {
      setLoading(false)
    }
  }

  const items = [
    {
      key: 'form',
      label: 'Form Edit',
      children: (
        <PolicyForm
          namespaces={namespaces}
          action={action}
          onSubmit={handleFormSubmit}
          loading={loading}
        />
      ),
    },
    {
      key: 'yaml',
      label: 'YAML Edit',
      children: (
        <YamlEditor
          initialValue={initialYaml}
          onSubmit={handleYamlSubmit}
          loading={loading}
        />
      ),
    },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header>
        <Title level={4} style={{ color: '#fff', margin: 0 }}>Sentinel</Title>
      </Header>
      <Content style={{ padding: 24 }}>
        <Breadcrumb
          items={[
            { title: <a onClick={() => navigate('/')}>Policies</a> },
            { title: isNew ? 'New Policy' : name },
          ]}
          style={{ marginBottom: 16 }}
        />
        <Tabs defaultActiveKey="form" items={items} />
      </Content>
    </Layout>
  )
}
