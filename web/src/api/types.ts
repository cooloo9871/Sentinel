export interface ProcessRule {
  binaries: string[]
}

export interface FileRule {
  paths: string[]
  operation: 'read' | 'write' | 'open'
}

export interface NetworkRule {
  protocol: 'TCP' | 'UDP'
  cidr: string
  port?: number
}

export interface PolicyFormInput {
  name: string
  namespace?: string
  podSelector?: Record<string, string>
  process?: ProcessRule[]
  file?: FileRule[]
  network?: NetworkRule[]
}

export interface PolicyRecord {
  name: string
  namespace?: string
  scope: 'cluster' | 'namespaced'
  createdAt: string
  rawYaml: string
}

export type Mode = 'Monitoring' | 'Protect' | 'Mixed'

export interface CreatePolicyPayload {
  source: 'form' | 'yaml'
  form?: PolicyFormInput
  action?: string
  rawYaml?: string
}
