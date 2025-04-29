import { useEffect, useState } from 'react'
import { DataTable, Column } from '../components/DataTable'
import api from '@/lib/axios'

type SystemLog = {
  LogID: number
  CreatedAt: string
  UserID?: number
  UserRole?: string
  Action: string
  Resource: string
  Status: string
  IPAddress?: string
  HTTPMethod: string
  Endpoint: string
  // StatusCode?:string
  // XForwardedFor?:string
  // UserAgent?:string
  // Referer?:string
  // Origin?:string
  // ClientApp?:string
  // BranchID?:string
  // Details?: Record<string, any>   
  // Metadata?: Record<string, any> 
}

const columns: Column<SystemLog>[] = [
  { header: 'Log ID', accessor: 'LogID' },
  { header: 'Time', accessor: 'CreatedAt' },
  { header: 'User', accessor: 'UserID' },
  { header: 'Role', accessor: 'UserRole' },
  { header: 'Action', accessor: 'Action' },
  { header: 'Resource', accessor: 'Resource' },
  { header: 'Status', accessor: 'Status' },
  { header: 'IP Address', accessor: 'IPAddress' },
  { header: 'Method', accessor: 'HTTPMethod' },
  { header: 'Endpoint', accessor: 'Endpoint' },
  // { header: 'StatusCode', accessor: 'StatusCode' },
  // { header: 'XForwardedFor', accessor: 'XForwardedFor' },
  // { header: 'UserAgent', accessor: 'UserAgent' },
  // { header: 'Referer', accessor: 'Referer' },
  // { header: 'Origin', accessor: 'Origin' },
  // { header: 'ClientApp', accessor: 'ClientApp' },
  // { header: 'BranchID', accessor: 'BranchID' },
  // { header: 'Details', accessor: 'Details' },
  // { header: 'Metadata', accessor: 'Metadata' },
]

export default function LogTablePage() {
  const [logs, setLogs] = useState<SystemLog[]>([])
  const [total, setTotal] = useState<number>(0)
  const [loading, setLoading] = useState<boolean>(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    api
      .get<{ logs: SystemLog[]; total: number }>('/admin/system_logs', {
        params: { page: 1, limit: 10 },
      })
      .then((res) => {
        setLogs(res.data.logs)
        setTotal(res.data.total)
      })
      .catch((err) => {
        setError(err.message)
      })
      .finally(() => {
        setLoading(false)
      })
  }, [])

  if (loading) return <p>Loading logs...</p>
  if (error) return <p className="text-red-600">Failed to load logs: {error}</p>

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-2">System Logs</h1>
      <p className="mb-4">Total logs: {total}</p>
      <DataTable<SystemLog>
        columns={columns}
        data={logs}
        onRowClick={(row) => console.log('Clicked log:', row)}
      />
    </div>
  )
}
