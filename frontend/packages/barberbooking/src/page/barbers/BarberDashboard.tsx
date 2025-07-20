import { useState, useEffect } from "react";
import {
  Clock,
  CheckCircle,
  XCircle,
  Loader,
  CalendarDays,
} from "lucide-react";
import {
  AppointmentBrief,
  getAppointmentsByBarber,
} from "../../api/appointment";

// 🔁 แปลง status จาก backend → ไทย
function translateStatus(status: string): string {
  switch (status) {
    case "PENDING":
    case "CONFIRMED":
      return "รอให้บริการ";
    case "COMPLETED":
      return "เสร็จแล้ว";
    case "CANCELLED":
    case "NO_SHOW":
    case "RESCHEDULED":
      return "ยกเลิก";
    default:
      return "กำลังให้บริการ";
  }
}

const statusIcon: Record<string, JSX.Element> = {
  "รอให้บริการ": <Clock className="w-4 h-4 text-yellow-600 mr-1" />,
  "กำลังให้บริการ": (
    <Loader className="w-4 h-4 text-blue-600 mr-1 animate-spin" />
  ),
  "เสร็จแล้ว": <CheckCircle className="w-4 h-4 text-green-600 mr-1" />,
  "ยกเลิก": <XCircle className="w-4 h-4 text-red-600 mr-1" />,
};

export default function BarberDashboard() {
  const [selectedDate, setSelectedDate] = useState(
    new Date().toISOString().slice(0, 10)
  );
  const [appointments, setAppointments] = useState<AppointmentBrief[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const BARBER_ID = 1;

  useEffect(() => {
    async function fetchAppointments() {
      setIsLoading(true);
      try {
        const resp = await getAppointmentsByBarber(BARBER_ID, {
          start: selectedDate,
          end: selectedDate,
        });
        setAppointments(resp.data ?? []);
      } catch (error) {
        console.error("โหลดข้อมูลล้มเหลว", error);
        setAppointments([]);
      } finally {
        setIsLoading(false);
      }
    }

    fetchAppointments();
  }, [selectedDate]);

  const updateStatus = (id: number, newStatus: string) => {
    setAppointments((prev) =>
      prev.map((a) =>
        a.id === id ? { ...a, status: newStatus.toUpperCase() } : a
      )
    );
  };

  return (
    <div className="p-6 space-y-6">
      <h2 className="text-2xl font-bold">ตารางนัดหมาย</h2>

      {/* เลือกวัน */}
      <div className="flex items-center gap-2">
        <CalendarDays className="w-5 h-5 text-gray-500" />
        <input
          type="date"
          className="border p-2 rounded"
          value={selectedDate}
          onChange={(e) => setSelectedDate(e.target.value)}
        />
      </div>

      {/* Summary */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        <SummaryCard label="จำนวนคิวทั้งหมด" count={appointments.length} />
        <SummaryCard
          label="เสร็จแล้ว"
          count={
            appointments.filter((a) => translateStatus(a.status) === "เสร็จแล้ว").length
          }
          color="green"
        />
        <SummaryCard
          label="ยกเลิก"
          count={
            appointments.filter((a) => translateStatus(a.status) === "ยกเลิก").length
          }
          color="red"
        />
        <SummaryCard
          label="กำลังให้บริการ"
          count={
            appointments.filter((a) => translateStatus(a.status) === "กำลังให้บริการ").length
          }
          color="blue"
        />
      </div>

      {/* Appointments */}
      <ul className="space-y-4">
        {isLoading ? (
          <p className="text-center text-gray-500">กำลังโหลด...</p>
        ) : appointments.length === 0 ? (
          <p className="text-center text-gray-500">ไม่มีคิวในวันนี้</p>
        ) : (
          appointments.map((appt) => {
            const translated = translateStatus(appt.status);
            return (
              <li
                key={appt.id}
                className="border rounded p-4 shadow-sm bg-white flex justify-between items-center"
              >
                <div>
                  <div className="text-sm text-gray-500">
                    {appt.start} - {appt.end}
                  </div>
                  <div className="text-lg font-semibold">
                    {appt.customer.name}
                  </div>
                  <div className="text-sm text-gray-700">
                    {appt.service.name}
                  </div>
                </div>

                <div className="text-right space-y-1">
                  <div
                    className={`text-sm font-medium px-2 py-1 rounded inline-flex items-center ${
                      translated === "รอให้บริการ"
                        ? "bg-yellow-100 text-yellow-800"
                        : translated === "กำลังให้บริการ"
                        ? "bg-blue-100 text-blue-800"
                        : translated === "เสร็จแล้ว"
                        ? "bg-green-100 text-green-800"
                        : "bg-red-100 text-red-800"
                    }`}
                  >
                    {statusIcon[translated]}
                    <span>{translated}</span>
                  </div>

                  {/* ปุ่มเปลี่ยนสถานะ */}
                  {translated === "รอให้บริการ" && (
                    <button
                      className="text-sm text-blue-600 underline"
                      onClick={() => updateStatus(appt.id, "CONFIRMED")}
                    >
                      เริ่มให้บริการ
                    </button>
                  )}
                  {translated === "กำลังให้บริการ" && (
                    <button
                      className="text-sm text-green-600 underline"
                      onClick={() => updateStatus(appt.id, "COMPLETED")}
                    >
                      เสร็จสิ้น
                    </button>
                  )}
                </div>
              </li>
            );
          })
        )}
      </ul>
    </div>
  );
}

function SummaryCard({
  label,
  count,
  color = "gray",
}: {
  label: string;
  count: number;
  color?: "gray" | "green" | "red" | "blue";
}) {
  const colorClass = {
    gray: "text-gray-500",
    green: "text-green-600",
    red: "text-red-600",
    blue: "text-blue-600",
  }[color];

  return (
    <div className="bg-white p-4 rounded shadow text-center">
      <div className={`text-sm ${colorClass}`}>{label}</div>
      <div className="text-xl font-bold">{count}</div>
    </div>
  );
}
