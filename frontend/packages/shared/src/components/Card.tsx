export interface CardProps {
  onView?: () => void;
  onEdit?: () => void;
  onDelete?: () => void;
  showActions?: boolean;
  className?: string;
  children: React.ReactNode;
}

export const Card: React.FC<CardProps> = ({
  children,
  onView,
  onEdit,
  onDelete,
  showActions = true,
  className = "",
}) => {
  return (
    <div className={`rounded-lg shadow-md bg-white p-4 flex flex-col ${className}`}>
      {/* เนื้อหาหลัก */}
      <div className="flex-1">{children}</div>

      {/* ปุ่ม Actions */}
      {showActions && (() => {
        const buttons = [
          onView ? "view" : null,
          onEdit ? "edit" : null,
          onDelete ? "delete" : null,
        ].filter(Boolean);

        const isTwoButtons = buttons.length === 2;
        const isThreeButtons = buttons.length === 3;

        return (
          <div
            className={`mt-4 flex ${isTwoButtons ? "justify-between" : "justify-between"
              } items-center`}
          >
            {/* ปุ่ม: รายละเอียด */}
            {onView && (
              <button
                className="text-sm text-green-600 underline"
                onClick={onView}
              >
                รายละเอียด
              </button>
            )}

            {/* ปุ่ม: แก้ไข */}
            {onEdit && (
              <button
                className={`text-sm text-yellow-500 underline ${isTwoButtons && !onView ? "ml-0" : ""
                  }`}
                onClick={onEdit}
              >
                แก้ไข
              </button>
            )}

            {/* ปุ่ม: ลบ */}
            {onDelete && (
              <button
                className={`text-sm text-red-500 underline ${isTwoButtons && !onView ? "ml-auto" : ""
                  }`}
                onClick={onDelete}
              >
                ลบ
              </button>
            )}
          </div>
        );
      })()}

    </div>
  );
};
