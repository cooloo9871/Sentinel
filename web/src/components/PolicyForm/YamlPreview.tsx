interface Props {
  yamlText: string
}

export function YamlPreview({ yamlText }: Props) {
  return (
    <pre
      style={{
        background: '#1e1e1e',
        color: '#d4d4d4',
        padding: 16,
        borderRadius: 4,
        overflow: 'auto',
        fontSize: 12,
        minHeight: 200,
        maxHeight: 500,
      }}
    >
      {yamlText || '# Fill in the form to preview YAML'}
    </pre>
  )
}
