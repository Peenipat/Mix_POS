export type Column<T> = {
  header: string;
  accessor: keyof T | ((row: T, rowIndex: number) => any);
};

export type Action<T> = {
  label: string;
  onClick: (row: T) => void;
  className?: string; // optional custom styling
};

type DataTableProps<T> = {
  data: T[];
  columns: Column<T>[];
  onRowClick?: (row: T) => void;
  showEdit?: boolean;
  showDelete?: boolean;
  onEdit?: (row: T) => void;
  onDelete?: (row: T) => void;
  actions?: Action<T>[];
  page?: number;
  perPage?: number;
  total?: number;
  onPageChange?: (page: number) => void;
};

export function DataTable<T extends Record<string, any>>({
  data,
  columns,
  onRowClick,
  showEdit = true,
  showDelete = true,
  onEdit,
  onDelete,
  actions = [],
  page,
  perPage,
  total,
  onPageChange,
}: DataTableProps<T>) {
  const hasAnyAction = showEdit || showDelete || actions.length > 0;

  function getPaginationPages(current: number, totalPages: number): (number | "...")[] {
    const pages: (number | "...")[] = [];

    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      if (current <= 4) {
        pages.push(1, 2, 3, 4, 5, "...", totalPages);
      } else if (current >= totalPages - 3) {
        pages.push(1, "...", totalPages - 4, totalPages - 3, totalPages - 2, totalPages - 1, totalPages);
      } else {
        pages.push(1, "...", current - 1, current, current + 1, "...", totalPages);
      }
    }

    return pages;
  }

  return (
    <div className="relative overflow-x-auto shadow-md sm:rounded-lg">
      <table className="w-full text-sm text-left text-gray-500 dark:text-gray-400">
        <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
          <tr>
            {columns.map((col) => (
              <th
                key={String(col.accessor)}
                scope="col"
                className="px-6 py-3"
              >
                {col.header}
              </th>
            ))}
            {hasAnyAction && (
              <th scope="col" className="px-6 py-3 text-right">
                <span className="sr-only">Actions</span>
              </th>
            )}
          </tr>
        </thead>
        <tbody>
          {data.map((row, rowIndex) => (
            <tr
              key={rowIndex}
              className="bg-white border-b dark:bg-gray-800 dark:border-gray-700 border-gray-200 hover:bg-gray-50"
              onClick={() => onRowClick && onRowClick(row)}
            >
              {columns.map((col) => {
                // 1) แยกเคสว่า accessor เป็น string key หรือ function
                let rawValue: any;
                if (typeof col.accessor === "function") {
                  rawValue = col.accessor(row, rowIndex);
                } else {
                  rawValue = (row as any)[col.accessor];
                }

                // 2) แปลง rawValue เป็น string เพื่อนำไปแสดง
                let display = "";
                if (rawValue === null || rawValue === undefined) {
                  display = "";
                } else if (typeof rawValue === "object") {
                  try {
                    display = JSON.stringify(rawValue);
                  } catch {
                    display = String(rawValue);
                  }
                } else {
                  display = String(rawValue);
                }

                return (
                  <td
                    key={
                      typeof col.accessor === "function"
                        ? `func-${rowIndex}`
                        : String(col.accessor)
                    }
                    className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
                  >
                    {display}
                  </td>
                );
              })}

              {hasAnyAction && (
                <td className="px-6 py-4 text-right space-x-2">
                  {actions.map((action, idx) => (
                    <button
                      key={idx}
                      onClick={(e) => {
                        e.stopPropagation();
                        action.onClick(row);
                      }}
                      className={`font-medium hover:underline ${action.className ?? ""}`}
                    >
                      {action.label}
                    </button>
                  ))}
                  {showEdit && onEdit && (
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        onEdit(row);
                      }}
                      className="font-medium text-blue-600 dark:text-blue-500 hover:underline"
                    >
                      แก้ไข
                    </button>
                  )}
                  {showDelete && onDelete && (
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        onDelete(row);
                      }}
                      className="font-medium text-red-600 dark:text-red-500 hover:underline"
                    >
                      ลบ
                    </button>
                  )}
                </td>
              )}
            </tr>
          ))}

        </tbody>

      </table>
      {typeof page === "number" && typeof total === "number" && typeof perPage === "number" && (
        <div className="flex items-center justify-between px-6 py-3 bg-white border-t">
          <span className="text-sm text-gray-600">
            แสดง {Math.min((page - 1) * perPage + 1, total)} - {Math.min(page * perPage, total)} จากทั้งหมด {total} รายการ
          </span>

          <div className="flex items-center gap-1">
            <button
              disabled={page === 1}
              onClick={() => onPageChange?.(page - 1)}
              className="px-3 py-1 text-sm border rounded disabled:opacity-50"
            >
              ก่อนหน้า
            </button>

            {getPaginationPages(page, Math.ceil(total / perPage)).map((p, i) =>
              p === "..." ? (
                <span key={i} className="px-2 py-1 text-gray-400">...</span>
              ) : (
                <button
                  key={p}
                  onClick={() => onPageChange?.(p as number)}
                  className={`px-3 py-1 text-sm border rounded ${p === page ? "bg-gray-200 font-semibold" : ""
                    }`}
                >
                  {p}
                </button>
              )
            )}

            <button
              disabled={page * perPage >= total}
              onClick={() => onPageChange?.(page + 1)}
              className="px-3 py-1 text-sm border rounded disabled:opacity-50"
            >
              ถัดไป
            </button>
          </div>
        </div>
      )}


    </div>
  );
}