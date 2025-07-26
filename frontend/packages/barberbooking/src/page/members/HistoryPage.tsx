export const mockAppointmentHistory = [

    {
        id: "a001",
        barberName: "ช่างบอล",
        serviceName: "ตัดผมชาย",
        appointmentDate: "2025-06-30",
        appointmentTime: "14:30",
        branchName: "เซ็นทรัลลาดพร้าว",
        status: "เสร็จสิ้น",
        review: { rating: 5, comment: "บริการดีมาก" },
    },
    {
        id: "a002",
        barberName: "ช่างกุ้ง",
        serviceName: "ทำสีผม",
        appointmentDate: "2025-06-15",
        appointmentTime: "10:00",
        branchName: "เมเจอร์รัชโยธิน",
        status: "เสร็จสิ้น",
        review: { rating: 4, comment: "สีสวยแต่รอนาน" },
    },
    {
        id: "a003",
        barberName: "ช่างเบียร์",
        serviceName: "โกนหนวด",
        appointmentDate: "2025-09-01",
        appointmentTime: "18:00",
        branchName: "ฟิวเจอร์พาร์ค",
        status: "ยกเลิก",
        review: null,
    },
    {
        id: "a004",
        barberName: "ช่างหนึ่ง",
        serviceName: "ตัดผมชาย",
        appointmentDate: "2025-12-02",
        appointmentTime: "11:00",
        branchName: "ซีคอนศรีนครินทร์",
        status: "เสร็จสิ้น",
        review: null,
    },
    {
        id: "a005",
        barberName: "ช่างสอง",
        serviceName: "สระไดร์",
        appointmentDate: "2025-06-20",
        appointmentTime: "09:30",
        branchName: "แฟชั่นไอส์แลนด์",
        status: "เสร็จสิ้น",
        review: { rating: 5, comment: "สะอาดดี" },
    },
    {
        id: "a006",
        barberName: "ช่างแดง",
        serviceName: "ทำสีผม",
        appointmentDate: "2025-05-28",
        appointmentTime: "13:00",
        branchName: "เดอะมอลล์บางกะปิ",
        status: "เสร็จสิ้น",
        review: { rating: 3, comment: "สีหลุดเร็ว" },
    },
    {
        id: "a007",
        barberName: "ช่างน้ำ",
        serviceName: "ตัดผมหญิง",
        appointmentDate: "2025-07-01",
        appointmentTime: "16:00",
        branchName: "มาบุญครอง",
        status: "เสร็จสิ้น",
        review: null,
    },
    {
        id: "a008",
        barberName: "ช่างบาส",
        serviceName: "ตัดผมชาย",
        appointmentDate: "2025-06-18",
        appointmentTime: "12:00",
        branchName: "ไอคอนสยาม",
        status: "ยกเลิก",
        review: null,
    },
    {
        id: "a009",
        barberName: "ช่างเฟิร์น",
        serviceName: "ยืดผม",
        appointmentDate: "2025-06-10",
        appointmentTime: "15:00",
        branchName: "เซ็นทรัลเวิลด์",
        status: "เสร็จสิ้น",
        review: { rating: 4, comment: "ตรงดี ไม่เสียผม" },
    },
    {
        id: "a010",
        barberName: "ช่างแนน",
        serviceName: "ดัดลอน",
        appointmentDate: "2025-06-05",
        appointmentTime: "17:30",
        branchName: "พารากอน",
        status: "เสร็จสิ้น",
        review: null,
    },
];
export type AppointmentHistoryItem = {
    id: string | number;
    barberName: string;
    serviceName: string;
    appointmentDate: string;     // รูปแบบ: yyyy-MM-dd
    appointmentTime: string;     // รูปแบบ: HH:mm
    branchName: string;
    status: string;
    review?: {
        rating: number;
        comment: string;
    };
};

