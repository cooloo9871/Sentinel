# Sentinel — 設計文件

**日期：** 2026-05-28  
**語言：** Go（後端）、React + TypeScript（前端）  
**狀態：** 已審核，待實作

---

## 1. 專案概述

Sentinel 是一個部署在 Kubernetes 叢集內的 Dashboard，讓使用者透過表單或 YAML 編輯器管理 Cilium Tetragon 的 `TracingPolicy` 與 `TracingPolicyNamespaced` CRD，無需直接操作 kubectl。

**核心功能：**
- 表單式建立 / 編輯 Policy（Process、File、Network 三種規則）
- 完整 YAML 編輯模式
- 全域模式切換：Monitoring（Post）或 Protect（Sigkill）
- 支援 Cluster-wide 與 Namespace-scoped Policy
- 基本帳號密碼認證

---

## 2. 整體架構

```
┌─────────────────────────────────────────────┐
│              Kubernetes Cluster              │
│                                             │
│  ┌──────────────┐     ┌──────────────────┐  │
│  │   Sentinel   │     │    Tetragon      │  │
│  │  (Pod)       │     │    (DaemonSet)   │  │
│  │              │     │                  │  │
│  │ ┌──────────┐ │     │  TracingPolicy   │  │
│  │ │Go Server │ │────▶│  TracingPolicy   │  │
│  │ │:8080     │ │     │  Namespaced      │  │
│  │ └────┬─────┘ │     └──────────────────┘  │
│  │      │       │                           │
│  │ ┌────▼─────┐ │                           │
│  │ │  React   │ │                           │
│  │ │  (embed) │ │                           │
│  │ └──────────┘ │                           │
│  └──────────────┘                           │
└─────────────────────────────────────────────┘
```

**關鍵決策：**
- Go Server 同時 serve React 靜態檔案（`embed.FS`），單一 Docker image
- 透過 in-cluster ServiceAccount + RBAC 存取 k8s API
- 認證憑證儲存在 Kubernetes Secret（bcrypt hash），不需外部 DB
- Kubernetes CRD 本身作為 Policy 的唯一 storage

---

## 3. 後端 API（Go）

### 目錄結構

```
sentinel/
├── cmd/server/main.go
├── internal/
│   ├── auth/          # JWT + bcrypt
│   ├── handler/       # HTTP handlers
│   ├── k8s/           # client-go wrapper
│   └── policy/        # TracingPolicy struct 與 builder
├── web/               # React build output（embed）
└── Dockerfile
```

### REST API

**認證**
```
POST /api/auth/login     # 驗證帳號密碼，回傳 JWT
POST /api/auth/logout
```

**Policy 管理**
```
GET    /api/policies            # 列出所有 policies（cluster + namespaced）
POST   /api/policies            # 建立 policy
GET    /api/policies/:name      # 取得單一 policy
PUT    /api/policies/:name      # 更新 policy
DELETE /api/policies/:name      # 刪除 policy
GET    /api/namespaces          # 列出可用 namespace
```

**模式切換**
```
GET    /api/mode                # 取得目前模式
PUT    /api/mode                # 切換模式（批次更新所有 policy 的 action）
```

---

## 4. 前端（React + TypeScript + Ant Design）

### 頁面結構

**Policy 列表頁（首頁）**
- 表格顯示：名稱、類型、範圍（cluster/namespaced）、namespace、建立時間
- 右上角全域模式 Toggle（Monitoring ↔ Protect）
- 每筆 policy：編輯 / 刪除按鈕

**Policy 編輯頁**

上方 Tab 切換兩種編輯模式：

```
[ Form 編輯 ]  [ YAML 編輯 ]
```

**Form 編輯 Tab：**

| 規則類型 | 可設定欄位 |
|---------|-----------|
| Process | 監控的 binary 路徑、arguments 條件 |
| File | 路徑（支援 glob）、操作類型（read/write/open） |
| Network | protocol（TCP/UDP）、port、CIDR |

每個 section 可新增多筆規則。右側 panel 即時預覽對應的完整 YAML。

**YAML 編輯 Tab：**
- Monaco Editor 直接編輯完整 YAML
- YAML parse 語法驗證
- Submit 直接套用至叢集

---

## 5. 資料模型

### Policy 頂層 Metadata

```yaml
metadata:
  name: <policy-name>
  namespace: <namespace>   # 僅 namespaced policy
spec:
  podSelector:             # 空白 = 套用全部 Pod
    matchLabels:
      app: my-app
```

### Process 規則（kprobe on sys_execve）

```yaml
kprobes:
- call: "sys_execve"
  syscall: true
  args:
  - index: 0
    type: "string"
  selectors:
  - matchBinaries:
    - operator: "In"
      values: ["/bin/bash", "/bin/sh"]
    matchActions:
    - action: Post      # Monitoring 模式
    # action: Sigkill   # Protect 模式
```

### File 規則（kprobe on sys_write / sys_read）

```yaml
kprobes:
- call: "sys_write"
  syscall: true
  selectors:
  - matchArgs:
    - index: 0
      operator: "Prefix"
      values: ["/etc/passwd"]
    matchActions:
    - action: Post
```

### Network 規則（kprobe on tcp_connect）

```yaml
kprobes:
- call: "tcp_connect"
  syscall: false
  selectors:
  - matchArgs:
    - index: 0
      operator: "Equal"
      values: ["10.0.0.0/8:8080"]
    matchActions:
    - action: Post
```

### 模式對應

| Dashboard 模式 | Tetragon action | 行為 |
|--------------|----------------|------|
| Monitoring | `Post` | 記錄事件，不阻擋 |
| Protect | `Sigkill` | 終止違規 process |

切換模式時，後端批次將所有 policy 的 `action` 欄位替換為對應值。

---

## 6. 部署（raw YAML + Kustomize）

### 目錄結構

```
deploy/
├── base/
│   ├── kustomization.yaml
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── serviceaccount.yaml
│   ├── clusterrole.yaml
│   ├── clusterrolebinding.yaml
│   └── secret.yaml          # 範本，實際值由 overlay 覆蓋
└── overlays/
    └── production/
        ├── kustomization.yaml
        └── secret-patch.yaml  # 設定實際帳號密碼
```

### 部署指令

```bash
kubectl apply -k deploy/overlays/production
```

### RBAC 所需權限

```yaml
rules:
- apiGroups: ["cilium.io"]
  resources: ["tracingpolicies", "tracingpoliciesnamespaced"]
  verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]
- apiGroups: [""]
  resources: ["namespaces"]
  verbs: ["get", "list"]
- apiGroups: [""]
  resources: ["secrets"]
  resourceNames: ["sentinel-credentials"]
  verbs: ["get"]
```

---

## 7. 測試策略

| 層級 | 工具 | 範圍 |
|------|------|------|
| 後端單元測試 | `go test` | policy builder（YAML 產生邏輯）、auth（JWT/bcrypt） |
| 後端整合測試 | `envtest`（controller-runtime） | k8s API CRUD，不需真實叢集 |
| 前端單元測試 | Vitest | Form → YAML 轉換邏輯 |
| E2E | 手動 or kind cluster | 完整流程驗證 |

最關鍵的測試：policy builder 單元測試，給定 Form 欄位輸入，驗證產生的 YAML 與預期 TracingPolicy spec 完全吻合。
