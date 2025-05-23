import { Navigate } from "react-router-dom"
import { ReactNode } from "react"

export default function ProtectedRoute({ children }: { children: ReactNode }) {
  const token = localStorage.getItem("token")
  if (!token) {
    return <Navigate to="/" replace />
  }
  return children
}
