import { useMemo } from "react";


export const mockMonthlyIncome = [
    {
        id: "2025-07",
        month: "กรกฎาคม 2568",
        totalIncome: 12500,
        totalAppointments: 54,
        slipUrl: "/slips/slip_2025-07.pdf",
    },
    {
        id: "2025-06",
        month: "มิถุนายน 2568",
        totalIncome: 10200,
        totalAppointments: 48,
        slipUrl: "/slips/slip_2025-06.pdf",
    },
    {
        id: "2025-05",
        month: "พฤษภาคม 2568",
        totalIncome: 8900,
        totalAppointments: 43,
        slipUrl: "/slips/slip_2025-05.pdf",
    },
];



export default function BarberIncome() {
    const currentMonthId = new Date().toISOString().slice(0, 7); // eg. '2025-07'

    const currentMonth = useMemo(() => {
        return mockMonthlyIncome.find((m) => m.id === currentMonthId);
    }, [currentMonthId]);

    return (
        <div className="space-y-6">
            {/* ✅ Header: สรุปรายได้เดือนปัจจุบัน */}
            <div className="bg-white rounded-lg shadow p-6">
                <h2 className="text-xl font-bold mb-2">สรุปรายได้เดือนนี้</h2>
                {currentMonth ? (
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div className="text-lg text-gray-700">
                            🧾 รายได้รวม:{" "}
                            <span className="font-semibold text-green-600">
                                {currentMonth.totalIncome.toLocaleString()} บาท
                            </span>
                        </div>
                        <div className="text-lg text-gray-700">
                            ✂️ จำนวนคิวทั้งหมด:{" "}
                            <span className="font-semibold">
                                {currentMonth.totalAppointments} คิว
                            </span>
                        </div>
                    </div>
                ) : (
                    <p className="text-gray-500">ไม่มีข้อมูลรายได้สำหรับเดือนนี้</p>
                )}
            </div>

            {/* ✅ ตาราง: รายได้ย้อนหลัง */}
            <div className="bg-white rounded-lg shadow p-6">
                <h3 className="text-lg font-bold mb-4">ประวัติรายได้ย้อนหลัง</h3>
                <table className="w-full text-left border border-gray-200">
                    <thead className="bg-gray-100 text-gray-700">
                        <tr>
                            <th className="px-4 py-2 border-b">เดือน</th>
                            <th className="px-4 py-2 border-b">จำนวนคิว</th>
                            <th className="px-4 py-2 border-b">รายได้รวม</th>
                            <th className="px-4 py-2 border-b">สลิป</th>
                        </tr>
                    </thead>
                    <tbody>
                        {mockMonthlyIncome.map((month) => (
                            <tr key={month.id} className="hover:bg-gray-50">
                                <td className="px-4 py-2 border-b">{month.month}</td>
                                <td className="px-4 py-2 border-b">{month.totalAppointments} คิว</td>
                                <td className="px-4 py-2 border-b text-green-700 font-semibold">
                                    {month.totalIncome.toLocaleString()} บาท
                                </td>
                                <td className="px-4 py-2 border-b">
                                    <a
                                        href={month.slipUrl}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className="text-sm text-blue-600 hover:underline"
                                    >
                                        ดูสลิป
                                    </a>
                                </td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
}
