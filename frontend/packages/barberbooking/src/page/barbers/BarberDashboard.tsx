import { useState } from "react";
import {
  Clock,
  CheckCircle,
  XCircle,
  Loader,
  CalendarDays,
  ChevronRight,
} from "lucide-react";

export const initialAppointments = [
  {
    id: 1,
    time: "10:00 - 10:30",
    customerName: "คุณบอย",
    serviceName: "ตัดผมชาย",
    status: "รอให้บริการ",
  },
  {
    id: 2,
    time: "11:00 - 11:45",
    customerName: "คุณตูน",
    serviceName: "ตัด+เซ็ตผม",
    status: "กำลังให้บริการ",
  },
  {
    id: 3,
    time: "13:00 - 14:00",
    customerName: "คุณเมย์",
    serviceName: "ย้อมผม",
    status: "ยกเลิก",
  },
  {
    id: 4,
    time: "15:00 - 15:45",
    customerName: "คุณกานต์",
    serviceName: "ตัดผม+โกนหนวด",
    status: "เสร็จแล้ว",
  },
];
const mockAppointmentsByDate: Record<string, typeof initialAppointments> = {
    "2025-07-03": [
      {
        id: 1,
        time: "10:00 - 10:30",
        customerName: "คุณบอย",
        serviceName: "ตัดผมชาย",
        status: "รอให้บริการ",
      },
      {
        id: 2,
        time: "11:00 - 11:45",
        customerName: "คุณตูน",
        serviceName: "ตัด+เซ็ตผม",
        status: "กำลังให้บริการ",
      },
    ],
    "2025-07-04": [
      {
        id: 3,
        time: "13:00 - 14:00",
        customerName: "คุณเมย์",
        serviceName: "ย้อมผม",
        status: "ยกเลิก",
      },
      {
        id: 4,
        time: "15:00 - 15:45",
        customerName: "คุณกานต์",
        serviceName: "ตัดผม+โกนหนวด",
        status: "เสร็จแล้ว",
      },
    ],
  };
  
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
  
    const appointments = mockAppointmentsByDate[selectedDate] || [];
  
    const updateStatus = (id: number, newStatus: string) => {
      mockAppointmentsByDate[selectedDate] = appointments.map((a) =>
        a.id === id ? { ...a, status: newStatus } : a
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
            count={appointments.filter((a) => a.status === "เสร็จแล้ว").length}
            color="green"
          />
          <SummaryCard
            label="ยกเลิก"
            count={appointments.filter((a) => a.status === "ยกเลิก").length}
            color="red"
          />
          <SummaryCard
            label="กำลังให้บริการ"
            count={
              appointments.filter((a) => a.status === "กำลังให้บริการ").length
            }
            color="blue"
          />
        </div>
  
        {/* Appointments */}
        <ul className="space-y-4">
          {appointments.map((appt) => (
            <li
              key={appt.id}
              className="border rounded p-4 shadow-sm bg-white flex justify-between items-center"
            >
              <div>
                <div className="text-sm text-gray-500">{appt.time}</div>
                <div className="text-lg font-semibold">{appt.customerName}</div>
                <div className="text-sm text-gray-700">{appt.serviceName}</div>
              </div>
  
              <div className="text-right space-y-1">
                <div
                  className={`text-sm font-medium px-2 py-1 rounded inline-flex items-center ${
                    appt.status === "รอให้บริการ"
                      ? "bg-yellow-100 text-yellow-800"
                      : appt.status === "กำลังให้บริการ"
                      ? "bg-blue-100 text-blue-800"
                      : appt.status === "เสร็จแล้ว"
                      ? "bg-green-100 text-green-800"
                      : "bg-red-100 text-red-800"
                  }`}
                >
                  {statusIcon[appt.status]}
                  <span>{appt.status}</span>
                </div>
  
                {/* ปุ่มเปลี่ยนสถานะ */}
                {appt.status === "รอให้บริการ" && (
                  <button
                    className="text-sm text-blue-600 underline"
                    onClick={() => updateStatus(appt.id, "กำลังให้บริการ")}
                  >
                    เริ่มให้บริการ
                  </button>
                )}
                {appt.status === "กำลังให้บริการ" && (
                  <button
                    className="text-sm text-green-600 underline"
                    onClick={() => updateStatus(appt.id, "เสร็จแล้ว")}
                  >
                    เสร็จสิ้น
                  </button>
                )}
              </div>
            </li>
          ))}
        </ul>
  
        {appointments.length === 0 && (
          <p className="text-gray-500 text-center">ไม่มีคิวในวันนี้</p>
        )}
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
  