import { useEffect, useState } from "react";
import { getAppointmentsByPhone } from "../../api/appointment";
import { format } from "date-fns/format";
export default function HistoryPage() {
    const [phone,setPhone] = useState('')
    const [search, setSearch] = useState("");
    const [startDate, setStartDate] = useState("");
    const [endDate, setEndDate] = useState("");
    const [onlyUpcoming, setOnlyUpcoming] = useState(false);
    const [appointments, setAppointments] = useState<AppointmentHistoryItem[]>([]);

    async function loadAppointmentsByPhone(phone: string) {
        try {
            const raw = await getAppointmentsByPhone(phone)
            const mapped = transformAppointments(raw);
            setAppointments(mapped);
        } catch (err) {
            console.error("Failed to load appointments:", err);
            setAppointments([]);
        }
    }

    const today = new Date().toISOString().split("T")[0];
    function translateStatus(status: string): string {
        switch (status.toUpperCase()) {
            case "COMPLETED":
                return "เสร็จสิ้น";
            case "CONFIRMED":
                return "ยืนยันแล้ว";
            case "PENDING":
                return "รอยืนยัน";
            case "CANCELLED":
                return "ยกเลิก";
            case "NO_SHOW":
                return "ไม่มา";
            default:
                return status;
        }
    }


    function transformAppointments(raw: any[]): AppointmentHistoryItem[] {
        return raw.map((a) => {
            const fullDateTimeStr = `${a.date}T${a.start}`; // "2025-07-27T10:30"
            const startDate = new Date(fullDateTimeStr);

            return {
                id: a.id,
                barberName: a.barber?.username || "ไม่ระบุ",
                serviceName: a.service?.name || "ไม่ระบุ",
                appointmentDate: format(startDate, "yyyy-MM-dd"),
                appointmentTime: format(startDate, "HH:mm"),
                branchName: a.branch?.name || "ไม่ระบุ", // ยังไม่มี? ใส่ "ไม่ระบุ"
                status: translateStatus(a.status),
                review: a.review
                    ? {
                        rating: a.review.rating,
                        comment: a.review.comment,
                    }
                    : undefined,
            };
        });
    }


    useEffect(() => {
        if (phone.length >= 10) {
            loadAppointmentsByPhone(phone);
        }
    }, [phone]);
    return (
        <div className="p-4 space-y-4">
            <h2 className="text-xl font-bold">ประวัติการจองคิว</h2>
            <input
                id="phone-box"
                type="text"
                value={phone}
                onChange={(e) => setPhone(e.target.value)}
                placeholder="กรอกเบอร์โทร"
                className="input input-bordered"
            />

            <div className="flex flex-wrap gap-2 items-center">
                <input
                    type="text"
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                    placeholder="ค้นหาช่างหรือบริการ"
                    className="input input-bordered"
                />
                <input
                    type="date"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                    className="input input-bordered"
                />
                <input
                    type="date"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                    className="input input-bordered"
                />
                <label className="flex items-center gap-1 text-sm">
                    <input
                        type="checkbox"
                        checked={onlyUpcoming}
                        onChange={(e) => setOnlyUpcoming(e.target.checked)}
                        className="checkbox"
                    />
                    แสดงเฉพาะคิวล่วงหน้า
                </label>
            </div>
            <div className="grid gap-2 lg:grid-cols-2">
                {appointments.length === 0 ? (
                    <p className="text-gray-500 italic">ไม่พบข้อมูลที่ตรงกับเงื่อนไข</p>
                ) : (
                    appointments.map((appt) => (
                        <div
                            key={appt.id}
                            className="p-4 border rounded-md shadow-sm bg-white space-y-1"
                        >
                            <div className="font-semibold">
                                {appt.serviceName} กับ {appt.barberName}
                            </div>
                            <div className="text-sm text-gray-600">
                                วันที่: {format(new Date(appt.appointmentDate), "dd-MM-yyyy")} เวลา: {appt.appointmentTime} <br />
                                สาขา: {appt.branchName} | สถานะ: {appt.status}
                            </div>
                            {appt.review ? (
                                <div className="text-yellow-600 text-sm">
                                    ⭐ {appt.review.rating} - "{appt.review.comment}"
                                </div>
                            ) : (
                                <div className="text-sm text-gray-400 italic">ยังไม่ได้รีวิว</div>
                            )}
                        </div>
                    ))
                )}

            </div>

        </div>
    );
}
