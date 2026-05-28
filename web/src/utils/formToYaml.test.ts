import { describe, it, expect } from 'vitest'
import { formToYaml } from './formToYaml'
import type { PolicyFormInput } from '../api/types'

describe('formToYaml', () => {
  it('generates process kprobe for cluster-wide policy', () => {
    const input: PolicyFormInput = {
      name: 'block-shells',
      process: [{ binaries: ['/bin/bash', '/bin/sh'] }],
    }
    const yaml = formToYaml(input, 'Post')
    expect(yaml).toContain('name: block-shells')
    expect(yaml).toContain('kind: TracingPolicy')
    expect(yaml).toContain('sys_execve')
    expect(yaml).toContain('Post')
    expect(yaml).toContain('/bin/bash')
  })

  it('sets kind to TracingPolicyNamespaced when namespace is set', () => {
    const input: PolicyFormInput = {
      name: 'ns-policy',
      namespace: 'production',
      process: [{ binaries: ['/bin/sh'] }],
    }
    const yaml = formToYaml(input, 'Post')
    expect(yaml).toContain('kind: TracingPolicyNamespaced')
    expect(yaml).toContain('namespace: production')
  })

  it('generates file kprobe for write operation', () => {
    const input: PolicyFormInput = {
      name: 'watch-etc',
      file: [{ paths: ['/etc/passwd'], operation: 'write' }],
    }
    const yaml = formToYaml(input, 'Post')
    expect(yaml).toContain('sys_write')
    expect(yaml).toContain('/etc/passwd')
  })

  it('generates network kprobe with CIDR and port', () => {
    const input: PolicyFormInput = {
      name: 'net-policy',
      network: [{ protocol: 'TCP', cidr: '10.0.0.0/8', port: 8080 }],
    }
    const yaml = formToYaml(input, 'Sigkill')
    expect(yaml).toContain('tcp_connect')
    expect(yaml).toContain('10.0.0.0/8:8080')
    expect(yaml).toContain('Sigkill')
  })

  it('generates network kprobe with CIDR only (no port)', () => {
    const input: PolicyFormInput = {
      name: 'cidr-only',
      network: [{ protocol: 'TCP', cidr: '192.168.0.0/16' }],
    }
    const yaml = formToYaml(input, 'Post')
    expect(yaml).toContain('192.168.0.0/16')
    expect(yaml).not.toContain('192.168.0.0/16:')
  })

  it('includes podSelector when provided', () => {
    const input: PolicyFormInput = {
      name: 'scoped',
      podSelector: { app: 'myapp' },
      process: [{ binaries: ['/bin/sh'] }],
    }
    const yaml = formToYaml(input, 'Post')
    expect(yaml).toContain('podSelector')
    expect(yaml).toContain('myapp')
  })

  it('returns multiple kprobes for mixed rules', () => {
    const input: PolicyFormInput = {
      name: 'multi',
      process: [{ binaries: ['/bin/bash'] }],
      file: [{ paths: ['/etc'], operation: 'open' }],
      network: [{ protocol: 'TCP', cidr: '0.0.0.0/0' }],
    }
    const yaml = formToYaml(input, 'Post')
    expect(yaml).toContain('sys_execve')
    expect(yaml).toContain('sys_openat')
    expect(yaml).toContain('tcp_connect')
  })
})
