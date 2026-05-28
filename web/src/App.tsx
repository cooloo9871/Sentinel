import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { LoginPage } from './pages/LoginPage'
import { PolicyListPage } from './pages/PolicyListPage'
import { PolicyEditPage } from './pages/PolicyEditPage'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/" element={<PolicyListPage />} />
        <Route path="/policies/new" element={<PolicyEditPage />} />
        <Route path="/policies/:name" element={<PolicyEditPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}
