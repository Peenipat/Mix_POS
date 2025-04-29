// Generic types for columns and data
export interface Column<T> { //กำหนดให้ column ให้รับข้อมูลได้ทุกประเภท
  header: string // table name
  accessor: keyof T // key
}

export interface DataTableProps<T> {
  data: T[] // array ข้อมูลที่แสดง
  columns: Column<T>[] // กำหนด key และ haeder column
  onRowClick?: (row: T) => void
}

// Generic DataTable component using Tailwind CSS + DaisyUI
export function DataTable<T extends Record<string, any>>({ data, columns, onRowClick }: DataTableProps<T>) {
  return (
    <div className="overflow-x-auto w-full">
      <table className="table table-zebra w-full">
        <thead>
          <tr>
            {columns.map((col) => (
              <th key={String(col.accessor)}>{col.header}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.map((row, rowIndex) => (
            <tr
              key={rowIndex}
              className={onRowClick ? 'cursor-pointer' : ''}
              onClick={() => onRowClick && onRowClick(row)}
            >
              {columns.map((col) => {
                const value = row[col.accessor]   // ค่า raw จากข้อมูล
                let display: string
                if (value === null || value === undefined) {
                  display = ''
                } else if (typeof value === 'object') {
                  try {
                    display = JSON.stringify(value)
                  } catch {
                    display = String(value)
                  }
                } else {
                  display = String(value)
                }
                return (
                  <td key={String(col.accessor)} className="capitalize">
                    {display}
                  </td>
                )
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
